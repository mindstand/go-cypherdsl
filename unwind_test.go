package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnwindConfig_ToString(t *testing.T) {
	req := require.New(t)

	//{a,b} AS A
	conf := UnwindConfig{
		Slice: []interface{}{
			ParamString("a"),
			ParamString("b"),
		},
		As: "A",
	}
	cypher, err := conf.ToString()
	req.Nil(err)
	req.EqualValues("[a,b] AS A", cypher)

	//{{a}, a} AS B
	conf = UnwindConfig{
		Slice: []interface{}{
			[]interface{}{ParamString("a")},
			ParamString("b"),
		},
		As: "A",
	}
	cypher, err = conf.ToString()
	req.Nil(err)
	req.EqualValues("[[a],b] AS A", cypher)

	//fail
	conf = UnwindConfig{}
	_, err = conf.ToString()
	req.NotNil(err)
}