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

type Match interface {
	Match(string, error) Cypher
}

type Create interface {
	Create(CreateQuery, error) Cypher
}

type Where interface {
	Where(WhereQuery, error) Cypher
}

type Merge interface {
	Merge(MergeQuery, error) Cypher
}

type Return interface {
	Return(ReturnQuery, error) Cypher
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
