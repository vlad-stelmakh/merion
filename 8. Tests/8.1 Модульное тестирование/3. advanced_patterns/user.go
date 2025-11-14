package main

import "regexp"

type User struct {
	Email string
	Name  string
	Age   int
}

func NewUser(name, email string, age int) *User {
	return &User{
		Name:  name,
		Email: email,
		Age:   age,
	}
}

func (u User) IsValidEmail() bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(u.Email)
}
