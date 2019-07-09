package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRemoveConfig_ToString(t *testing.T) {
	req := require.New(t)
	var err error
	var cypher string

	//test field
	t1 := RemoveConfig{
		Name: "n",
		Field: "f",
	}
	cypher, err = t1.ToString()
	req.Nil(err)
	req.EqualValues("n.f", cypher)

	//test label
	t2 := RemoveConfig{
		Name: "n",
		Labels: []string{"f"},
	}
	cypher, err = t2.ToString()
	req.Nil(err)
	req.EqualValues("n:f", cypher)

	//test many labels
	t3 := RemoveConfig{
		Name: "n",
		Labels: []string{"1", "2"},
	}
	cypher, err = t3.ToString()
	req.Nil(err)
	req.EqualValues("n:1:2", cypher)

	//error empty name
	e1 := RemoveConfig{}
	_, err = e1.ToString()
	req.NotNil(err)

	//error label and field
	e2 := RemoveConfig{
		Labels: []string{"dfd"},
		Name: "dafd",
		Field: "dasd",
	}
	_, err = e2.ToString()
	req.NotNil(err)

	//error error no label or field
	e3 := RemoveConfig{
		Name: "dfad",
	}
	_, err = e3.ToString()
	req.NotNil(err)
}
