package go_cypherdsl

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Params struct{
	params map[string]string
}

func ParamsFromMap(m map[string]interface{}) (*Params, error){
	if m == nil || len(m) == 0{
		return nil, errors.New("map can not be empty or nil")
	}

	p := &Params{}

	for k, v := range m{
		err := p.Set(k, v)
		if err != nil{
			return nil, err
		}
	}

	return p, nil
}

func (p *Params) IsEmpty() bool{
	return p.params == nil || len(p.params) == 0
}

func (p *Params) Set(key string, value interface{}) error{
	if p.params == nil{
		p.params = map[string]string{}
	}

	//get interface kind
	k := reflect.TypeOf(value).Kind()

	//check string
	if k == reflect.String{
		p.params[key] = fmt.Sprintf("%s:\"%s\"", key, value.(string))
		return nil
	}

	//check if primitive numeric type
	if k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 ||
		k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 ||
		k == reflect.Float32 || k == reflect.Float64{
		p.params[key] = fmt.Sprintf("%s:%v", key, value)
		return nil
	}

	//check bool
	if k == reflect.Bool{
		p.params[key] = fmt.Sprintf("%s:%t", key, value.(bool))
		return nil
	}

	return fmt.Errorf("unknown type %s", k.String())
}

func (p *Params) ToCypherMap() string{

	if p.params == nil || len(p.params) == 0{
		return "{}"
	}

	q := ""

	for _, v := range p.params{
		q += fmt.Sprintf("%s,", v)
	}

	return "{" + strings.TrimSuffix(q, ",") + "}"
}

