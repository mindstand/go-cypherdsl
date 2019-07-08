package go_cypherdsl

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

func (q *QueryBuilder) Merge(m MergeQuery, err error) Cypher {
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(string(m))

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

func (q *QueryBuilder) Delete(d DeleteQuery, err error) Cypher {
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(string(d))
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

