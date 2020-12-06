package acg

import (
	"html/template"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/the-NZA/acg/internal/app/helpers"
	"github.com/the-NZA/acg/internal/app/store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tpl *template.Template
var pstsPerPage int64 = 15 // Number of posts per each page

// postspage represent struct for each page, which show posts
type postspage struct {
	Page       models.Page
	Posts      []models.Post
	Categories []models.Category
	Pagination []helpers.PaginationLink
	PagesCnt   int
	PageNum    string
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
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{}
		vars := mux.Vars(r)
		var pageNum int64

		bs, err := s.store.FindOne("pages", bson.M{"slug": "/posts"})
		if err != nil {
			s.logger.Info(err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		// Parse mux url vars to get pageNum
		if p, exist := vars["page"]; exist {
			pageNum, _ = strconv.ParseInt(p, 10, 64)
		} else {
			pageNum = 1
		}

		// Find options to deal with pagination
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

		// Calcs number of pages
		numOfPgs := math.Ceil(float64(pstsCnt) / float64(pstsPerPage))

		// Redirect to the first page if out of range
		if pageNum > int64(numOfPgs) {
			s.logger.Warn("Triggered unexisting page")
			http.Redirect(w, r, "/posts", http.StatusTemporaryRedirect)
			return
		}

		// Generate pagination slice
		pagiArr := helpers.GeneratePagination(int(pageNum), int(numOfPgs))

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

		// Fix for active menu link
		m.Slug = "/posts"

		postsContent := &postspage{
			Page:       m,
			Posts:      psts,
			Categories: cats,
			Pagination: pagiArr,
			PageNum:    strconv.Itoa(int(pageNum)),
			PagesCnt:   int(numOfPgs),
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
