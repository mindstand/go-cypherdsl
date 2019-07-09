package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryBuilder(t *testing.T){
	req := require.New(t)

	//match (n) return n
	cypher, err := QB().Match(Path().V(V{Name:"n"}).Done()).Return(ReturnPart{Name: "n"}).ToCypher()
	req.Nil(err)
	req.EqualValues("MATCH (n) RETURN n", cypher)

}
