package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/vitoraalmeida/lenslocked/controllers"
	"github.com/vitoraalmeida/lenslocked/migrations"
	"github.com/vitoraalmeida/lenslocked/models"
	"github.com/vitoraalmeida/lenslocked/templates"
	"github.com/vitoraalmeida/lenslocked/views"
)

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// aplica migrations automaticamente utlizando arquivos
	// de migration embutidos no binário, assim não precisamos
	// copiar arquivos de migration para o local de produção
	// manualmente
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	r.Get("/signup", usersC.New)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/users", usersC.Create)
	r.Get("/users/me", usersC.CurrentUser)
	r.Post("/signout", usersC.ProcessSignOut)

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting ther server on :3000...")
	// chave necessária para o gorilla csrf criar um token aleatório
	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying
		csrf.Secure(false), // a proteção necessita que seja usada numa conexão com https(prod)
	)
	http.ListenAndServe(":3000", csrfMw(r)) // utilzia a proteção csrf em todas as requisições
	// exige que seja passado um token nas requisições que garantem que a requisição para o servidor
	// foi originada de um formulário (ou outra forma de interação) que foi criada pelo próprio
	// servidor. Se algum atacante tentar fazer uma cópia do sistema adicionando alguma interação
	// com o usuário que por debaixo dos panos tenta performa uma ação indevida em nome do usuário
	// no servidor se utilizando de sessões já criadas (CSRF - Cookies), o servidor tera como saber
	// que é uma requisição inválida, pois esse token é colocado no frontend final do cliente
	// autentico e apenas requisições partindo do form gerado pelo servidor terão esse token
}

// middlewares recebem como argumento a função que originalmente deve ser executada
// podem executar código antes, executam a função original e depois podem executar
// mais código.

// func TimerMiddleware(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		h(w, r)
// 		fmt.Println("Request time:", time.Since(start))
// 	}
// }
