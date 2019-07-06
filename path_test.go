package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMatchBuilder(t *testing.T)  {
	req := require.New(t)

	//(a)
	cypher, err := NewPath().V(V{Name: StrPtr("a")}).ToCypher()
	req.Nil(err)
	req.EqualValues("(a)", cypher)

	//(a)--(b)
	cypher, err = NewPath().
		V(V{Name:StrPtr("a")}).
		E(E{}).
		V(V{Name: StrPtr("b")}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("(a)--(b)", cypher)

	//(a:type)--(b:type)
	cypher, err = NewPath().
		V(V{Name:StrPtr("a"), Type: StrPtr("type")}).
		E(E{}).
		V(V{Name: StrPtr("b"), Type: StrPtr("type")}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("(a:type)--(b:type)", cypher)

	//(a:type)-[e:type]->(b:type)
	cypher, err = NewPath().
		V(V{Name:StrPtr("a"), Type: StrPtr("type")}).
		E(E{Direction: DirectionPtr(Outgoing), Name: StrPtr("e"), Types:[]string{"type"}}).
		V(V{Name: StrPtr("b"), Type: StrPtr("type")}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("(a:type)-[e:type]->(b:type)", cypher)

	//p=(a:type),(b:type)--(c:type)
	cypher, err = NewPath().
		P().
		V(V{Name:StrPtr("a"), Type: StrPtr("type")}, V{Name: StrPtr("b"), Type: StrPtr("type")}).
		E(E{}).
		V(V{Name: StrPtr("c"), Type: StrPtr("type")}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("p=(a:type),(b:type)--(c:type)", cypher)

	//p=(a:type)-[e:type]->(b:type)
	cypher, err = NewPath().
		P().
		V(V{Name:StrPtr("a"), Type: StrPtr("type")}).
		E(E{Direction: DirectionPtr(Outgoing), Name: StrPtr("e"), Types:[]string{"type"}}).
		V(V{Name: StrPtr("b"), Type: StrPtr("type")}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("p=(a:type)-[e:type]->(b:type)", cypher)

	//(a)--(b)-[e:type*2..5 ]->(c:type{q:1})
	params, err := ParamsFromMap(map[string]interface{}{
		"q": 1,
	})
	req.Nil(err)

	cypher, err = NewPath().
		V(V{Name: StrPtr("a")}).
		E(E{}).
		V(V{Name:StrPtr("b")}).
		E(E{
			Name:StrPtr("e"),
			Types: []string{"type"},
			Direction: DirectionPtr(Outgoing),
			MaxJumps: IntPtr(5),
			MinJumps: IntPtr(2),
		}).
		V(V{
			Name: StrPtr("c"),
			Type:StrPtr("type"),
			Params: params,
		}).ToCypher()
	req.Nil(err)
	req.EqualValues("(a)--(b)-[e:type*2..5]->(c:type{q:1})", cypher)

	//p=(a)--()-[e:type*1..5]->(c:type {q:1})
	cypher, err = NewPath().
		P().
		V(V{Name: StrPtr("a")}).
		E(E{}).
		V(V{}).
		E(E{
			Name:StrPtr("e"),
			Types: []string{"type"},
			Direction: DirectionPtr(Outgoing),
			MaxJumps: IntPtr(5),
		}).
		V(V{
			Name: StrPtr("c"),
			Type:StrPtr("type"),
			Params: params,
		}).ToCypher()
	req.Nil(err)
	req.EqualValues("p=(a)--()-[e:type*1..5]->(c:type{q:1})", cypher)
}