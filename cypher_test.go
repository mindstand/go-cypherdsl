package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryBuilder(t *testing.T){
	req := require.New(t)

	//match (n) return n
	cypher, err := QB().Match(Path().V(V{Name:"n"}).Build()).Return(false, ReturnPart{Name: "n"}).ToCypher()
	req.Nil(err)
	req.EqualValues("MATCH (n) RETURN n", cypher)

	//MATCH (n) WHERE n.age = 21 RETURN n
	cypher, err = QB().
		Match(Path().V(V{Name: "n"}).Build()).Where(C(&ConditionConfig{
			Name: "n",
			Field: "age",
			ConditionOperator: EqualToOperator,
			Check: 21,
		})).
		Return(false, ReturnPart{Name: "n"}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("MATCH (n) WHERE n.age = 21 RETURN n", cypher)

	//MATCH (n) RETURN n ORDER BY n.age DESC LIMIT 5
	cypher, err = QB().
		Match(Path().V(V{Name: "n"}).Build()).
		Return(false, ReturnPart{Name: "n"}).
		OrderBy(OrderByConfig{Name: "n", Member: "age", Desc: true}).
		Limit(5).
		ToCypher()
	req.Nil(err)
	req.EqualValues("MATCH (n) RETURN n ORDER BY n.age DESC LIMIT 5", cypher)

	//CREATE (n:Type{name:'Eric',age:21}) RETURN n
	params, err := ParamsFromMap(map[string]interface{}{
		"name": "Eric",
		"age": 21,
	})
	req.Nil(err)
	req.NotNil(params)

	cypher, err = QB().
		Create(NewNode(Path().V(V{
			Name: "n",
			Type: "Type",
			Params: params,
		}).Build())).
		Return(false, ReturnPart{Name:"n"}).
		ToCypher()
	req.Nil(err)
	req.Contains(cypher, "CREATE (n:Type{")
	req.Contains(cypher, "}) RETURN n")
	req.Contains(cypher, "name:'Eric'")
	req.Contains(cypher, "age:21")

	//MATCH (n) DETACH DELETE n
	cypher, err = QB().
		Match(Path().V(V{Name:"n"}).Build()).
		Delete(true, "n").
		ToCypher()
	req.Nil(err)
	req.EqualValues("MATCH (n) DETACH DELETE n", cypher)

	//MATCH (n) SET n.name=5
	cypher, err = QB().
		Match(Path().V(V{Name: "n"}).Build()).
		Set(SetConfig{
			Name: "n",
			Member: "name",
			Operation: SetEqualTo,
			Target: 5,
		}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("MATCH (n) SET n.name = 5", cypher)


	path, err := Path().V(V{Name: "city", Type:"City"}).ToCypher()
	req.Nil(err)

	//MERGE (city:City) REMOVE city.name
	cypher, err = QB().
		Merge(&MergeConfig{
			Path: path,
		}).
		Remove(RemoveConfig{
			Name: "city",
			Field: "name",
		}).
		ToCypher()
	req.Nil(err)
	req.EqualValues("MERGE (city:City) REMOVE city.name", cypher)
}
