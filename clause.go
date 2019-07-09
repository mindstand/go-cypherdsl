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

//complete
type Merge interface {
	Merge(mergeConf *MergeConfig) Cypher
}

//complete
type Return interface {
	Return(parts ...ReturnPart) Cypher
}

//complete
type Delete interface {
	Delete(detach bool, params ...string) Cypher
}

//complete
type Set interface {
	Set(sets ...SetConfig) Cypher
}

type Remove interface {
	Remove(removes ...RemoveConfig) Cypher
}

type OrderBy interface {
	OrderBy() Cypher
}

type Limit interface {
	Limit() Cypher
}