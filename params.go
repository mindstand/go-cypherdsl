package go_cypherdsl

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Params struct {
	params map[string]string
}

func ParamsFromMap(m map[string]interface{}) (*Params, error) {
	if m == nil || len(m) == 0 {
		return nil, errors.New("map can not be empty or nil")
	}

	p := &Params{}

	for k, v := range m {
		err := p.Set(k, v)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *Params) IsEmpty() bool {
	return p.params == nil || len(p.params) == 0
}

func (p *Params) Set(key string, value interface{}) error {
	if p.params == nil {
		p.params = map[string]string{}
	}

	str, err := cypherizeInterface(value)
	if err != nil {
		return err
	}

	p.params[key] = fmt.Sprintf("%s:%s", key, str)

	return nil
}

func (p *Params) ToCypherMap() string {

	if p.params == nil || len(p.params) == 0 {
		return "{}"
	}

	// Create and sort a slice of keys to ensure cypher queries are created with consistent ordering
	keys := make([]string, 0, len(p.params))
	for k, _ := range p.params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s,", p.params[k]))
	}

	return "{" + strings.TrimSuffix(sb.String(), ",") + "}"
}
