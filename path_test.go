package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMatchBuilder(t *testing.T) {
	req := require.New(t)

	//(a)
	cypher, err := NewPath().V(V{Name: "a"}).ToCypher()
	req.Nil(err)
	req.EqualValues("(a)", cypher)

	//(a)--(b)
	cypher, err = NewPath().
		V(V{Name: "a"}).
		E(E{Direction: DirectionNone}).
		V(V{Name: "b"}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("(a)--(b)", cypher)

	//(a:type)--(b:type)
	cypher, err = NewPath().
		V(V{Name: "a", Type: "type"}).
		E(E{Direction: DirectionNone}).
		V(V{Name: "b", Type: "type"}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("(a:type)--(b:type)", cypher)

	//(a:type)-[e:type]->(b:type)
	cypher, err = NewPath().
		V(V{Name: "a", Type: "type"}).
		E(E{Direction: DirectionOutgoing, Name: "e", Types: []string{"type"}}).
		V(V{Name: "b", Type: "type"}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("(a:type)-[e:type]->(b:type)", cypher)

	//p=(a:type),(b:type)--(c:type)
	cypher, err = NewPath().
		P().
		V(V{Name: "a", Type: "type"}, V{Name: "b", Type: "type"}).
		E(E{Direction: DirectionNone}).
		V(V{Name: "c", Type: "type"}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("p=(a:type),(b:type)--(c:type)", cypher)

	//p=(a:type)-[e:type]->(b:type)
	cypher, err = NewPath().
		P().
		V(V{Name: "a", Type: "type"}).
		E(E{Direction: DirectionOutgoing, Name: "e", Types: []string{"type"}}).
		V(V{Name: "b", Type: "type"}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("p=(a:type)-[e:type]->(b:type)", cypher)

	//(a)--(b)-[e:type*2..5 ]->(c:type{q:1})
	params, err := ParamsFromMap(map[string]interface{}{
		"q": 1,
	})
	req.Nil(err)

	cypher, err = NewPath().
		V(V{Name: "a"}).
		E(E{Direction: DirectionNone}).
		V(V{Name: "b"}).
		E(E{
			Name:      "e",
			Types:     []string{"type"},
			Direction: DirectionOutgoing,
			MaxJumps:  5,
			MinJumps:  2,
		}).
		V(V{
			Name:   "c",
			Type:   "type",
			Params: params,
		}).ToCypher()
	req.Nil(err)
	req.EqualValues("(a)--(b)-[e:type*2..5]->(c:type{q:1})", cypher)

	//p=(a)--()-[e:type*0..5]->(c:type {q:1})
	cypher, err = NewPath().
		P().
		V(V{Name: "a"}).
		E(E{Direction: DirectionNone}).
		V(V{}).
		E(E{
			Name:      "e",
			Types:     []string{"type"},
			Direction: DirectionOutgoing,
			MaxJumps:  5,
		}).
		V(V{
			Name:   "c",
			Type:   "type",
			Params: params,
		}).ToCypher()
	req.Nil(err)
	req.EqualValues("p=(a)--()-[e:type*0..5]->(c:type{q:1})", cypher)
}
