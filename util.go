package csvx

import (
	"fmt"
	"log"
	"reflect"
)

func traverseFields(target interface{}, createMissingStructs bool, fn func(fieldCsvTag string, field reflect.Value /*field reflect.StructField*/) error) error {
	rv := reflect.ValueOf(target)
	rt := reflect.TypeOf(target)
	if rt.Kind() == reflect.Pointer {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	println("type::", rv.Type().String(), rv.CanAddr(), rv.Type().String())

	for i := 0; i < rt.NumField(); i++ {
		fieldV := rv.Field(i)
		fieldT := rt.Field(i)

		const csvTagName = "csv"

		csvTag := fieldT.Tag.Get(csvTagName)
		if csvTag != "" {
			err := fn(csvTag, fieldV)
			if err != nil {
				return err
			}
		}

		fieldUnderlyingKind := getUnderlyingObject(fieldV).Kind()

		// if field is a struct, go into that struct and look for tags there
		if fieldUnderlyingKind == reflect.Struct {
			if csvTag != "" {
				return fmt.Errorf("csvx: %q tag on anonymous field not supported", csvTagName)
			}

			fieldOrCreatedObjectV := fieldV
			if createMissingStructs {
				println("Creating missing structs...")
				// https://github.com/robertkrimen/otto/issues/83
				// https: //go.dev/blog/laws-of-reflection
				fieldUnderlyingInterface := fieldV.Interface()
				// if fieldV.Kind() == reflect.Pointer {
				// 	fieldUnderlyingInterface = &fieldUnderlyingInterface
				// }
				fieldOrCreatedObjectV = reflect.New(reflect.Indirect(reflect.ValueOf(fieldUnderlyingInterface)).Type())

			}

			ni := fieldOrCreatedObjectV.Interface()
			println("going into::", fieldV.Type().Name(), "::", reflect.TypeOf(ni).String(), reflect.TypeOf(&ni).String())
			err := traverseFields(ni, createMissingStructs, fn)
			if err != nil {
				return err
			}

			if createMissingStructs {

				if fieldV.Kind() != reflect.Pointer {
					fieldOrCreatedObjectV = fieldOrCreatedObjectV.Elem()
				}

				log.Printf("type:: %T || %s\n", fieldOrCreatedObjectV.Interface(), fieldV.Kind().String())

				fieldV.Set(fieldOrCreatedObjectV)
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
