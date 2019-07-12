package go_cypherdsl

import (
	"errors"
	"fmt"
	neo "github.com/johnnadratowski/golang-neo4j-bolt-driver"
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

func RowsToStringArray(rows neo.Rows) ([]string, error){
	data, _, err := rows.All()
	if err != nil{
		return nil, err
	}

	//check to make sure its not empty
	if data == nil || len(data) == 0 || len(data[0]) == 0{
		//todo standard "not found" error
		return nil, errors.New("data can not be empty")
	}


	_, ok := data[0][0].(string)
	if !ok{
		return nil, errors.New("does not contain array of strings")
	}

	toReturn := make([]string, len(data))
	for i, v := range data{
		if len(v) == 0{
			return nil, errors.New("index %v is empty")
		}
		toReturn[i], ok = v[0].(string)
	}

	return toReturn, nil
}

func RowsTo2dStringArray(rows neo.Rows) ([][]string, error){
	data, _, err := rows.All()
	if err != nil{
		return nil, err
	}

	if len(data) != 0 && len(data[0]) != 0{
		toReturn := make([][]string, len(data))

		var ok bool

		for i, v := range data{

			toReturn[i] = make([]string, len(v))

			for j, v1 := range v{
				toReturn[i][j], ok = v1.(string)
				if !ok{
					return nil, errors.New("failed to cast value to string")
				}
			}
		}
		return toReturn, rows.Close()
	} else {
		return nil, errors.New("unknown response type")
	}
}

func RowsTo2DInterfaceArray(rows neo.Rows) ([][]interface{}, error){
	data, _, err := rows.All()
	return data, err
}