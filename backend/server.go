package server

import (
	"fmt"
	// "io/ioutil"
	
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/justinas/alice"

	"github.com/is0405/controller"
	"github.com/is0405/db"
	"github.com/is0405/middleware"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
)

type Server struct {
	db           *sqlx.DB
	router       *mux.Router
	jwtSecretKey []byte
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init(datasource string) error {
	cs := db.NewDB(datasource)
	dbcon, err := cs.Open()
	if err != nil {
		return fmt.Errorf("failed db init. %s", err)
	}
	s.db = dbcon

	s.router = s.Route()
	return nil
}

func (s *Server) Run(port int) {
	log.Printf("Listening on port %d", port)
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		handlers.CombinedLoggingHandler(os.Stdout, s.router),
	)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Route() *mux.Router {
	authMiddleware := middleware.NewAuth(s.jwtSecretKey, s.db)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
	})
	
	commonChain := alice.New(
		middleware.RecoverMiddleware,
		corsMiddleware.Handler,
	)
	authChain := commonChain.Append(
		authMiddleware.Handler,
	)
	
	r := mux.NewRouter()
	Controlloer := controller.NewToDO(s.db)
	r.Methods(http.MethodPost).Path("/todo").Handler(authChain.Then(AppHandler{Controlloer.Create}))
	r.Methods(http.MethodGet).Path("/todo").Handler(authChain.Then(AppHandler{Controlloer.Get}))
	r.Methods(http.MethodPatch).Path("/todo/{id}").Handler(authChain.Then(AppHandler{Controlloer.Update}))
	r.Methods(http.MethodDelete).Path("/todo/{id}").Handler(authChain.Then(AppHandler{Controlloer.Delete}))

	loginController := controller.NewLogin(s.db, s.jwtSecretKey)
	r.Methods(http.MethodPost).Path("/todo/login").Handler(commonChain.Then(AppHandler{loginController.Login}))

	accountController := controller.NewAccount(s.db)
	r.Methods(http.MethodPost).Path("/todo/account").Handler(commonChain.Then(AppHandler{accountController.Create}))
	
	return r
}
