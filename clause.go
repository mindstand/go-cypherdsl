package go_cypherdsl

type Cypher interface {
	Match
	Create
	Where
	Merge
	Return
	Delete
	Set
	Remove
}

//complete
type Match interface {
	Match(string, error) Cypher
}

//complete
type Create interface {
	Create(CreateQuery, error) Cypher
}

//complete
type Where interface {
	Where(WhereQuery, error) Cypher
}

type Merge interface {
	Merge(MergeQuery, error) Cypher
}

//complete
type Return interface {
	Return(parts ...ReturnPart) Cypher
}

type Delete interface {
	Delete(DeleteQuery, error) Cypher
}

type Set interface {
	Set(SetQuery, error) Cypher
}

type Remove interface {
	Remove(RemoveQuery, error) Cypher
}

type OrderBy interface {
	OrderBy() Cypher
}

type Limit interface {
	Limit() Cypher
}