package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetConfig_ToString(t *testing.T) {
	req := require.New(t)
	var err error
	var cypher string

	//check conditional
	t1 := &SetConfig{
		Operation: SetEqualTo,
		Target:    1,
		Member:    "t",
		Name:      "f",
		Condition: C(&ConditionConfig{
			Name:              "t",
			Field:             "s",
			ConditionOperator: EqualToOperator,
			Check:             4,
		}),
	}
	cypher, err = t1.ToString()
	req.Nil(err)
	req.EqualValues("(CASE WHEN t.s = 4 THEN f END).t = 1", cypher)

	//check labels
	t2 := &SetConfig{
		Name:  "t",
		Label: []string{"l1", "l2"},
	}
	cypher, err = t2.ToString()
	req.Nil(err)
	req.EqualValues("t:l1:l2", cypher)

	//check mutate
	param, err := ParamsFromMap(map[string]interface{}{
		"test": 1,
	})
	req.Nil(err)
	req.NotNil(param)

	t3 := &SetConfig{
		Name:      "t",
		Operation: SetMutate,
		TargetMap: param,
	}
	cypher, err = t3.ToString()
	req.Nil(err)
	req.EqualValues("t += {test:1}", cypher)

	//check set node to map
	t4 := &SetConfig{
		Name:      "t",
		Operation: SetEqualTo,
		TargetMap: param,
	}
	cypher, err = t4.ToString()
	req.Nil(err)
	req.EqualValues("t = {test:1}", cypher)

	//check set node member
	t5 := &SetConfig{
		Operation: SetEqualTo,
		Target:    1,
		Member:    "t",
		Name:      "f",
	}
	cypher, err = t5.ToString()
	req.Nil(err)
	req.EqualValues("f.t = 1", cypher)

	//check set target to param string
	t6 := &SetConfig{
		Operation: SetEqualTo,
		Target:    ParamString("$props"),
		Name:      "f",
	}
	cypher, err = t6.ToString()
	req.Nil(err)
	req.EqualValues("f = $props", cypher)

	//error - name not defined
	e1 := &SetConfig{}
	_, err = e1.ToString()
	req.NotNil(err)

	//error - operation not defined
	e2 := &SetConfig{Name: "asdfasd"}
	_, err = e2.ToString()
	req.NotNil(err)

	//error - no targets defined
	e3 := &SetConfig{Name: "asdfasd", Operation: SetEqualTo}
	_, err = e3.ToString()
	req.NotNil(err)

	//error - 2 targets defined
	e4 := &SetConfig{
		Name:      "asdfasd",
		Operation: SetEqualTo,
		Target:    1,
		TargetMap: param,
	}
	_, err = e4.ToString()
	req.NotNil(err)

	//error - 3 targets defined
	e5 := &SetConfig{
		Name:           "asdfasd",
		Operation:      SetEqualTo,
		Target:         1,
		TargetMap:      param,
		TargetFunction: &FunctionConfig{Name: "s"},
	}
	_, err = e5.ToString()
	req.NotNil(err)

	//error - trying to set labels on conditional
	e6 := &SetConfig{
		Condition: C(&ConditionConfig{
			Name:              "t",
			Field:             "s",
			ConditionOperator: EqualToOperator,
			Check:             4,
		}),
		Name:  "check",
		Label: []string{"asdfasd"},
	}
	_, err = e6.ToString()
	req.NotNil(err)

	//error - trying to set labels on mutate operation
	e7 := &SetConfig{
		Name:      "t",
		Operation: SetMutate,
		Label:     []string{"dasdf"},
		Target:    1,
	}
	_, err = e7.ToString()
	req.NotNil(err)

	//error - trying to use member on mutate operation
	e8 := &SetConfig{
		Name:      "t",
		Operation: SetMutate,
		Member:    "asdfasd",
		Target:    1,
	}
	_, err = e8.ToString()
	req.NotNil(err)

	//error - not using target map on mutate operation
	e9 := &SetConfig{
		Name:      "t",
		Operation: SetMutate,
		Target:    1,
	}
	_, err = e9.ToString()
	req.NotNil(err)

	//error - trying to set node equal to plain target
	e10 := &SetConfig{
		Name:      "t",
		Operation: SetEqualTo,
		Target:    1,
	}
	_, err = e10.ToString()
	req.NotNil(err)
}
