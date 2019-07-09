package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFunctionConfig_ToString(t *testing.T) {
	name := FunctionConfig{
		Name:"test",
	}

	oneP := FunctionConfig{
		Name: "test",
		Params: []interface{}{1},
	}

	twoP := FunctionConfig{
		Name: "test",
		Params: []interface{}{1, ParamString("test")},
	}

	errC := FunctionConfig{}

	req := require.New(t)
	var cypher string
	var err error

	//name only
	cypher, err = name.ToString()
	req.Nil(err)
	req.EqualValues("test()", cypher)

	//name with one param
	cypher, err = oneP.ToString()
	req.Nil(err)
	req.EqualValues("test(1)", cypher)

	//name with 2 params
	cypher, err = twoP.ToString()
	req.Nil(err)
	req.EqualValues("test(1,test)", cypher)

	//nothing configured
	_, err = errC.ToString()
	req.NotNil(err)
}