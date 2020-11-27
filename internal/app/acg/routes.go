package acg

import (
	"html/template"
	"net/http"

	"github.com/the-NZA/acg/internal/app/store/models"
	"go.mongodb.org/mongo-driver/bson"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("views/*.gohtml"))
}

// Homepage
func (s *Server) handleHomePage() http.HandlerFunc {
	type homepage struct {
		Page     models.Page
		Services []models.Service
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// s.logger.Printf("Host %s, Path %s\n", r.Host, r.URL.Path)
		m := models.Page{}
		srv := make([]models.Service, 0)

		bs, err := s.store.FindOne("pages", bson.M{"slug": r.URL.Path})
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		srv, err = s.store.FindAllServices()
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		bsb, err := bson.Marshal(bs)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		err = bson.Unmarshal(bsb, &m)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		homeContent := &homepage{m, srv}

		tpl.ExecuteTemplate(w, "index.gohtml", homeContent)
	}
}

// Posts page
func (s *Server) handlePostsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{}

		bs, err := s.store.FindOne("pages", bson.M{"slug": r.URL.Path})
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		bsb, err := bson.Marshal(bs)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		err = bson.Unmarshal(bsb, &m)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		tpl.ExecuteTemplate(w, "posts.gohtml", &m)
	}
}

// Services page
func (s *Server) handleServicesPage() http.HandlerFunc {
	type services struct {
		Page     models.Page
		Services []models.Service
	}

	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{}
		srv := make([]models.Service, 0)

		bs, err := s.store.FindOne("pages", bson.M{"slug": r.URL.Path})
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		srv, err = s.store.FindAllServices()
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		bsb, err := bson.Marshal(bs)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		err = bson.Unmarshal(bsb, &m)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		pageContect := &services{m, srv}

		tpl.ExecuteTemplate(w, "services.gohtml", pageContect)
	}
}

// Materials page
func (s *Server) handleMaterialsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{}

		bs, err := s.store.FindOne("pages", bson.M{"slug": r.URL.Path})
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		bsb, err := bson.Marshal(bs)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		err = bson.Unmarshal(bsb, &m)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		tpl.ExecuteTemplate(w, "materials.gohtml", &m)
	}
}

// About page
func (s *Server) handleAboutPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{}

		bs, err := s.store.FindOne("pages", bson.M{"slug": r.URL.Path})
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		bsb, err := bson.Marshal(bs)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		err = bson.Unmarshal(bsb, &m)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		tpl.ExecuteTemplate(w, "singlepage.gohtml", &m)
	}
}

// Contacts page
func (s *Server) handleContactsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{}

		bs, err := s.store.FindOne("pages", bson.M{"slug": r.URL.Path})
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		bsb, err := bson.Marshal(bs)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		err = bson.Unmarshal(bsb, &m)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		tpl.ExecuteTemplate(w, "singlepage.gohtml", &m)
	}
}
