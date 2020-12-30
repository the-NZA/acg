package acg

import (
	"fmt"
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

	s.router.HandleFunc("/posts", s.handlePostsPage())
	s.router.HandleFunc("/posts/", s.handlePostsPage())
	s.router.HandleFunc("/posts/{page}", s.handlePostsPage())

	s.router.HandleFunc("/materials", s.handleMaterialsPage())
	s.router.HandleFunc("/services", s.handleServicesPage())
	s.router.HandleFunc("/about", s.handleAboutPage())
	s.router.HandleFunc("/contacts", s.handleContactsPage())

	s.router.HandleFunc("/category", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/posts", http.StatusTemporaryRedirect)
	})
	s.router.HandleFunc("/category/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/posts", http.StatusTemporaryRedirect)
	})
	s.router.HandleFunc("/category/{cat}", s.handleCategoryPage())
	s.router.HandleFunc("/category/{cat}/{post}", s.handleSinglePost())

	/* API ROUTES */
	// Upload API route
	s.router.HandleFunc("/api/upload", s.handleUploadFile()).Methods("POST")

	// Auth API routes
	s.router.HandleFunc("/api/auth/reg", s.handleRegistration()).Methods("POST")
	s.router.HandleFunc("/api/auth/login", s.handleLogin()).Methods("POST")

	// Pages API routes
	s.router.HandleFunc("/api/pages", s.handleGetPages()).Methods("GET")
	s.router.HandleFunc("/api/pages", s.handleCreatePage()).Methods("POST")
	s.router.HandleFunc("/api/pages", s.handleUpdatePage()).Methods("PUT")

	// Services API routes
	s.router.HandleFunc("/api/services", s.handleGetServices()).Methods("GET")
	s.router.HandleFunc("/api/services", s.handleCreateService()).Methods("POST")
	s.router.HandleFunc("/api/services", s.handleUpdateService()).Methods("PUT")

	// Posts API routes
	s.router.HandleFunc("/api/posts", s.handleGetPosts()).Methods("GET")
	s.router.HandleFunc("/api/posts", s.handleCreatePost()).Methods("POST")
	s.router.HandleFunc("/api/posts", s.handleUpdatePost()).Methods("PUT")

	// Categories API routes
	s.router.HandleFunc("/api/categories", s.handleGetCategories()).Methods("GET")
	s.router.HandleFunc("/api/categories", s.handleCreateCategory()).Methods("POST")
	s.router.HandleFunc("/api/categories", s.handleUpdateCategory()).Methods("PUT")

	// Materials API routes
	s.router.HandleFunc("/api/materials", s.handleGetMaterials()).Methods("GET")
	s.router.HandleFunc("/api/materials", s.handleCreateMaterial()).Methods("POST")
	// s.router.HandleFunc("/api/materials", s.handleUpdateMaterial()).Methods("PUT")

	// MatCategories API routes
	s.router.HandleFunc("/api/matcategory", s.handleGetOneMatcat()).Methods("GET")
	s.router.HandleFunc("/api/matcategories", s.handleGetMatcat()).Methods("GET")
	s.router.HandleFunc("/api/matcategories", s.handleCreateMatcat()).Methods("POST")
	s.router.HandleFunc("/api/matcategories", s.handleUpdateMatcat()).Methods("PUT")

	s.router.HandleFunc("/404", notFound)
	// 404 Handler
	// TODO: Add good page for 404 with view and some additional info
	s.router.NotFoundHandler = http.HandlerFunc(notFound)
}

func (s *Server) configureStore() error {
	st := store.New(s.config.DatabaseURL)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st
	return nil
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>This is 404 page. Sorry...</h1>\n")
}
