package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOrderByConfig_ToString(t *testing.T) {
	req := require.New(t)
	var err error
	var cypher string

	//name not defined
	t1 := OrderByConfig{
		Member: "da",
	}
	_, err = t1.ToString()
	req.NotNil(err)

	//member not defined
	t2 := OrderByConfig{
		Name: "da",
	}
	_, err = t2.ToString()
	req.NotNil(err)

	//both member and name not defined
	t3 := OrderByConfig{}
	_, err = t3.ToString()
	req.NotNil(err)

	//proper
	t4 := OrderByConfig{
		Name: "n",
		Member: "m",
	}
	cypher, err = t4.ToString()
	req.Nil(err)
	req.EqualValues("n.m", cypher)

	//proper
	t5 := OrderByConfig{
		Name: "n",
		Member: "m",
		Desc: true,
	}
	cypher, err = t5.ToString()
	req.Nil(err)
	req.EqualValues("n.m DESC", cypher)
}