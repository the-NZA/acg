package acg

import (
	"html/template"
	"net/http"

	"github.com/the-NZA/acg/internal/app/store/models"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("views/*.gohtml"))
}

// Homepage
func (s *Server) handleHomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{
			Title:    "Добро пожаловать, пожаловать добро!",
			MetaDesc: "This is awesome mock description",
		}

		tpl.ExecuteTemplate(w, "index.gohtml", &m)
	}

}

// Posts page
func (s *Server) handlePostsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{
			Title:    "Новости и публикации",
			MetaDesc: "Здесь вы найдете новые публикации и новостные сводки от нашей команды.",
		}

		tpl.ExecuteTemplate(w, "posts.gohtml", &m)
	}
}

// Services page
func (s *Server) handleServicesPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{
			Title:    "Наши услуги",
			MetaDesc: "Весь спектр услуг представлен в удобном виде",
		}

		tpl.ExecuteTemplate(w, "services.gohtml", &m)
	}

}

// Materials page
func (s *Server) handleMaterialsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{
			Title:    "Материалы и обсуждение",
			MetaDesc: "Здесь вы найдете выборку и примеры готовых файлов по нужным тематикам",
		}

		tpl.ExecuteTemplate(w, "materials.gohtml", &m)
	}
}

// About page
func (s *Server) handleAboutPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{
			Title:    "О компании",
			MetaDesc: "Вся информация о нашей компании может быть найдена здесь",
		}

		tpl.ExecuteTemplate(w, "singlepage.gohtml", &m)
	}
}

// Contacts page
func (s *Server) handleContactsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := models.Page{
			Title:    "Контакты",
			MetaDesc: "Подробная информация о способах связи с нами",
		}

		tpl.ExecuteTemplate(w, "singlepage.gohtml", &m)
	}
}
