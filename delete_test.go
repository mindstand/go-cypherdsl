package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteToString(t *testing.T) {
	req := require.New(t)
	var err error
	var cypher string

	//delete one param
	cypher, err = deleteToString(false, "n")
	req.Nil(err)
	req.EqualValues("DELETE n", cypher)

	//detach delete one param
	cypher, err = deleteToString(true, "n")
	req.Nil(err)
	req.EqualValues("DETACH DELETE n", cypher)

	//delete one param
	cypher, err = deleteToString(false, "n", "b", "C")
	req.Nil(err)
	req.EqualValues("DELETE n, b, C", cypher)

	//delete no params (error)
	_, err = deleteToString(false)
	req.NotNil(err)
}
