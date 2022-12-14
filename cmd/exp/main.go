package main

import (
	"html/template"
	"os"
	//"text/template"
)

type User struct {
	Name string
	Age  int
	Bio  string // se não quiser que encode usar o tipo template.HTML
	Meta UserMeta
}

type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("hello.tmpl")
	if err != nil {
		panic(err)
	}

	// usando text/template, os caracteres <> são interpretados literalmente
	// então o código script vai ser executado e o comando js vai rodar
	// XSS
	// o html/template fazer o escape dos caracteres especiais, evitando isso

	user := User{
		Name: "John Smith",
		Age:  19,
		Bio:  `<script>alert("haha, you have been hacked!");</script>`,
		Meta: UserMeta{
			Visits: 4,
		},
	}

	user2 := struct {
		Name string
		Age  int
		Bio  string
		Meta struct {
			Visits int
		}
	}{
		Name: "Susan Smith",
		Age:  29,
		Bio:  `<script>alert("haha, you have been h4x0r3d!");</script>`,
		Meta: struct {
			Visits int
		}{
			Visits: 19,
		},
	}
	// recebe onde será escrito o resultado e o dado que será usando no template
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, user2)
	if err != nil {
		panic(err)
	}
}
