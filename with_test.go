package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithConfig_ToString(t *testing.T) {
	req := require.New(t)

	//one
	conf := WithConfig{
		Parts: []WithPart{
			{
				Name: "test",
			},
		},
	}
	cypher, err := conf.ToString()
	req.Nil(err)
	req.EqualValues("test", cypher)

	//two
	conf = WithConfig{
		Parts: []WithPart{
			{
				Name: "test",
			},
			{
				Name: "t",
				As:   "woo",
			},
		},
	}
	cypher, err = conf.ToString()
	req.Nil(err)
	req.EqualValues("test, t AS woo", cypher)

	//empty
	conf = WithConfig{}
	_, err = conf.ToString()
	req.NotNil(err)
}

func TestWithPart_ToString(t *testing.T) {
	req := require.New(t)

	//test function
	part := WithPart{
		Function: &FunctionConfig{
			Name: "test",
		},
	}

	cypher, err := part.ToString()
	req.Nil(err)
	req.EqualValues("test()", cypher)

	//test function as
	part = WithPart{
		Function: &FunctionConfig{
			Name: "test",
		},
		As: "TEST",
	}

	cypher, err = part.ToString()
	req.Nil(err)
	req.EqualValues("test() AS TEST", cypher)

	//test name
	part = WithPart{
		Name: "test",
	}

	cypher, err = part.ToString()
	req.Nil(err)
	req.EqualValues("test", cypher)

	//test name field
	part = WithPart{
		Name:  "test",
		Field: "t",
	}

	cypher, err = part.ToString()
	req.Nil(err)
	req.EqualValues("test.t", cypher)

	//test name field as
	part = WithPart{
		Name:  "test",
		Field: "t",
		As:    "TEST",
	}

	cypher, err = part.ToString()
	req.Nil(err)
	req.EqualValues("test.t AS TEST", cypher)

	//test nothing defined
	part = WithPart{}
	_, err = part.ToString()
	req.NotNil(err)

	//test function with name (error)
	part = WithPart{
		Function: &FunctionConfig{
			Name: "t",
		},
		Name: "test",
	}
	_, err = part.ToString()
	req.NotNil(err)
}
