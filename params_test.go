package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParams(t *testing.T) {
	params, err := ParamsFromMap(map[string]interface{}{
		"val1": "string",
		"val2": int(1),
		"val3": true,
		"val4": float64(43.444),
		"val5": FuncString("datetime.realtime()"),
	})
	require.Nil(t, err)
	require.NotNil(t, params)

	require.Contains(t, params.ToCypherMap(), "val4:43.444")
	require.Contains(t, params.ToCypherMap(), "val1:'string'")
	require.Contains(t, params.ToCypherMap(), "val2:1")
	require.Contains(t, params.ToCypherMap(), "val3:true")
	require.Contains(t, params.ToCypherMap(), "val5:datetime.realtime()")

	//error test
	err = params.Set("val", struct {
	}{})

	require.NotNil(t, err)
}
