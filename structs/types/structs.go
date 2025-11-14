package types

type Person struct {
	Name string
	age  int
}

type Employee struct {
	Person
	Position string
}
