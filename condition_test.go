package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetStrVersion(t *testing.T)  {
	req := require.New(t)

	cypher, err := cypherizeInterface("test")
	req.Nil(err)
	req.EqualValues("'test'", cypher)

	cypher, err = cypherizeInterface(1)
	req.Nil(err)
	req.EqualValues("1", cypher)

	cypher, err = cypherizeInterface(true)
	req.Nil(err)
	req.EqualValues("true", cypher)

	_, err = cypherizeInterface(struct {

	}{})
	req.NotNil(err)
}

func TestNewCondition(t *testing.T) {
	//good condition checks
	req := require.New(t)

	//exists(type.value)
	cypher, err := NewCondition(&ConditionConfig{
		ConditionFunction: "exists",
		Name: "type",
		Field: "value",
	})
	req.Nil(err)
	req.EqualValues("exists(type.value)", cypher.ToString())

	//exists(toLower(type.value)
	cypher, err = NewCondition(&ConditionConfig{
		ConditionFunction: "exists",
		FieldManipulationFunction: "toLower",
		Name: "type",
		Field: "value",
	})
	req.Nil(err)
	req.EqualValues("exists(toLower(type.value))", cypher.ToString())

	//type.value >= 1
	cypher, err = NewCondition(&ConditionConfig{
		Name: "type",
		Field: "value",
		ConditionOperator: GreaterThanOrEqualToOperator,
		Check: 1,
	})
	req.Nil(err)
	req.EqualValues("type.value >= 1", cypher.ToString())

	//type.value = 'test'
	cypher, err = NewCondition(&ConditionConfig{
		Name: "type",
		Field: "value",
		ConditionOperator: EqualToOperator,
		Check: "test",
	})
	req.Nil(err)
	req.EqualValues("type.value = 'test'", cypher.ToString())

	//type.value in ['test','test2']
	cypher, err = NewCondition(&ConditionConfig{
		Name: "type",
		Field: "value",
		ConditionOperator: InOperator,
		CheckSlice: []interface{}{"test", "test2"},
	})
	req.Nil(err)
	req.EqualValues("type.value IN ['test','test2']", cypher.ToString())

	//type.value IS NULL
	cypher, err = NewCondition(&ConditionConfig{
		Name: "type",
		Field: "value",
		ConditionOperator: IsOperator,
		Check: nil,
	})
	req.Nil(err)
	req.EqualValues("type.value IS NULL", cypher.ToString())

	//error checks
	//nil check
	_, err = NewCondition(nil)
	req.NotNil(err)

	//checks label and field not defined
	_, err = NewCondition(&ConditionConfig{Name: "sfgs"})
	req.NotNil(err)

	//checks label and field both defined
	_, err = NewCondition(&ConditionConfig{
		Name:"adfa",
		Label: "dasdf",
		Field: "dafsd",
	})
	req.NotNil(err)

	//check operator and function not defined
	_, err = NewCondition(&ConditionConfig{
		Name:"adfa",
		Field: "dafsd",
	})
	req.NotNil(err)

	//check operator and function both defined
	_, err = NewCondition(&ConditionConfig{
		Name:"adfa",
		Field: "dafsd",
		ConditionOperator: "adfasd",
		ConditionFunction: "asdfasd",
	})
	req.NotNil(err)

	//check IN slice is nil
	_, err = NewCondition(&ConditionConfig{
		Name:"adfa",
		Field: "dafsd",
		ConditionOperator: InOperator,
		CheckSlice: nil,
	})
	req.NotNil(err)

	//check IN non slice check is not nil
	_, err = NewCondition(&ConditionConfig{
		Name:"adfa",
		Field: "dafsd",
		ConditionOperator: InOperator,
		Check: 45,
	})
	req.NotNil(err)

	//check IN slice is empty
	_, err = NewCondition(&ConditionConfig{
		Name:"adfa",
		Field: "dafsd",
		ConditionOperator: InOperator,
		CheckSlice: []interface{}{},
	})
	req.NotNil(err)

	//check invalid generic
	_, err = NewCondition(&ConditionConfig{
		Name:"adfa",
		Field: "dafsd",
		ConditionOperator: EqualToOperator,
		Check: struct {

		}{},
	})
	req.NotNil(err)
}

func TestConditionBuilder(t *testing.T){
	req := require.New(t)

	//(name.type = 1)
	cypher, err := C(&ConditionConfig{
		Name: "name",
		Field: "type",
		ConditionOperator: EqualToOperator,
		Check: 1,
	}).Build()
	req.Nil(err)
	req.EqualValues("name.type = 1", cypher.ToString())

	//name.type = 1 AND exists(name.type)
	cypher, err = C(&ConditionConfig{
		Name: "name",
		Field: "type",
		ConditionOperator: EqualToOperator,
		Check: 1,
	}).And(&ConditionConfig{
		Name: "name",
		Field: "type",
		ConditionFunction: "exists",
	}).Build()
	req.Nil(err)
	req.EqualValues("name.type = 1 AND exists(name.type)", cypher.ToString())

	//name.type = 1 AND (name.otherType >= 1 OR name.str STARTS WITH 'test')
	cypher, err = C(&ConditionConfig{
		Name: "name",
		Field: "type",
		ConditionOperator: EqualToOperator,
		Check: 1,
	}).AndNested(C(
		&ConditionConfig{
			Name: "name",
			Field: "otherType",
			ConditionOperator: GreaterThanOrEqualToOperator,
			Check: 1,
		}).Or(&ConditionConfig{
			Name:              "name",
			Field:             "str",
			ConditionOperator: StartsWithOperator,
			Check: "test",
		}).Build()).
	Build()
	req.Nil(err)
	req.EqualValues("name.type = 1 AND (name.otherType >= 1 OR name.str STARTS WITH 'test')", cypher.ToString())

	//todo fail tests
}