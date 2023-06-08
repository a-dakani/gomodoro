package main

type ErrorMsg error

type Team struct {
	Name  string
	Slug  string
	Focus int64
	Pause int64
}

func (t *Team) IsEmpty() bool {
	return t.Name == "" && t.Slug == "" && t.Focus == 0 && t.Pause == 0
}
