package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/vitoraalmeida/lenslocked/controllers"
	"github.com/vitoraalmeida/lenslocked/migrations"
	"github.com/vitoraalmeida/lenslocked/models"
	"github.com/vitoraalmeida/lenslocked/templates"
	"github.com/vitoraalmeida/lenslocked/views"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	// TODO: Read the PSQL values from an ENV variable
	cfg.PSQL = models.DefaultPostgresConfig()

	// TODO: SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	// chave necessária para o gorilla csrf criar um token aleatório
	// TODO: Read the CSRF values from an ENV variable
	cfg.CSRF.Key = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	cfg.CSRF.Secure = false

	// TODO: Read the server values from an ENV variable
	cfg.Server.Address = ":3000"

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Setup the database
	db, err := models.Open(cfg.PSQL)
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

	// setup services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	pwResetService := models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)

	// setup middlewares
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
	)

	// setup controllers
	usersC := controllers.Users{
		UserService:          &userService,
		SessionService:       &sessionService,
		PasswordResetService: &pwResetService,
		EmailService:         emailService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS, "signup.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS, "signin.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml",
	))

	// setup router
	r := chi.NewRouter()
	// utilzia a proteção csrf e o middleware de recuperação de usuário na requisição em todas as requisições. Primeiro aplica a recuperação do usuário no contexto e depois o csrf
	// o middleware que é registrado primeiro é o middleware que englobará
	// todos os restantes
	r.Use(csrfMw)
	r.Use(umw.SetUser)
	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))
	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))
	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	r.Get("/signup", usersC.New)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/users", usersC.Create)
	r.Post("/signout", usersC.ProcessSignOut)
	// cria um prefixo que possui rotas específicas em si e midlewares que tem
	// de ser usados para acessar determinados recursos
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}

}

// o middleware csrfMw exige que seja passado um token nas requisições que garantem que a requisição para o servidor
// foi originada de um formulário (ou outra forma de interação) que foi criada pelo próprio
// servidor. Se algum atacante tentar fazer uma cópia do sistema adicionando alguma interação
// com o usuário que por debaixo dos panos tenta performa uma ação indevida em nome do usuário
// no servidor se utilizando de sessões já criadas (CSRF - Cookies), o servidor tera como saber
// que é uma requisição inválida, pois esse token é colocado no frontend final do cliente
// autentico e apenas requisições partindo do form gerado pelo servidor terão esse token

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
