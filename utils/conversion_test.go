package utils

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

const TEST_MSG_INVALID_CONVERSION_VALUE = "Invalid conversion value"

func TestConvertInterfaceValue_String(t *testing.T) {
	v := "This is a string value"
	assert.Equal(t, v, ConvertInterfaceValue(v), TEST_MSG_INVALID_CONVERSION_VALUE)
}

func TestConvertInterfaceValue_Bool(t *testing.T) {
	assert.Equal(t, true, ConvertInterfaceValue(true), TEST_MSG_INVALID_CONVERSION_VALUE)
	assert.Equal(t, false, ConvertInterfaceValue(false), TEST_MSG_INVALID_CONVERSION_VALUE)
}

func TestConvertInterfaceValue_Int(t *testing.T) {
	assert.Equal(t, 10, ConvertInterfaceValue(10), TEST_MSG_INVALID_CONVERSION_VALUE)
	assert.Equal(t, 1000000000000, ConvertInterfaceValue(1000000000000), TEST_MSG_INVALID_CONVERSION_VALUE)
}

func TestConvertInterfaceValue_List(t *testing.T) {
	l1 := []interface{}{1, 2, 3, 4}
	assert.Equal(t, l1, ConvertInterfaceValue(l1), TEST_MSG_INVALID_CONVERSION_VALUE)
	l2 := []interface{}{map[string]interface{}{"payload": "one,two,three"}, map[string]interface{}{"payload": "one,two,three", "separator": ","}}
	assert.Equal(t, l2, ConvertInterfaceValue(l2), TEST_MSG_INVALID_CONVERSION_VALUE)
}

func TestConvertInterfaceValue_Map(t *testing.T) {
	json_1 := map[string]interface{}{"payload": "one,two,three"}
	assert.Equal(t, json_1, ConvertInterfaceValue(json_1), TEST_MSG_INVALID_CONVERSION_VALUE)

	json_2 := map[string]interface{}{"payload": "one,two,three", "separator": ","}
	assert.Equal(t, json_2, ConvertInterfaceValue(json_2), TEST_MSG_INVALID_CONVERSION_VALUE)

	json_3 := map[string]interface{}{"payload": "one,two,three", "lines": []interface{}{"one", "two", "three"}}
	assert.Equal(t, json_3, ConvertInterfaceValue(json_3), TEST_MSG_INVALID_CONVERSION_VALUE)

	json_4 := map[string]interface{}{"p": map[string]interface{}{"a": 1}}
	assert.Equal(t, json_4, ConvertInterfaceValue(json_4), TEST_MSG_INVALID_CONVERSION_VALUE)

	json_5 := map[string]interface{}{"p": map[string]interface{}{"a": 1, "b": 2}}
	assert.Equal(t, json_5, ConvertInterfaceValue(json_5), TEST_MSG_INVALID_CONVERSION_VALUE)

	json_6 := map[string]interface{}{"p": map[string]interface{}{"a": 1, "b": map[string]interface{}{"c": 2}}}
	assert.Equal(t, json_6, ConvertInterfaceValue(json_6), TEST_MSG_INVALID_CONVERSION_VALUE)

	json_7 := map[string]interface{}{"p": map[string]interface{}{"a": 1, "b": map[string]interface{}{"c": 2, "d": 3}}}
	assert.Equal(t, json_7, ConvertInterfaceValue(json_7), TEST_MSG_INVALID_CONVERSION_VALUE)

	json_8 := map[string]interface{}{"p": map[string]interface{}{"a": 1, "b": map[string]interface{}{"c": 2, "d": []interface{}{3, 4}}}}
	assert.Equal(t, json_8, ConvertInterfaceValue(json_8), TEST_MSG_INVALID_CONVERSION_VALUE)

	json_9 := map[string]interface{}{"p": map[string]interface{}{"a": 99.99}}
	assert.Equal(t, json_9, ConvertInterfaceValue(json_9), TEST_MSG_INVALID_CONVERSION_VALUE)

	json_10 := map[string]interface{}{"p": map[string]interface{}{"a": true}}
	assert.Equal(t, json_10, ConvertInterfaceValue(json_10), TEST_MSG_INVALID_CONVERSION_VALUE)
}
