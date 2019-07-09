package go_cypherdsl

import "errors"

//todo this might need to be renamed
type QueryBuilder struct {
	Start *queryPartNode
	Current *queryPartNode
	errors []error
}

func (q *QueryBuilder) addNext(s string) {
	node := &queryPartNode{
		Part: s,
	}

	if q.Start == nil{
		q.Start = node
		q.Current = node
	} else {
		q.Current.Next = node
		q.Current = node
	}
}

func (q *QueryBuilder) addError(err error){
	if q.errors == nil{
		q.errors = []error{}
	}

	q.errors = append(q.errors, err)
}

func (q *QueryBuilder) hasErrors() bool{
	return q.errors != nil && len(q.errors) > 0
}

type queryPartNode struct {
	Part string
	Next *queryPartNode
}

func (q *QueryBuilder) Match(s string, err error) Cypher {
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(s)
	return q
}

func (q *QueryBuilder) Create(c CreateQuery, err error) Cypher {
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(string(c))
	return q
}

func (q *QueryBuilder) Where(w WhereQuery, err error) Cypher {
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(string(w))
	return q
}

func (q *QueryBuilder) Merge(mergeConf *MergeConfig) Cypher {
	if mergeConf == nil{
		q.addError(errors.New("mergeConf can not be nil"))
		return q
	}
	cypher, err := mergeConf.ToString()
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(cypher)

	return q
}

func (q *QueryBuilder) Return(parts ...ReturnPart) Cypher {
	str, err := NewReturnClause(parts...)
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(string(str))
	return q
}

func (q *QueryBuilder) Delete(detach bool, params ...string) Cypher {
	cypher, err := deleteToString(detach, params...)
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(cypher)
	return q
}

func (q *QueryBuilder) Set(s SetQuery, err error) Cypher {
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(string(s))
	return q
}

func (q *QueryBuilder) Remove(r RemoveQuery, err error) Cypher {
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(string(r))
	return q
}

