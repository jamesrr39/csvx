package dal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoder_Decode(t *testing.T) {
	type namedType string

	type targetType struct {
		String      string    `csv:"string"`
		Int         int       `csv:"int"`
		Int64       int64     `csv:"int64"`
		NamedType   namedType `csv:"namedType"`
		PtrInt      *int      `csv:"ptrInt"`
		PtrIntNull  *int      `csv:"ptrIntNull"`
		PtrBool     *bool     `csv:"ptrBool"`
		PtrString   *string   `csv:"ptrString"`
		Float64     float64   `csv:"float64"`
		PtrFloat64  *float64  `csv:"ptrFloat64"`
		NonCSVField string
	}

	ptrIntVal := 21
	ptrBoolVal := true
	ptrStringVal := "PtrStringVal..."
	ptrFloat64Val := 50.5

	type fields struct {
		Fields []string
	}
	type args struct {
		values []string
		target interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wanted  *targetType
		wantErr bool
	}{
		{
			fields: fields{
				Fields: []string{"string", "int", "int64", "namedType", "ptrInt", "ptrIntNull", "ptrBool", "ptrString", "float64", "ptrFloat64"},
			},
			args: args{
				values: []string{
					"Hello World!",
					"50",
					"50",
					"Hello World 2!",
					fmt.Sprintf("%d", *&ptrIntVal),
					"",
					"true",
					ptrStringVal,
					"50.5",
					fmt.Sprintf("%v", ptrFloat64Val),
				},
				target: new(targetType),
			},
			wanted: &targetType{
				String:     "Hello World!",
				Int:        50,
				Int64:      50,
				NamedType:  "Hello World 2!",
				PtrInt:     &ptrIntVal,
				PtrIntNull: nil,
				PtrBool:    &ptrBoolVal,
				PtrString:  &ptrStringVal,
				Float64:    50.5,
				PtrFloat64: &ptrFloat64Val,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decoder{
				Fields: tt.fields.Fields,
			}
			if err := d.Decode(tt.args.values, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.wanted, tt.args.target)
		})
	}
}
