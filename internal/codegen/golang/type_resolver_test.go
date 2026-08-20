package golang

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mypricehealth/pggen/internal/casing"
	"github.com/mypricehealth/pggen/internal/codegen/golang/gotype"
	"github.com/mypricehealth/pggen/internal/difftest"
	"github.com/mypricehealth/pggen/internal/pg"
	"github.com/stretchr/testify/assert"
)

func TestTypeResolver_Resolve(t *testing.T) {
	testPkgPath := "github.com/mypricehealth/pggen/internal/codegen/golang/test_resolve"
	caser := casing.NewCaser()
	caser.AddAcronym("ios", "IOS")
	caser.AddAcronym("macos", "MacOS")
	caser.AddAcronym("id", "ID")
	pgDeviceEnum := pg.EnumType{Name: "device_type", Labels: []string{"macos", "ios", "web"}}
	pgReqIntDomain := pg.DomainType{Name: "req_int", IsNotNull: true, Elem: pg.Int8}
	goDeviceEnum := &gotype.EnumType{
		PgEnum: pgDeviceEnum,
		Name:   "DeviceType",
		Labels: []string{"DeviceTypeMacOS", "DeviceTypeIOS", "DeviceTypeWeb"},
		Values: []string{"macos", "ios", "web"},
	}
	tests := []struct {
		name      string
		overrides map[string]string
		pgType    pg.Type
		nullable  bool
		want      gotype.Type
	}{
		{
			name:   "enum",
			pgType: pgDeviceEnum,
			want:   &gotype.ImportType{PkgPath: testPkgPath, Type: goDeviceEnum},
		},
		{
			name:   "enum array",
			pgType: pg.ArrayType{Name: "_device_type", Elem: pgDeviceEnum},
			want: &gotype.ArrayType{
				PgArray: pg.ArrayType{Name: "_device_type", Elem: pgDeviceEnum},
				Elem:    &gotype.ImportType{PkgPath: testPkgPath, Type: goDeviceEnum},
			},
		},
		{
			name:     "array element ignores the column marker",
			pgType:   pg.ArrayType{Name: "_device_type", Elem: pgDeviceEnum},
			nullable: true,
			want: &gotype.ArrayType{
				PgArray: pg.ArrayType{Name: "_device_type", Elem: pgDeviceEnum},
				Elem:    &gotype.ImportType{PkgPath: testPkgPath, Type: goDeviceEnum},
			},
		},
		{
			name:     "not null domain element survives the column marker",
			pgType:   pg.ArrayType{Name: "_req_int", Elem: pgReqIntDomain},
			nullable: true,
			want: &gotype.ArrayType{
				PgArray: pg.ArrayType{Name: "_req_int", Elem: pgReqIntDomain},
				Elem: &gotype.DomainType{
					Name:     "ReqInt",
					PgDomain: pgReqIntDomain,
					Elem:     &gotype.OpaqueType{PgType: pg.Int8, Name: "int"},
				},
			},
		},
		{
			name:   "void",
			pgType: pg.VoidType{},
			want:   &gotype.VoidType{},
		},
		{
			name:      "override",
			overrides: map[string]string{"custom_type": "example.com/custom.QualType"},
			pgType:    pg.BaseType{Name: "custom_type"},
			want: &gotype.ImportType{
				PkgPath: "example.com/custom",
				Type:    &gotype.OpaqueType{PgType: pg.BaseType{Name: "custom_type"}, Name: "QualType"},
			},
		},
		{
			name:      "override pointer",
			overrides: map[string]string{"custom_type": "*example.com/custom.QualType"},
			pgType:    pg.BaseType{Name: "custom_type"},
			want: &gotype.PointerType{
				Elem: &gotype.ImportType{
					PkgPath: "example.com/custom",
					Type:    &gotype.OpaqueType{PgType: pg.BaseType{Name: "custom_type"}, Name: "QualType"},
				},
			},
		},
		{
			name:      "override pointer array",
			overrides: map[string]string{"_custom_type": "[]*example.com/custom.QualType"},
			pgType:    pg.ArrayType{Name: "_custom_type", Elem: pg.BaseType{Name: "custom_type"}},
			want: &gotype.ArrayType{
				PgArray: pg.ArrayType{Name: "_custom_type", Elem: pg.BaseType{Name: "custom_type"}},
				Elem: &gotype.PointerType{
					Elem: &gotype.ImportType{
						PkgPath: "example.com/custom",
						Type:    &gotype.OpaqueType{Name: "QualType"},
					},
				},
			},
		},
		{
			name:     "known nonNullable empty",
			pgType:   pg.BaseType{Name: "point", ID: pgtype.PointOID},
			nullable: false,
			want: &gotype.ImportType{
				PkgPath: "github.com/jackc/pgx/v5/pgtype",
				Type: &gotype.OpaqueType{
					PgType: pg.BaseType{Name: "point", ID: pgtype.PointOID},
					Name:   "Point",
				},
			},
		},
		{
			name:     "known nullable",
			pgType:   pg.BaseType{Name: "point", ID: pgtype.PointOID},
			nullable: true,
			want: &gotype.PointerType{Elem: &gotype.ImportType{
				PkgPath: "github.com/jackc/pgx/v5/pgtype",
				Type: &gotype.OpaqueType{
					PgType: pg.BaseType{Name: "point", ID: pgtype.PointOID},
					Name:   "Point",
				},
			}},
		},
		{
			name:     "known nullable pointer variant is not double wrapped",
			pgType:   pg.Text,
			nullable: true,
			want:     &gotype.PointerType{Elem: &gotype.OpaqueType{Name: "string", PgType: pg.Text}},
		},
		{
			name:     "enum nullable",
			pgType:   pgDeviceEnum,
			nullable: true,
			want:     &gotype.PointerType{Elem: &gotype.ImportType{PkgPath: testPkgPath, Type: goDeviceEnum}},
		},
		{
			name:      "bigint - int8",
			overrides: map[string]string{"bigint": "example.com/custom.QualType"},
			pgType:    pg.BaseType{Name: "int8", ID: pgtype.Int8OID},
			want: &gotype.ImportType{
				PkgPath: "example.com/custom",
				Type: &gotype.OpaqueType{
					PgType: pg.BaseType{Name: "int8", ID: pgtype.Int8OID},
					Name:   "QualType",
				},
			},
		},
		{
			name:      "_bigint - _int8",
			overrides: map[string]string{"_bigint": "[]uint16"},
			pgType:    pg.ArrayType{Name: "_int8", Elem: pg.BaseType{Name: "int8", ID: pgtype.Int8OID}},
			want: &gotype.ArrayType{
				PgArray: pg.ArrayType{Name: "_int8", Elem: pg.BaseType{Name: "int8", ID: pgtype.Int8OID}},
				Elem:    &gotype.OpaqueType{Name: "uint16"},
			},
		},
		{
			name:   "_bigint - _int8",
			pgType: pg.BaseType{Name: "jsonb", ID: pgtype.JSONBOID},
			want: &gotype.ImportType{
				PkgPath: "encoding/json",
				Type: &gotype.OpaqueType{
					PgType: pg.BaseType{Name: "jsonb", ID: pgtype.JSONBOID},
					Name:   "RawMessage",
				},
			},
		},
		{
			name:      "_real - _float4 custom type",
			overrides: map[string]string{"_real": "[]example.com/custom.F32"},
			pgType:    pg.ArrayType{ID: pgtype.Float4ArrayOID, Name: "_float4", Elem: pg.BaseType{Name: "_float4", ID: pgtype.Float4OID}},
			want: &gotype.ArrayType{
				PgArray: pg.ArrayType{ID: pgtype.Float4ArrayOID, Name: "_float4", Elem: pg.BaseType{Name: "_float4", ID: pgtype.Float4OID}},
				Elem: &gotype.ImportType{
					PkgPath: "example.com/custom",
					Type:    &gotype.OpaqueType{Name: "F32"},
				},
			},
		},
		{
			name: "composite",
			pgType: pg.CompositeType{
				Name:        "qux",
				ColumnNames: []string{"id", "foo"},
				ColumnTypes: []pg.Type{pg.Text, pg.Int8},
			},
			nullable: true,
			want: &gotype.PointerType{Elem: &gotype.ImportType{
				PkgPath: testPkgPath,
				Type: &gotype.CompositeType{
					PgComposite: pg.CompositeType{
						Name:        "qux",
						ColumnNames: []string{"id", "foo"},
						ColumnTypes: []pg.Type{pg.Text, pg.Int8},
					},
					Name:       "Qux",
					FieldNames: []string{"ID", "Foo"},
					FieldTypes: []gotype.Type{
						&gotype.PointerType{Elem: &gotype.OpaqueType{Name: "string", PgType: pg.Text}},
						&gotype.PointerType{Elem: &gotype.OpaqueType{Name: "int", PgType: pg.Int8}},
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewTypeResolver(caser, tt.overrides)
			got, err := resolver.Resolve(tt.pgType, tt.nullable, testPkgPath)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestType_QualifyRel(t *testing.T) {
	caser := casing.NewCaser()
	tests := []struct {
		typ          gotype.Type
		otherPkgPath string
		want         string
	}{
		{
			typ: gotype.NewEnumType(
				"example.com/foo",
				pg.EnumType{Name: "device", Labels: []string{"macos"}},
				caser,
			),
			otherPkgPath: "example.com/bar",
			want:         "foo.Device",
		},
		{
			typ: gotype.NewEnumType(
				"example.com/bar",
				pg.EnumType{Name: "device", Labels: []string{"macos"}},
				caser,
			),
			otherPkgPath: "example.com/bar",
			want:         "Device",
		},
		{
			typ:          gotype.MustParseOpaqueType("example.com/bar.Baz"),
			otherPkgPath: "example.com/bar",
			want:         "Baz",
		},
		{
			typ:          gotype.MustParseKnownType("string", pg.Text),
			otherPkgPath: "example.com/bar",
			want:         "string",
		},
		{
			typ:          gotype.MustParseKnownType("string", pg.Text),
			otherPkgPath: "",
			want:         "string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.typ.Import()+"."+tt.typ.BaseName(), func(t *testing.T) {
			got := gotype.QualifyType(tt.typ, tt.otherPkgPath)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateCompositeType(t *testing.T) {
	caser := casing.NewCaser()
	resolver := NewTypeResolver(caser, nil)
	pgImage := pg.CompositeType{
		Name:        "image",
		ColumnNames: []string{"source"},
		ColumnTypes: []pg.Type{pg.Text},
	}
	pgImageArray := pg.ArrayType{Name: "_image", Elem: pgImage}
	pgImageSet := pg.CompositeType{
		Name:        "image_set",
		ColumnNames: []string{"orig", "rest"},
		ColumnTypes: []pg.Type{pgImage, pgImageArray},
	}
	goImage := &gotype.ImportType{
		PkgPath: "example.com/foo",
		Type: &gotype.CompositeType{
			PgComposite: pgImage,
			Name:        "Image",
			FieldNames:  []string{"Source"},
			FieldTypes: []gotype.Type{
				&gotype.PointerType{Elem: &gotype.OpaqueType{PgType: pg.Text, Name: "string"}},
			},
		},
	}
	pgReqInt := pg.DomainType{Name: "req_int", IsNotNull: true, Elem: pg.Int8}
	pgOptInt := pg.DomainType{Name: "opt_int", Elem: pg.Int8}
	pgTimeOfDay := pg.CompositeType{
		Name:        "time_of_day_type",
		ColumnNames: []string{"hour", "minute"},
		ColumnTypes: []pg.Type{pgReqInt, pgOptInt},
	}
	pgAppointment := pg.CompositeType{
		Name:           "appointment",
		ColumnNames:    []string{"label", "note"},
		ColumnTypes:    []pg.Type{pg.Int8, pg.Text},
		ColumnNotNulls: []bool{true, false},
	}
	tests := []struct {
		pkgPath string
		pgType  pg.CompositeType
		want    gotype.Type
	}{
		{
			pkgPath: "example.com/foo",
			pgType: pg.CompositeType{
				Name:        "qux",
				ColumnNames: []string{"one", "two_a"},
				ColumnTypes: []pg.Type{pg.Text, pg.Int8},
			},
			want: &gotype.ImportType{
				PkgPath: "example.com/foo",
				Type: &gotype.CompositeType{
					PgComposite: pg.CompositeType{
						Name:        "qux",
						ColumnNames: []string{"one", "two_a"},
						ColumnTypes: []pg.Type{pg.Text, pg.Int8},
					},
					Name:       "Qux",
					FieldNames: []string{"One", "TwoA"},
					FieldTypes: []gotype.Type{
						&gotype.PointerType{Elem: &gotype.OpaqueType{PgType: pg.Text, Name: "string"}},
						&gotype.PointerType{Elem: &gotype.OpaqueType{PgType: pg.Int8, Name: "int"}},
					},
				},
			},
		},
		{
			pkgPath: "example.com/foo",
			pgType:  pgImageSet,
			want: &gotype.ImportType{
				PkgPath: "example.com/foo",
				Type: &gotype.CompositeType{
					PgComposite: pgImageSet,
					Name:        "ImageSet",
					FieldNames:  []string{"Orig", "Rest"},
					FieldTypes: []gotype.Type{
						&gotype.PointerType{Elem: goImage},
						&gotype.ArrayType{PgArray: pgImageArray, Elem: goImage},
					},
				},
			},
		},
		{
			pkgPath: "example.com/foo",
			pgType:  pgTimeOfDay,
			want: &gotype.ImportType{
				PkgPath: "example.com/foo",
				Type: &gotype.CompositeType{
					PgComposite: pgTimeOfDay,
					Name:        "TimeOfDayType",
					FieldNames:  []string{"Hour", "Minute"},
					FieldTypes: []gotype.Type{
						&gotype.DomainType{
							Name:     "ReqInt",
							PgDomain: pgReqInt,
							Elem:     &gotype.OpaqueType{PgType: pg.Int8, Name: "int"},
						},
						&gotype.DomainType{
							Name:     "OptInt",
							PgDomain: pgOptInt,
							Elem:     &gotype.PointerType{Elem: &gotype.OpaqueType{PgType: pg.Int8, Name: "int"}},
						},
					},
				},
			},
		},
		{
			pkgPath: "example.com/foo",
			pgType:  pgAppointment,
			want: &gotype.ImportType{
				PkgPath: "example.com/foo",
				Type: &gotype.CompositeType{
					PgComposite: pgAppointment,
					Name:        "Appointment",
					FieldNames:  []string{"Label", "Note"},
					FieldTypes: []gotype.Type{
						&gotype.OpaqueType{PgType: pg.Int8, Name: "int"},
						&gotype.PointerType{Elem: &gotype.OpaqueType{PgType: pg.Text, Name: "string"}},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.pkgPath+" "+tt.pgType.Name, func(t *testing.T) {
			got, err := CreateCompositeType(tt.pkgPath, tt.pgType, resolver, caser)
			assert.NoError(t, err)
			difftest.AssertSame(t, tt.want, got)
		})
	}
}
