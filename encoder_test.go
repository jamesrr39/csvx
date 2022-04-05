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

	age2 := uint(50)
	object2 := &objectType{
		ID:          21,
		Age:         &age2,
		Name:        "Test2",
		Score:       2.5,
		IsAdult:     false,
		IgnoreField: true,
	}

	fields, err = encoder.Encode(object2)
	require.NoError(t, err)
	assert.Equal(t, []string{"21", "50", "Test2", "2.5", "false"}, fields)
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
