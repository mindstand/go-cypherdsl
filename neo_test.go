package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type TestSerialize struct{
	Name string `json:"name"`
	Age int `json:"age"`
}

func TestIndexs(t *testing.T){
	if !testing.Short(){
		t.SkipNow()
		return
	}

	err := Init(&ConnectionConfig{
		Username: "neo4j",
		Password: "password",
		Host: "0.0.0.0",
		Port: 7687,
		PoolSize: 15,
	})
	require.Nil(t, err)

	sess := NewSession()

	rows, err := sess.QueryReadOnly().Cypher("CALL db.constraints()").Query(nil)
	require.Nil(t, err)

	vals, err := RowsToStringArray(rows)
	require.Nil(t, err)
	for _, v := range vals{
		log.Println(v)
	}

}

//this  is purely to demonstrate usage
func TestNeo(t *testing.T){
	//comment out to actually run
	if !testing.Short(){
		t.SkipNow()
		return
	}

	req := require.New(t)

	err := Init(&ConnectionConfig{
		Username: "neo4j",
		Password: "password",
		Host: "mindstand.com",
		Port: 7687,
		PoolSize: 15,
	})

	sess := NewSession()
	defer sess.Close()
	//err = sess.Begin()
	//req.Nil(err)
	//
	//ericParams, err := ParamsFromMap(map[string]interface{}{
	//	"name": "Eric",
	//	"age": 21,
	//})
	//req.Nil(err)
	//
	//nikitaParams, err := ParamsFromMap(map[string]interface{}{
	//	"name": "Nikita",
	//	"age": 21,
	//})
	//req.Nil(err)
	//
	//path := Path().P().V(V{Type:"TEST", Params:ericParams}).E(E{Types: []string{"CONN"}, Direction:DirectionPtr(DirectionOutgoing)}).V(V{Type: "TEST", Params: nikitaParams}).Build()
	//
	//res, err := sess.Query().Create(NewNode(path)).Return(ReturnPart{Name:"p"}).Exec(nil)
	//req.Nil(err)
	//req.NotNil(res)
	//
	//err = sess.Commit()
	//req.Nil(err)
	//
	//t.Log(res.RowsAffected())
	//t.Log(res.LastInsertId())
	//t.Log(res.Metadata())

	rows, err := sess.QueryReadOnly().
		Match(Path().V(V{Name: "n", Type:"OrganizationNode"}).Build()).
		With(&WithConfig{
			Parts: []WithPart{
				{Name: "n"},
			},
		}).
		Match(Path().P().V(V{Name: "n"}).E(E{Name: "e", MaxJumps:2}).V(V{Name: "c"}).Build()).Return(true, ReturnPart{Name: "p"}).Query(nil)
	req.Nil(err)

	arr, meta, err := rows.All()
	req.Nil(err)

	t.Log(arr)
	t.Log(meta)

	//notes, we will cast the first output from rows.all to graph.Path then use the ogm to convert that to structs

	t.Log("done")
}