package acg

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/the-NZA/acg/internal/app/store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tpl *template.Template
var pstsPerPage int64 = 1 // Number of posts per each page

// postspage represent struct for each page, which show posts
type postspage struct {
	Page       models.Page
	Posts      []models.Post
	Categories []models.Category
	Pagination interface{}
}

func init() {
	tpl = template.Must(template.ParseGlob("views/*.gohtml"))
}

// Homepage
func (s *Server) handleHomePage() http.HandlerFunc {
	type homepage struct {
		Page     models.Page
		Services []models.Service
		Posts    []models.Post
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

		findOptions := options.Find()
		findOptions.SetLimit(3)
		findOptions.SetSort(bson.M{"time": -1})
		psts, err := s.store.FindAllPosts(findOptions)

		bsb, err := bson.Marshal(bs)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		err = bson.Unmarshal(bsb, &m)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		homeContent := &homepage{m, srv, psts}

		tpl.ExecuteTemplate(w, "index.gohtml", homeContent)
	}
}

// Posts page
func (s *Server) handlePostsPage() http.HandlerFunc {
	type pagination struct {
		first string
		last  string
		next  string
		prev  string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{}
		vars := mux.Vars(r)
		var pageNum int64

		s.logger.Debug(vars)

		bs, err := s.store.FindOne("pages", bson.M{"slug": "/posts"})
		if err != nil {
			s.logger.Info(err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		if p, exist := vars["page"]; exist {
			s.logger.Info(vars["page"])
			pageNum, _ = strconv.ParseInt(p, 10, 64)
		} else {
			pageNum = 1
		}

		findOptions := options.Find()
		findOptions.SetLimit(pstsPerPage)
		findOptions.SetSort(bson.M{"time": -1})
		findOptions.SetSkip((pageNum - 1) * pstsPerPage)

		psts, err := s.store.FindAllPosts(findOptions)
		if err != nil {
			s.logger.Error(err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
			return
		}

		pstsCnt, err := s.store.CountAllPosts()
		if err != nil {
			s.logger.Error(err)
		}

		if pstsCnt < pstsPerPage*pageNum {
			s.logger.Warn("triggered unexisting page")
			http.Redirect(w, r, "/posts", http.StatusTemporaryRedirect)
			return
		}

		cats, err := s.store.FindAllCategories(bson.M{})
		if err != nil {
			s.logger.Error(err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		bsb, err := bson.Marshal(bs)
		if err != nil {
			s.logger.Error(err)
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		err = bson.Unmarshal(bsb, &m)
		if err != nil {
			s.logger.Error(err)
			http.Redirect(w, r, "/404", http.StatusInternalServerError)
		}

		postsContent := &postspage{
			Page:       m,
			Posts:      psts,
			Categories: cats,
		}

		tpl.ExecuteTemplate(w, "posts.gohtml", postsContent)
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
