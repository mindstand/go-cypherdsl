package go_cypherdsl

import (
	"errors"
	"fmt"
	"strings"
)

//v represents a vertex query
type V struct{
	//name of the vertex, omit if null
	Name string

	//type of edge, omit if null
	Type string

	//params for edge to map to, omit if null
	Params *Params
}

func (v *V) ToCypher() (string, error) {
	//if nothing is specified, its just an empty vertex
	if v.Name == "" && v.Type == "" && (v.Params == nil || v.Params.IsEmpty()){
		return "()", nil
	}

	str := "("

	//specify variable name if its there
	if v.Name != ""{
		str += v.Name
	}

	//specify type if its there
	if v.Type != ""{
		str += ":" + v.Type
	}

	//add params if its there
	if v.Params != nil{
		str += v.Params.ToCypherMap()
	}

	str += ")"

	return str, nil
}

// E represents an edge
type E struct {
	//direction of the edge, if null default to any
	Direction *Direction

	//variable name for constraint queries, omit if null
	Name string

	//names in the case that the edge is named or the query could be on multiple edges
	Types []string

	//min jumps to the next node, if null omit
	MinJumps int

	//max jumps to the next node, if null omit
	MaxJumps int

	//params for edges across individual jumps
	Params *Params
}

func (e *E) ToCypher() (string, error) {
	//check if the edge has anything specific
	if e.Name == "" && (e.Types == nil || len(e.Types) == 0) && e.MinJumps == 0 && e.MaxJumps == 0 && (e.Params == nil || e.Params.IsEmpty()){
		if e.Direction == nil{
			return "--", nil
		} else {
			if *e.Direction == Incoming {
				return "<--", nil
			} else if *e.Direction == Any {
				return "--", nil
			} else {
				return "-->", nil
			}
		}
	}

	core := "["

	if e.Name != ""{
		core += e.Name
	}

	if e.Types != nil && len(e.Types) != 0{
		if len(e.Types) == 1{
			core += ":" + e.Types[0]
		} else {
			q := ""
			for _, v := range e.Types {
				q += v + "|"
			}

			q = strings.TrimSuffix(q, "|")
			core += ":" + q
		}
	}

	if e.MinJumps != 0 && e.MaxJumps != 0{
		if (e.MinJumps >= e.MaxJumps) || e.MinJumps <= 0 || e.MaxJumps <= 0{
			return "", errors.New("min jumps can not be greater than or equal to max jumps, also can not be less than 0")
		}
		q := fmt.Sprintf("*%v..%v", e.MinJumps, e.MaxJumps)
		core += q
	} else if e.MinJumps != 0{
		if e.MinJumps <= 0{
			return "", errors.New("min jumps can not be less than 0")
		}
		q := fmt.Sprintf("*%v", e.MinJumps)
		core += q
	} else if e.MaxJumps != 0{
		if e.MaxJumps < 0{
			return "", errors.New("max jumps can not be less than 0")
		}
		q := fmt.Sprintf("*0..%v", e.MaxJumps)
		core += q
	}

	if e.Params != nil{
		core += e.Params.ToCypherMap()
	}

	core += "]"

	if e.Direction == nil{
		return fmt.Sprintf("-%s-", core), nil
	} else {
		if *e.Direction == Incoming{
			return fmt.Sprintf("<-%s-", core), nil
		} else {
			return fmt.Sprintf("-%s->", core), nil
		}
	}
}

type Direction int

const (
	Any Direction = 1
	Outgoing Direction = 2
	Incoming Direction = 3

)

type EdgeConfig struct {
	Type string
	StartNode int64
	EndNode int64
}