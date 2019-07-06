package go_cypherdsl

type CreateQuery string

func (c *CreateQuery) ToString() string{
	return string(*c)
}

type WhereQuery string

func (c *WhereQuery) ToString() string{
	return string(*c)
}

type MergeQuery string

func (c *MergeQuery) ToString() string{
	return string(*c)
}

type ReturnQuery string

func (c *ReturnQuery) ToString() string{
	return string(*c)
}

type DeleteQuery string

func (c *DeleteQuery) ToString() string{
	return string(*c)
}

type SetQuery string

func (c *SetQuery) ToString() string{
	return string(*c)
}

type RemoveQuery string

func (c *RemoveQuery) ToString() string{
	return string(*c)
}