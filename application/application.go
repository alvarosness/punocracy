package application

import (
	"net/http"

	"github.com/carbocation/interpose"
	_ "github.com/go-sql-driver/mysql"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	"github.com/alvarosness/punocracy/handlers"
	"github.com/alvarosness/punocracy/middlewares"
)

// New is the constructor for Application struct.
func New(config *viper.Viper) (*Application, error) {
	dsn := config.Get("dsn").(string)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	cookieStoreSecret := config.Get("cookie_secret").(string)

	app := &Application{}
	app.config = config
	app.dsn = dsn
	app.db = db
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, nil
}

// Application is the application object that runs HTTP server.
type Application struct {
	config       *viper.Viper
	dsn          string
	db           *sqlx.DB
	sessionStore sessions.Store
}

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetDB(app.db))
	middle.Use(middlewares.SetSessionStore(app.sessionStore))

	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *gorilla_mux.Router {
	MustLogin := middlewares.MustLogin

	router := gorilla_mux.NewRouter()

	router.Handle("/", http.HandlerFunc(handlers.GetHome)).Methods("GET")

	router.HandleFunc("/submit", handlers.GetSubmit).Methods("GET")
	router.HandleFunc("/submit", handlers.PostSubmit).Methods("POST")

	router.HandleFunc("/history", handlers.GetHistory).Methods("GET")
	router.HandleFunc("/history", handlers.PostHistory).Methods("POST")

	router.HandleFunc("/words", handlers.GetWords).Methods("GET")
	router.HandleFunc("/words", handlers.PostWords).Methods("POST")

	router.HandleFunc("/queuerater", handlers.GetCurator).Methods("GET")
	router.HandleFunc("/queuerater", handlers.PostCurator).Methods("POST")

	router.HandleFunc("/about", handlers.GetAbout).Methods("GET")
	router.HandleFunc("/about", handlers.PostAbout).Methods("POST")

	router.HandleFunc("/signup", handlers.GetSignup).Methods("GET")
	router.HandleFunc("/signup", handlers.PostSignup).Methods("POST")

	router.HandleFunc("/login", handlers.GetLogin).Methods("GET")
	router.HandleFunc("/login", handlers.PostLogin).Methods("POST")

	router.HandleFunc("/logout", handlers.GetLogout).Methods("GET")

	router.Handle("/users/{id:[0-9]+}", MustLogin(http.HandlerFunc(handlers.PostPutDeleteUsersID))).Methods("POST", "PUT", "DELETE")

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return router
}
