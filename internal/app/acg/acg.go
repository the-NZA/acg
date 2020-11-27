package acg

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/the-NZA/acg/internal/app/store"
)

// Server ...
type Server struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

// New ...
func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start ...
func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info("Starting server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *Server) configureRouter() {
	// Website routes
	s.router.HandleFunc("/", s.handleHomePage())
	s.router.HandleFunc("/materials", s.handleMaterialsPage())
	s.router.HandleFunc("/posts", s.handlePostsPage())
	s.router.HandleFunc("/services", s.handleServicesPage())
	s.router.HandleFunc("/about", s.handleAboutPage())
	s.router.HandleFunc("/contacts", s.handleContactsPage())

	// Pages API routes
	s.router.HandleFunc("/api/pages", s.handleGetPages()).Methods("GET")
	s.router.HandleFunc("/api/pages", s.handleCreatePage()).Methods("POST")
	s.router.HandleFunc("/api/pages", s.handleUpdatePage()).Methods("PUT")

	// Services API routes
	s.router.HandleFunc("/api/services", s.handleGetServices()).Methods("GET")
	s.router.HandleFunc("/api/services", s.handleCreateService()).Methods("POST")
	s.router.HandleFunc("/api/services", s.handleUpdateService()).Methods("PUT")

	// Static files
	// TODO: Deliver this to NGINX later, with proxy and ssl
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
}

func (s *Server) configureStore() error {
	st := store.New(s.config.DatabaseURL)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st
	return nil
}
