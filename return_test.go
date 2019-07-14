package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReturnPart_ToString(t *testing.T) {
	req := require.New(t)

	literal := ReturnPart{Literal: "test"}
	boolExpr := ReturnPart{BooleanExpression: "test"}
	path := ReturnPart{Path: "test-test"}
	nameOnly := ReturnPart{Name: "name"}
	nameType := ReturnPart{Name: "name", Type: "whatever"}
	nameTypeAlias := ReturnPart{Name: "name", Type: "whatever", Alias: "test"}
	empty := ReturnPart{}

	var err error
	var cypher string

	//test literal
	cypher, err = literal.ToString()
	req.Nil(err)
	req.EqualValues("'test'", cypher)

	//test bool expression
	cypher, err = boolExpr.ToString()
	req.Nil(err)
	req.EqualValues("test", cypher)

	//test path
	cypher, err = path.ToString()
	req.Nil(err)
	req.EqualValues("test-test", cypher)

	//test name only
	cypher, err = nameOnly.ToString()
	req.Nil(err)
	req.EqualValues("name", cypher)

	//test name.whatever
	cypher, err = nameType.ToString()
	req.Nil(err)
	req.EqualValues("name.whatever", cypher)

	//test name.whatever as test
	cypher, err = nameTypeAlias.ToString()
	req.Nil(err)
	req.EqualValues("name.whatever AS test", cypher)

	//test nothing defined
	cypher, err = empty.ToString()
	req.NotNil(err)
	req.EqualValues("", cypher)
}

func TestNewReturnClause(t *testing.T) {
	req := require.New(t)

	//empty test
	cypher, err := NewReturnClause()
	req.NotNil(err)
	req.EqualValues("", cypher)

	//one part
	cypher, err = NewReturnClause(ReturnPart{
		Name: "t",
	})
	req.Nil(err)
	req.EqualValues("RETURN t", cypher.ToString())

	//two parts
	cypher, err = NewReturnClause(ReturnPart{
		Name: "t1",
	}, ReturnPart{
		Name: "t2",
	})
	req.Nil(err)
	req.EqualValues("RETURN t1, t2", cypher.ToString())
}