package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestV_ToCypher(t *testing.T) {
	req := require.New(t)

	empty := V{}
	withVarName := V{Name: "test_str"}
	withVarNameAndType := V{Name: "test_str", Type: "test_type"}

	pms, err := ParamsFromMap(map[string]interface{}{
		"test_var": "test",
	})
	req.Nil(err)

	withAll := V{
		Name:   "test_str",
		Type:   "test_type",
		Params: pms,
	}

	//test empty
	cypher, err := empty.ToCypher()
	req.Nil(err)
	req.EqualValues(cypher, "()")

	//test with name
	cypher, err = withVarName.ToCypher()
	req.Nil(err)
	req.EqualValues(cypher, "(test_str)")

	//test with name and type
	cypher, err = withVarNameAndType.ToCypher()
	req.Nil(err)
	req.EqualValues(cypher, "(test_str:test_type)")

	//test with name
	cypher, err = withAll.ToCypher()
	req.Nil(err)
	req.EqualValues(cypher, "(test_str:test_type{test_var:'test'})")
}

func TestE_ToCypher(t *testing.T) {
	req := require.New(t)

	//not specified
	nothingSpecified := E{}
	onlyDirectionIncoming := E{Direction: DirectionIncoming}
	onlyDirectionOutgoing := E{Direction: DirectionOutgoing}
	onlyDirectionNone := E{Direction: DirectionNone}

	//var types
	varNameOnly := E{Name: "var", Direction: DirectionNone}
	varNameWithOneType := E{Name: "var", Types: []string{"test"}, Direction: DirectionNone}
	varWithManyTypes := E{Name: "var", Types: []string{"test1", "test2"}, Direction: DirectionNone}

	//jumps
	jumpsMinMax := E{MinJumps: 2, MaxJumps: 5, Direction: DirectionNone}
	jumpsMin := E{MinJumps: 2, Direction: DirectionNone}
	jumpsMax := E{MaxJumps: 3, Direction: DirectionNone}

	pms, err := ParamsFromMap(map[string]interface{}{
		"test": 1,
	})
	req.Nil(err)

	//params
	params := E{
		Params:    pms,
		Direction: DirectionNone,
	}

	//everything
	all := E{
		Direction: DirectionOutgoing,
		MaxJumps:  20,
		MinJumps:  4,
		Params:    pms,
		Types:     []string{"one", "two"},
		Name:      "q",
	}

	//error test
	errorJumps := E{MinJumps: 5, MaxJumps: 2}
	errorNegMaxJumps := E{MinJumps: -1}
	errorNegMinJumps := E{MaxJumps: -1}

	cypher, err := nothingSpecified.ToCypher()
	req.Nil(err)
	req.EqualValues("-->", cypher)

	cypher, err = onlyDirectionIncoming.ToCypher()
	req.Nil(err)
	req.EqualValues("<--", cypher)

	cypher, err = onlyDirectionOutgoing.ToCypher()
	req.Nil(err)
	req.EqualValues("-->", cypher)

	cypher, err = onlyDirectionNone.ToCypher()
	req.Nil(err)
	req.EqualValues("--", cypher)

	cypher, err = varNameOnly.ToCypher()
	req.Nil(err)
	req.EqualValues("-[var]-", cypher)

	cypher, err = varNameWithOneType.ToCypher()
	req.Nil(err)
	req.EqualValues("-[var:test]-", cypher)

	cypher, err = varWithManyTypes.ToCypher()
	req.Nil(err)
	req.EqualValues("-[var:test1|test2]-", cypher)

	cypher, err = jumpsMinMax.ToCypher()
	req.Nil(err)
	req.EqualValues("-[*2..5]-", cypher)

	cypher, err = jumpsMin.ToCypher()
	req.Nil(err)
	req.EqualValues("-[*2]-", cypher)

	cypher, err = jumpsMax.ToCypher()
	req.Nil(err)
	req.EqualValues("-[*0..3]-", cypher)

	cypher, err = params.ToCypher()
	req.Nil(err)
	req.EqualValues("-[{test:1}]-", cypher)

	cypher, err = all.ToCypher()
	req.Nil(err)
	req.EqualValues("-[q:one|two*4..20{test:1}]->", cypher)

	cypher, err = errorJumps.ToCypher()
	req.NotNil(err)
	req.EqualValues("", cypher)

	cypher, err = errorNegMaxJumps.ToCypher()
	req.NotNil(err)
	req.EqualValues("", cypher)

	cypher, err = errorNegMinJumps.ToCypher()
	req.NotNil(err)
	req.EqualValues("", cypher)
}
