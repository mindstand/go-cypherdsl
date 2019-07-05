package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)
func TestV_ToCypher(t *testing.T) {
	req := require.New(t)

	empty := V{}
	withVarName := V{Name: StrPtr("test_str")}
	withVarNameAndType := V{Name: StrPtr("test_str"), Type:StrPtr("test_type")}

	pms, err := ParamsFromMap(map[string]interface{}{
		"test_var": "test",
	})
	req.Nil(err)

	withAll := V{
		Name:   StrPtr("test_str"),
		Type:   StrPtr("test_type"),
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
	req.EqualValues(cypher, "(test_str:test_type{test_var:\"test\"})")
}

func TestE_ToCypher(t *testing.T) {
	req := require.New(t)

	//not specified
	nothingSpecified := E{}
	onlyDirectionIncoming := E{Direction: DirectionPtr(Incoming)}
	onlyDirectionOutgoing := E{Direction: DirectionPtr(Outgoing)}
	onlyDirectionAny := E{Direction: DirectionPtr(Any)}

	//var types
	varNameOnly := E{Name: StrPtr("var")}
	varNameWithOneType := E{Name: StrPtr("var"), Types: []string{"test"}}
	varWithManyTypes := E{Name: StrPtr("var"), Types: []string{"test1", "test2"}}

	//jumps
	jumpsMinMax := E{MinJumps: IntPtr(2), MaxJumps: IntPtr(5)}
	jumpsMin := E{MinJumps: IntPtr(2)}
	jumpsMax := E{MaxJumps: IntPtr(3)}

	pms, err := ParamsFromMap(map[string]interface{}{
		"test": 1,
	})
	req.Nil(err)

	//params
	params := E{
		Params: pms,
	}

	//everything
	all := E{
		Direction: DirectionPtr(Outgoing),
		MaxJumps:  IntPtr(20),
		MinJumps:  IntPtr(4),
		Params:    pms,
		Types:     []string{"one", "two"},
		Name:      StrPtr("q"),
	}

	//error test
	errorJumps := E{MinJumps: IntPtr(5), MaxJumps: IntPtr(2)}
	errorNegMaxJumps := E{MinJumps: IntPtr(-1)}
	errorNegMinJumps := E{MaxJumps: IntPtr(-1)}



	cypher, err := nothingSpecified.ToCypher()
	req.Nil(err)
	req.EqualValues("--", cypher)

	cypher, err = onlyDirectionIncoming.ToCypher()
	req.Nil(err)
	req.EqualValues("<--", cypher)

	cypher, err = onlyDirectionOutgoing.ToCypher()
	req.Nil(err)
	req.EqualValues("-->", cypher)

	cypher, err = onlyDirectionAny.ToCypher()
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
	req.EqualValues("-[*1..3]-", cypher)

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