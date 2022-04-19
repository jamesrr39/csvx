package csvx

import (
	"fmt"
	"reflect"
)

func buildFieldIndexByName(target interface{}, fn func(fieldCsvTag string, field reflect.Value /*field reflect.StructField*/) error) error {
	rv := reflect.ValueOf(target)
	rt := reflect.TypeOf(target)
	if rt.Kind() == reflect.Pointer {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	for i := 0; i < rt.NumField(); i++ {
		fieldV := rv.Field(i)
		fieldT := rt.Field(i)

		csvTag := fieldT.Tag.Get("csv")
		if csvTag != "" {
			err := fn(csvTag, fieldV)
			if err != nil {
				return err
			}
		}

		if fieldT.Anonymous {
			if csvTag != "" {
				return fmt.Errorf("csv tag on anonymous field not supported")
			}

			err := buildFieldIndexByName(fieldV.Interface(), fn)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getUnderlyingObject(rv reflect.Value) reflect.Value {
	if rv.Kind() == reflect.Pointer {
		return rv.Elem()
	}

	return rv
}
