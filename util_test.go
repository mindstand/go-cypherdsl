package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCypherizeInterface(t *testing.T) {
	req := require.New(t)

	str, err := cypherizeInterface(nil)
	req.Nil(err)
	req.EqualValues("NULL", str)

	str, err = cypherizeInterface(ParamString("test"))
	req.Nil(err)
	req.EqualValues("test", str)

	str, err = cypherizeInterface("test")
	req.Nil(err)
	req.EqualValues("'test'", str)

	str, err = cypherizeInterface(1)
	req.Nil(err)
	req.EqualValues("1", str)

	str, err = cypherizeInterface(true)
	req.Nil(err)
	req.EqualValues("true", str)

	_, err = cypherizeInterface(struct {
	}{})
	req.NotNil(err)
}
