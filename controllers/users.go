package controllers

import (
	"fmt"
	"net/http"
)

// desacopla o controller das views, injetando a instância
// do template que queremos executar como parte do tipo
// Users controller. Não interessa a implementação do template
// desde que obedeça a interface template (controlers/template.go)
type Users struct {
	Templates struct {
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	/*
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, "Email: ", r.PostForm.Get("email"))
		fmt.Fprint(w, "Pass: ", r.PostForm.Get("password"))
	*/

	// não checará erro pois se esses valores não estivem presentes, não há nada
	// para fazer além de retornar erro
	fmt.Fprint(w, "Email: ", r.FormValue("email"))
	fmt.Fprint(w, "Pass: ", r.FormValue("password"))
}
