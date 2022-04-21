package csvx

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoder_Encode(t *testing.T) {
	type objectType struct {
		ID                      int64   `csv:"id"`
		Age                     *uint   `csv:"age"`
		Name                    string  `csv:"name"`
		Score                   float64 `csv:"score"`
		IsAdult                 bool    `csv:"isAdult"`
		IgnoreField             bool
		ignorePackageLocalField string
	}

	encoder := NewEncoderWithDefaultOpts([]string{"id", "age", "name", "score", "isAdult"})

	t.Run("object1", func(t *testing.T) {
		object1 := objectType{
			ID:          20,
			Age:         nil,
			Name:        "Test1",
			Score:       1.23,
			IsAdult:     false,
			IgnoreField: true,
		}

		fields, err := encoder.Encode(object1)
		require.NoError(t, err)
		assert.Equal(t, []string{"20", "null", "Test1", "1.23", "false"}, fields)
	})

	t.Run("object2", func(t *testing.T) {
		age2 := uint(50)
		object2 := &objectType{
			ID:          21,
			Age:         &age2,
			Name:        "Test2",
			Score:       2.5,
			IsAdult:     false,
			IgnoreField: true,
		}

		fields, err := encoder.Encode(object2)
		require.NoError(t, err)
		assert.Equal(t, []string{"21", "50", "Test2", "2.5", "false"}, fields)
	})
}

func ExampleEncoder() {
	type objectType struct {
		ID     int64   `csv:"id"`
		Name   string  `csv:"name"`
		Height float64 `csv:"height"`
	}
	obj := objectType{
		ID:     123,
		Name:   "Test 987",
		Height: 1.567,
	}

	encoder := NewEncoderWithDefaultOpts([]string{"id", "name", "height"})

	fields, err := encoder.Encode(obj)
	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(os.Stdout)
	err = w.Write(fields)
	if err != nil {
		panic(err)
	}
	w.Flush()

	err = w.Error()
	if err != nil {
		panic(err)
	}

	// Output:
	// 123,Test 987,1.567
}

func Test_embedded_struct_encode(t *testing.T) {
	type EmbeddedType struct {
		Field1 string `csv:"field1"`
	}

	type EmbeddedType2 struct {
		Field2 int     `csv:"field2"`
		Field3 float64 `csv:"field3"`
	}

	type SubType struct {
		Field4 string `csv:"field4"`
	}

	type SubType2 struct {
		Field5 string `csv:"field5"`
	}

	type myType struct {
		*EmbeddedType // pointer
		EmbeddedType2 // not pointer
		SubType       SubType
		SubType2      *SubType2
	}

	obj := myType{
		EmbeddedType:  &EmbeddedType{Field1: "Test1"},
		EmbeddedType2: EmbeddedType2{Field2: 50},
		SubType:       SubType{Field4: "Test2"},
		SubType2:      &SubType2{Field5: "Test3"},
	}

	encoder := NewEncoderWithDefaultOpts([]string{"field2", "field1", "field3", "field4", "field5"})
	fields, err := encoder.Encode(obj)
	require.NoError(t, err)

	expectedResults := []string{"50", "Test1", "0", "Test2", "Test3"}

	assert.Equal(t, expectedResults, fields)
}
