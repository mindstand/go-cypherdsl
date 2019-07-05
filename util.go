package go_cypherdsl


func StrPtr(s string) *string{
	return &s
}

func IntPtr(i int) *int{
	return &i
}

func DirectionPtr(d Direction) *Direction{
	return &d
}