package controllers

import (
	"fmt"
	"net/http"

	"github.com/vitoraalmeida/lenslocked/models"
)

// desacopla o controller das views, injetando a instância
// do template que queremos executar como parte do tipo
// Users controller. Não interessa a implementação do template
// desde que obedeça a interface template (controlers/template.go)
type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService *models.UserService
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
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User craeted: %+v", user)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	u.Templates.SignIn.Execute(w, nil)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	cookie := http.Cookie{
		Name:     "email",
		Value:    user.Email,
		Path:     "/",  /*Define em que rotas na aplicação o cookie é acessível*/
		HttpOnly: true, /* define que JS não pode acessar cookies (XSS) */
	}
	http.SetCookie(w, &cookie) /*adiciona o header set-cookie no response*/
	fmt.Fprintf(w, "User authenticated: %+v", user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("email")
	if err != nil {
		fmt.Fprint(w, "The email cookie could not be read.")
		return
	}
	fmt.Fprintf(w, "Email cookie: %s\n", email.Value)
	fmt.Fprintf(w, "Headers: %+v", r.Header)
}
