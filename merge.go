package go_cypherdsl

import (
	"errors"
	"reflect"
	"strings"
)

type MergeConfig struct {
	// the path its merging on
	Path string

	// what it does if its creating the node
	OnCreate *MergeSetConfig

	// what it does if its matching the node
	OnMatch *MergeSetConfig
}

func (m *MergeConfig) ToString() (string, error) {
	var sb strings.Builder

	if m.Path == "" {
		return "", errors.New("path can not be empty")
	}

	sb.WriteString(m.Path)

	if m.OnCreate != nil {
		str, err := m.OnCreate.ToString()
		if err != nil {
			return "", err
		}

		sb.WriteString(" ON CREATE SET ")
		sb.WriteString(str)
	}

	if m.OnMatch != nil {
		str, err := m.OnMatch.ToString()
		if err != nil {
			return "", err
		}

		sb.WriteString(" ON MATCH SET ")
		sb.WriteString(str)
	}

	return sb.String(), nil
}

type MergeSetConfig struct {
	// variable name
	Name string

	// member variable of node
	Member string

	// new value
	Target interface{}

	// new value if its a function, do not include
	TargetFunction *FunctionConfig
}

func (m *MergeSetConfig) ToString() (string, error) {
	var sb strings.Builder

	if m.Name == "" {
		return "", errors.New("name can not be empty")
	}

	if m.Target == nil && m.TargetFunction == nil {
		return "", errors.New("target or target function has to be defined")
	}

	if m.Target != nil && m.TargetFunction != nil {
		return "", errors.New("target and target function can not both be defined")
	}

	if m.Target != nil && (reflect.TypeOf(m.Target) == reflect.TypeOf(ParamString(""))) {
		sb.WriteString(m.Name)
		sb.WriteRune(' ')
		sb.WriteString(EqualToOperator.String())
		sb.WriteRune(' ')
	} else {
		if m.Member == "" {
			return "", errors.New("member can not be empty")
		}

		sb.WriteString(m.Name)
		sb.WriteRune('.')
		sb.WriteString(m.Member)
		sb.WriteRune(' ')
		sb.WriteString(EqualToOperator.String())
		sb.WriteRune(' ')
	}

	var str string
	var err error

	if m.Target != nil {
		str, err = cypherizeInterface(m.Target)
	} else {
		str, err = m.TargetFunction.ToString()
	}

	if err != nil {
		return "", err
	}

	sb.WriteString(str)
	return sb.String(), nil
}
