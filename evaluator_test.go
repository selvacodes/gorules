package gorules_test

import (
	"fmt"
	"gorules"
	"testing"

	"github.com/stretchr/testify/assert"
)

var evaluatorTestDataString = `{
  "id": 25,
  "zip5": 33076,
  "zip3": "333",
  "state": "FL",
  "country": "SOUTH A",
  "subtotal": "25.00",
  "promoamount": 1.00,
  "testobj":{
      "id": "3",
      "productId": 34354,
      "quantity": 3,
      "warehouse": "",
      "name": "test product 1",
      "availableInventory": "",
      "promos": []
   },
  "orderItems": [
    {
      "id": 3,
      "productId": 34354,
      "quantity": 3,
      "warehouse": "",
      "weight" : "11",
      "name": "test product 1",
      "availableInventory": "D",
      "promos": []
    },{
      "id": 3,
      "productId": 34354,
      "quantity": 3,
      "weight" : "11",
      "warehouse": "",
      "name": "test product 1",
      "availableInventory": "NC",
      "promos": []
    }
  ],
  "promos": []
}`

var evaluatorTestData = gorules.ParseStringToJSON(evaluatorTestDataString)

func TestCompareValueAndValue(t *testing.T) {
	result := gorules.DSLEvaluator("'SOUTH A' IsEqualTo 'SOUTH A'", evaluatorTestData)
	assert.True(t, result)
}

var testStringSlice = `country IsEqualTo '10' AND
                       country IsEqualTo '100' OR
                       country IsEqualTo 'USA' OR
					   country IsEqualTo 'CANADA'
                       AND
                       state IsEqualTo 'FL' AND
					   country IsEqualTo 'USA'`

func TestCompareValueAndPath(t *testing.T) {
	result := gorules.DSLEvaluator("'SOUTH A' IsEqualTo country", evaluatorTestData)
	assert.True(t, result)
}

func TestComparePathAndValue(t *testing.T) {
	result := gorules.DSLEvaluator("country IsEqualTo 'SOUTH A'", evaluatorTestData)
	assert.True(t, result)
}

func TestComparePathAndPath(t *testing.T) {
	result := gorules.DSLEvaluator("country IsEqualTo country AND country IsEqualTo country", evaluatorTestData)
	assert.True(t, result)
}

func TestAllSelectorPass(t *testing.T) {
	result := gorules.DSLEvaluator("ALL orderItems.weight IsEqualTo '11'", evaluatorTestData)
	assert.True(t, result)
}

func TestAllSelectorFail(t *testing.T) {
	result := gorules.DSLEvaluator("ALL orderItems.weight IsEqualTo '0'", evaluatorTestData)
	assert.False(t, result)
}

func TestAnySelectorPass(t *testing.T) {
	result := gorules.DSLEvaluator("ANY orderItems.availableInventory IsEqualTo 'NC'", evaluatorTestData)
	assert.True(t, result)
}

func TestAnySelectorFail(t *testing.T) {
	result := gorules.DSLEvaluator("ANY orderItems.weight IsEqualTo 'NV'", evaluatorTestData)
	assert.False(t, result)
}

func TestSingleConjunction(t *testing.T) {
	result := gorules.DSLEvaluator("OR", evaluatorTestData)
	assert.False(t, result)
}

func TestWithPrecedence(t *testing.T) {
	var testStringSlice = `country IsEqualTo 'CANADA'
					  	   AND
					   	   country IsEqualTo 'CANADA' AND
					        country IsEqualTo 'CANADA' AND
					   		 state IsEqualTo 'L'`
	result := gorules.EvaluateRules(testStringSlice, evaluatorTestData)
	fmt.Println(result)
	assert.False(t, false)
}

func TestWithPrecedenceOne(t *testing.T) {
	var testStringSlice = `country IsEqualTo 'CANADA'
					   	   AND
					   	   country IsEqualTo 'CANADA'
					   	   OR
					       state IsEqualTo 'L'`
	result := gorules.EvaluateRules(testStringSlice, evaluatorTestData)
	fmt.Println(result)
	assert.False(t, false)
}

func TestMaths(t *testing.T) {
	result := gorules.EvaluateRules("'SOUTH' IsEqualTo |TAKETILL '-' 'SOTH-'|", evaluatorTestData)
	assert.True(t, result)
}
