package go_cypherdsl

import (
	"encoding/json"
	"fmt"
	"strings"
)

//v represents a vertex query
type V struct{
	//name of the vertex, omit if null
	VarName *string

	//type of edge, omit if null
	Type *string

	//params for edge to map to, omit if null
	Params map[string]interface{}
}

func (v *V) ToCypher() (string, error) {
	//if nothing is specified, its just an empty vertex
	if v.VarName == nil && v.Type == nil && v.Params == nil{
		return "()", nil
	}

	str := "("

	//specify variable name if its there
	if v.VarName != nil{
		str += *v.VarName
	}

	//specify type if its there
	if v.Type != nil{
		str += ":" + *v.Type
	}

	//add params if its there
	if v.Params != nil{
		bytes, err := json.Marshal(v.Params)
		if err != nil{
			return "", err
		}

		str += " " + string(bytes)
	}

	str += ")"

	return str, nil
}

// E represents an edge
type E struct {
	//direction of the edge, if null default to any
	Direction *Direction

	//variable name for constraint queries, omit if null
	VarName *string

	//names in the case that the edge is named or the query could be on multiple edges
	Names []string

	//min jumps to the next node, if null omit
	MinJumps *int

	//max jumps to the next node, if null omit
	MaxJumps *int

	//params for edges across individual jumps
	Params map[string]interface{}
}

func (e *E) ToCypher() (string, error) {
	if e.VarName == nil && (e.Names == nil || len(e.Names) == 0) && e.MinJumps == nil || e.MaxJumps == nil || e.Params == nil{
		if e.Direction == nil{
			return "--", nil
		} else {
			if *e.Direction == Incoming{
				return "<--", nil
			} else {
				return "-->", nil
			}
		}
	}

	core := "["

	if e.VarName != nil{
		core += *e.VarName
	}

	if e.Names != nil && len(e.Names) != 0{
		if len(e.Names) == 1{
			core += ":" + e.Names[0]
		} else {
			q := ""
			for _, v := range e.Names{
				q += v + "|"
			}

			q = strings.TrimSuffix(q, "|")
			core += ":" + q
		}
	}

	if e.MinJumps != nil && e.MaxJumps != nil{
		q := fmt.Sprintf("*%v..%v", *e.MinJumps, *e.MaxJumps)
		core += q
	} else if e.MinJumps != nil{
		q := fmt.Sprintf("*%v", *e.MinJumps)
		core += q
	} else {
		q := fmt.Sprintf("*1..%v", *e.MaxJumps)
		core += q
	}

	if e.Params != nil{
		bytes, err := json.Marshal(e.Params)
		if err != nil{
			return "", err
		}

		core += " " + string(bytes)
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
	Any Direction = iota
	Outgoing
	Incoming

)