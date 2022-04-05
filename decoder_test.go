package csvx

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	fields := []string{"string", "int", "int64", "namedType", "ptrInt", "ptrIntNull", "ptrBool", "ptrString", "float64", "ptrFloat64"}
	decoder := NewDecoderWithDefaultOpts(fields)

	csvData := bytes.NewBufferString(`Hello World!,50,50,Hello World 2!,21,,true,PtrStringVal...,50.5,50.5`)

	wanted := &targetType{
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
	}

	var results []*targetType

	reader := csv.NewReader(csvData)
	for {
		valueStrings, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			require.NoError(t, err)
		}

		target := new(targetType)
		err = decoder.Decode(valueStrings, target)
		require.NoError(t, err)

		results = append(results, target)
	}

	require.Len(t, results, 1)
	assert.Equal(t, wanted, results[0])
}

func ExampleDecoder() {
	// setup types. Note "csv" field tag.
	type targetType struct {
		Name string `csv:"name"`
		Age  *int   `csv:"age"`
	}

	fields := []string{"name", "age"}
	decoder := NewDecoderWithDefaultOpts(fields)

	csvData := bytes.NewBufferString("John Smith,40\nJane Doe,")

	var results []*targetType

	// use stdlib csv reader to read line by line []string slices
	reader := csv.NewReader(csvData)
	for {
		valueStrings, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			// unexpected error
			panic(err)
		}

		target := new(targetType)
		err = decoder.Decode(valueStrings, target)
		if err != nil {
			panic(err)
		}

		results = append(results, target)
	}

	fmt.Printf("Found %d results\n", len(results))
	for _, result := range results {
		age := "nil"
		if result.Age != nil {
			age = fmt.Sprintf("%d", *result.Age)
		}
		fmt.Printf("%s: %s\n", result.Name, age)
	}

	// Output:
	// Found 2 results
	// John Smith: 40
	// Jane Doe: nil
}
