package csvx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Decoder struct {
	Fields []string
}

func NewDecoder(fields []string) *Decoder {
	return &Decoder{fields}
}

type InvalidUnmarshalError struct {
	goType reflect.Type
}

func (e InvalidUnmarshalError) Error() string {
	return fmt.Sprintf("csvx: unmarshal target must be a non-nil pointer, but got %s", e.goType)
}

func (d *Decoder) Decode(values []string, target interface{}) error {
	// check target is a non-nil pointer
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(target)}
	}

	if len(values) != len(d.Fields) {
		return fmt.Errorf("csvx: amount of fields (%d) does not match amount of values passed in (%d)", len(d.Fields), len(values))
	}

	elem := reflect.TypeOf(target).Elem()

	fieldIndexByName := make(map[string]int)
	for i := 0; i < rv.Elem().NumField(); i++ {
		field := elem.Field(i)

		fieldTag := field.Tag.Get("csv")
		if fieldTag == "" {
			// no "csv" tag set, skip this field
			continue
		}
		fieldIndexByName[fieldTag] = i
	}

	for i, fieldName := range d.Fields {
		fieldIndex, ok := fieldIndexByName[fieldName]
		if !ok {
			return fmt.Errorf("csv: could not find field %q in struct. Make sure the tag 'csv' is set.", fieldName)
		}

		field := rv.Elem().Field(fieldIndex)

		valueStr := values[i]

		err := setField(field, field.Kind(), valueStr, false)
		if err != nil {
			return fmt.Errorf("csvx: error setting field. Value: %q, field index: %d. Error: %s", err, i, valueStr)
		}
	}

	return nil
}

func setField(field reflect.Value, fieldKind reflect.Kind, valueStr string, isPtr bool) error {
	switch fieldKind {
	case reflect.String:
		if isPtr {
			field.Set(reflect.ValueOf(&valueStr))
		} else {
			field.SetString(valueStr)
		}
	case reflect.Int:
		val, err := strconv.Atoi(valueStr)
		if err != nil {
			return err
		}

		if isPtr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.SetInt(int64(val))
		}
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		val, err := strconv.ParseInt(valueStr, 10, bitSizeFromKind(fieldKind))
		if err != nil {
			return err
		}

		if isPtr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.SetInt(int64(val))
		}
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		val, err := strconv.ParseUint(valueStr, 10, bitSizeFromKind(fieldKind))
		if err != nil {
			return err
		}

		if isPtr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.SetUint(val)
		}
	case reflect.Float64, reflect.Float32:
		val, err := strconv.ParseFloat(valueStr, bitSizeFromKind(fieldKind))
		if err != nil {
			return err
		}

		if isPtr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.SetFloat(val)
		}
	case reflect.Bool:
		val, err := boolValueFromStr(valueStr)
		if err != nil {
			return err
		}
		if isPtr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.SetBool(val)
		}
	case reflect.Pointer:
		if valueStr == "" || valueStr == "null" {
			// leave field nil
			return nil
		}

		setField(field, field.Type().Elem().Kind(), valueStr, true)
	default:
		return fmt.Errorf("field type not implemented: %s", fieldKind)
	}

	return nil
}

func boolValueFromStr(valueStr string) (bool, error) {
	switch strings.ToLower(valueStr) {
	case "true", "yes", "1", "1.0":
		return true, nil
	case "false", "no", "0", "0.0":
		return false, nil
	}

	return false, fmt.Errorf("couldn't understand value that should be a boolean field")
}

func bitSizeFromKind(kind reflect.Kind) int {
	switch kind {
	case reflect.Int64, reflect.Float64:
		return 64
	case reflect.Int32, reflect.Float32:
		return 32
	case reflect.Int16:
		return 16
	case reflect.Int8:
		return 8
	}

	panic(fmt.Sprintf("kind not handled: %s", kind))
}
