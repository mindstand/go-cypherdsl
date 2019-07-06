package go_cypherdsl

import (
	"errors"
	"fmt"
)

type CreateBuilder struct {

}

type IndexBuilder struct {

}

type ConstraintConfig struct {
	Name string
	Type string
	Field string
	Unique bool
	Exists bool
}

func NewConstraint(constraint *ConstraintConfig) (CreateQuery, error){
	if constraint == nil{
		return "", errors.New("constraint can not be nil")
	}

	if constraint.Name == "" || constraint.Type == "" || constraint.Field == ""{
		return "", errors.New("name, type and field can not be empty")
	}

	if constraint.Unique == constraint.Exists{
		return "", errors.New("can only be unique or exists per call")

	}

	if !constraint.Unique || constraint.Exists{
		return "", errors.New("constraint has to be defined")
	}

	root := "CONSTRAINT ON ("

	root += fmt.Sprintf("%s:%s) ASSERT ", constraint.Name, constraint.Type)

	if constraint.Unique {
		root += fmt.Sprintf("%s.%s IS UNIQUE", constraint.Name, constraint.Field)
	} else {
		root += fmt.Sprintf("exists(%s.%s)", constraint.Name, constraint.Field)
	}

	return CreateQuery(root), nil
}