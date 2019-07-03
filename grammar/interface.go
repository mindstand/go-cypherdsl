package grammar

type Cypherize interface {
	ToCypher() (string, error)
}
