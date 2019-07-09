package go_cypherdsl

import (
	"errors"
	"fmt"
	"reflect"
)

func DirectionPtr(d Direction) *Direction{
	return &d
}

func cypherizeInterface(i interface{}) (string, error){
	if i == nil{
		return "NULL", nil
	}

	//get interface kind
	t := reflect.TypeOf(i)
	k := t.Kind()

	if t == reflect.TypeOf(ParamString("")){
		s := i.(ParamString)
		return string(s), nil
	}

	//check string
	if k == reflect.String{
		return fmt.Sprintf("'%s'", i.(string)), nil
	}

	//check if primitive numeric type
	if k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 ||
		k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 ||
		k == reflect.Float32 || k == reflect.Float64{
		return fmt.Sprintf("%v", i), nil
	}

	//check bool
	if k == reflect.Bool{
		return fmt.Sprintf("%t", i.(bool)), nil
	}

	return "", errors.New("invalid type " + k.String())
}