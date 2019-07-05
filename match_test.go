package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMatchBuilder(t *testing.T)  {
	req := require.New(t)

	//match (a)
	cypher, err := Match().V(V{Name: StrPtr("a")}).ToCypher()
	req.Nil(err)
	req.EqualValues("match (a)", cypher)

	//match (a)--(b)
	cypher, err = Match().
		V(V{Name:StrPtr("a")}).
		E(E{}).
		V(V{Name: StrPtr("b")}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("match (a)--(b)", cypher)

	//match (a:type)--(b:type)
	cypher, err = Match().
		V(V{Name:StrPtr("a"), Type: StrPtr("type")}).
		E(E{}).
		V(V{Name: StrPtr("b"), Type: StrPtr("type")}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("match (a:type)--(b:type)", cypher)

	//match (a:type)-[e:type]->(b:type)
	cypher, err = Match().
		V(V{Name:StrPtr("a"), Type: StrPtr("type")}).
		E(E{Direction: DirectionPtr(Outgoing), Name: StrPtr("e"), Types:[]string{"type"}}).
		V(V{Name: StrPtr("b"), Type: StrPtr("type")}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("match (a:type)-[e:type]->(b:type)", cypher)

	//match p=(a:type),(b:type)--(c:type)
	cypher, err = Match().
		P().
		V(V{Name:StrPtr("a"), Type: StrPtr("type")}, V{Name: StrPtr("b"), Type: StrPtr("type")}).
		E(E{}).
		V(V{Name: StrPtr("c"), Type: StrPtr("type")}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("match p=(a:type),(b:type)--(c:type)", cypher)

	//match p=(a:type)-[e:type]->(b:type)
	cypher, err = Match().
		P().
		V(V{Name:StrPtr("a"), Type: StrPtr("type")}).
		E(E{Direction: DirectionPtr(Outgoing), Name: StrPtr("e"), Types:[]string{"type"}}).
		V(V{Name: StrPtr("b"), Type: StrPtr("type")}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("match p=(a:type)-[e:type]->(b:type)", cypher)

	//match (a)--(b)-[e:type*2..5 ]->(c:type{q:1})
	params, err := ParamsFromMap(map[string]interface{}{
		"q": 1,
	})
	req.Nil(err)

	cypher, err = Match().
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
	req.EqualValues("match (a)--(b)-[e:type*2..5]->(c:type{q:1})", cypher)

	//match p=(a)--()-[e:type*1..5]->(c:type {q:1})
	cypher, err = Match().
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
	req.EqualValues("match p=(a)--()-[e:type*1..5]->(c:type{q:1})", cypher)
}