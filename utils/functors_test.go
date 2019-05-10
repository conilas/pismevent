package utils

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
  value := Sum(0,1)

  assert.Equal(t, value, float64(1), "The sum of the values should be equal")
}

func TestReduceNumbersWithSum(t *testing.T) {
  testedFunction  := Sum
  values          := [...]float64{1,2,3,4,5}

  //convert the values to []interface because.. like, you know.. golang
  interfaceValues := make([]interface{}, len(values))
  for i := range values {
    interfaceValues[i] = values[i]
  }

  reducedValues  := ReduceNumbers(interfaceValues, testedFunction)

  assert.Equal(t, reducedValues, float64(1+2+3+4+5), "The sum of the values should be equal")
}

func TestMap(t *testing.T) {

  type TwoValues struct {
    first string
    second string
  }

  type SingleValue struct {
    first string
  }

  expected := SingleValue{first: "hello,world"}

  mapFunction := func(v interface{}) interface{} {
    workingOn := v.(TwoValues)
    return SingleValue{first: workingOn.first + workingOn.second}
  }

  testedValues := [...] TwoValues {TwoValues{first: "hello,", second: "world"}}

  //convert the values to []interface because.. like, you know.. golang
  interfaceValues := make([]interface{}, len(testedValues))
  for i := range testedValues {
    interfaceValues[i] = testedValues[i]
  }

  mappedValues := Map(interfaceValues, mapFunction)

  assert.Equal(t, mappedValues[0].(SingleValue), expected, "The sum of the values should be equal")
}
