package acg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/the-NZA/acg/internal/app/helpers"
	"github.com/the-NZA/acg/internal/app/store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	errWrongPasswordOrLogin = errors.New("Wrong username or password. Be more accurate.")
	errAlreadyExists        = errors.New("Username already taken. Try something else.")
	errWrongRequest         = errors.New("You provided incorrect data or did not provide them at all.")
	errUnauthorized         = errors.New("Something wrong with token, credentials or etc.")
)

/*
 * Helpers for Respond
 */
// TODO: Switch to s.respond and s.error methods for returning values
func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}

}

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// tknCookie, err := r.Cookie("TKN")
		// if err != nil {
		// 	s.error(w, r, http.StatusBadRequest, errWrongRequest)
		// 	return
		// }

		// tknString := tknCookie.Value

		tknHeader := r.Header.Get("Authorization")
		if tknHeader == "" {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		splited := strings.Split(tknHeader, " ")
		if len(splited) != 2 {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		tknString := splited[1]

		tkn, err := jwt.ParseWithClaims(tknString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected token signing: %v\n", t.Header["alg"])
			}

			return []byte(s.config.SecretKey), nil
		})

		if err != nil || !tkn.Valid {
			s.logger.Errorf("Err after parsing token: %v", err)
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		next(w, r)
	}
}

func (s *Server) testMiddlewareRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("this is test route")
		w.WriteHeader(http.StatusOK)
	}
}

/*
 * Pages Handlers for GET, POST and UPDATE
 */
// Handle GET on /pages
func (s *Server) handleGetPages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages, err := s.store.FindAllPages()
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pjs, err := json.Marshal(pages)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(pjs)
	}
}

// Handle POST on /pages
func (s *Server) handleCreatePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		np := &models.Page{
			ID: primitive.NewObjectID(),
		}

		err := json.NewDecoder(r.Body).Decode(np)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		np.Slug = fmt.Sprintf("/%s", helpers.GenerateSlug(np.Title))

		res, err := s.store.InsertOne("pages", np)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(res.InsertedID)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// Handle PUT on /pages
func (s *Server) handleUpdatePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		np := &models.Page{}

		err := json.NewDecoder(r.Body).Decode(np)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		bsbytes, err := bson.Marshal(np)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		var bs bson.M
		err = bson.Unmarshal(bsbytes, &bs)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, exist := bs["_id"]; exist {
			delete(bs, "_id")
		}

		res, err := s.store.UpdateOne("pages", bson.M{"_id": np.ID}, bson.M{"$set": bs})
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if res.UpsertedID == nil {
			bsbytes, err = json.Marshal(np.ID)
			if err != nil {
				s.logger.Error(err)
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		s.respond(w, r, http.StatusOK, bsbytes)
	}
}

/*
 * Services Handlers for GET, POST and UPDATE
 */
// Handle POST on /services
func (s *Server) handleCreateService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ns := &models.Service{
			ID:      primitive.NewObjectID(),
			Deleted: false,
		}

		err := json.NewDecoder(r.Body).Decode(ns)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = ns.Validate(); err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := s.store.InsertOne("services", ns)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(res.InsertedID)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// Handle GET on /services
func (s *Server) handleGetServices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services, err := s.store.FindAllServices()
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sjs, err := json.Marshal(services)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(sjs)
	}
}

// Handle PUT on /services
func (s *Server) handleUpdateService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ns := &models.Service{}

		err := json.NewDecoder(r.Body).Decode(ns)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bsbytes, err := bson.Marshal(ns)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var bs bson.M
		err = bson.Unmarshal(bsbytes, &bs)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, exist := bs["_id"]; exist {
			delete(bs, "_id")
		}

		res, err := s.store.UpdateOne("services", bson.M{"_id": ns.ID}, bson.M{"$set": bs})
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if res.UpsertedID == nil {
			bsbytes, err = json.Marshal(ns.ID)
			if err != nil {
				s.logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bsbytes)
	}
}

// Handle DELETE on /services
func (s *Server) handleDeleteService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ns := &models.Service{}

		err := json.NewDecoder(r.Body).Decode(ns)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Set deleted flag for service
		ns.Deleted = true

		bsbytes, err := bson.Marshal(ns)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		var bs bson.M
		err = bson.Unmarshal(bsbytes, &bs)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, exist := bs["_id"]; exist {
			delete(bs, "_id")
		}

		_, err = s.store.UpdateOne("services", bson.M{"_id": ns.ID}, bson.M{"$set": bs})
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, "Service successfully deleted")
	}
}

/*
* Posts Handlers for CRUD operations
 */
// Handel GET all posts on /posts
func (s *Server) handleGetPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := s.store.FindAllPosts(bson.M{"deleted": false})
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pjs, err := json.Marshal(posts)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(pjs)
	}
}

// Handle POST on /posts
func (s *Server) handleCreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		np := &models.Post{
			ID:   primitive.NewObjectID(),
			Time: time.Now(),
		}

		np.TimeString = np.Time.Format("02.01.2006")

		err := json.NewDecoder(r.Body).Decode(np)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		np.URL = fmt.Sprintf("%s/%s", np.CategoryURL, helpers.GenerateSlug(np.Title))

		if err = np.Validate(); err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := s.store.InsertOne("posts", np)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.store.UpdateOne("categories", bson.M{"title": np.Category}, bson.M{"$push": bson.M{"posts": np}})

		js, err := json.Marshal(res.InsertedID)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// Handle PUT on /posts
func (s *Server) handleUpdatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nc := &models.Post{}

		err := json.NewDecoder(r.Body).Decode(nc)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bsbytes, err := bson.Marshal(nc)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var bs bson.M
		err = bson.Unmarshal(bsbytes, &bs)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Delete ID to map (for corrent update)
		if _, exist := bs["_id"]; exist {
			delete(bs, "_id")
		}

		res, err := s.store.UpdateOne("posts", bson.M{"_id": nc.ID}, bson.M{"$set": bs})
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return ID to map (for corrent category update)
		bs["_id"] = nc.ID
		_, err = s.store.UpdateOne("categories", bson.M{"url": nc.CategoryURL, "posts._id": nc.ID}, bson.M{"$set": bson.M{"posts.$": bs}})
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if res.UpsertedID == nil {
			bsbytes, err = json.Marshal(nc.ID)
			if err != nil {
				s.logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bsbytes)
	}
}

// Handle DELETE on /posts
func (s *Server) handleDeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nc := &models.Post{}

		err := json.NewDecoder(r.Body).Decode(nc)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		nc.Deleted = true

		bsbytes, err := bson.Marshal(nc)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		var bs bson.M
		err = bson.Unmarshal(bsbytes, &bs)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Delete ID to map (for corrent update)
		if _, exist := bs["_id"]; exist {
			delete(bs, "_id")
		}

		_, err = s.store.UpdateOne("posts", bson.M{"_id": nc.ID}, bson.M{"$set": bs})
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, "Post successfully deleted")
	}
}

/*
* Categories Handlers for CRUD operations
 */
// Handel GET all categories on /categories
func (s *Server) handleGetCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cats, err := s.store.FindAllCategories(bson.M{})
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pjs, err := json.Marshal(cats)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(pjs)
	}
}

// Handle POST on /categories
func (s *Server) handleCreateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		np := &models.Category{
			ID: primitive.NewObjectID(),
		}

		err := json.NewDecoder(r.Body).Decode(np)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		np.URL = "/category/" + helpers.GenerateSlug(np.Title)

		if err = np.Validate(); err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := s.store.InsertOne("categories", np)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(res.InsertedID)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// Handle PUT on /categories
func (s *Server) handleUpdateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nc := &models.Category{}

		err := json.NewDecoder(r.Body).Decode(nc)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bsbytes, err := bson.Marshal(nc)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var bs bson.M
		err = bson.Unmarshal(bsbytes, &bs)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, exist := bs["_id"]; exist {
			delete(bs, "_id")
		}

		res, err := s.store.UpdateOne("categories", bson.M{"_id": nc.ID}, bson.M{"$set": bs})
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if res.UpsertedID == nil {
			bsbytes, err = json.Marshal(nc.ID)
			if err != nil {
				s.logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bsbytes)
	}
}

// Handle DELETE on /categories
func (s *Server) handleDeleteCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nc := &models.Category{}

		err := json.NewDecoder(r.Body).Decode(nc)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		nc.Deleted = true

		bsbytes, err := bson.Marshal(nc)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		var bs bson.M
		err = bson.Unmarshal(bsbytes, &bs)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, exist := bs["_id"]; exist {
			delete(bs, "_id")
		}

		_, err = s.store.UpdateOne("categories", bson.M{"_id": nc.ID}, bson.M{"$set": bs})
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		cat, err := s.store.FindOne("categories", bson.M{"_id": nc.ID})
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.logger.Infoln(cat.Map()["url"])

		posts, err := s.store.FindAllPosts(bson.M{"categoryurl": cat.Map()["url"]})
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		for i := range posts {
			posts[i].Deleted = true

			postbytes, err := bson.Marshal(posts[i])
			if err != nil {
				s.logger.Error(err)
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			var postbs bson.M
			err = bson.Unmarshal(postbytes, &postbs)
			if err != nil {
				s.logger.Error(err)
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			// Delete ID to map (for corrent update)
			if _, exist := postbs["_id"]; exist {
				delete(postbs, "_id")
			}

			_, err = s.store.UpdateOne("posts", bson.M{"_id": posts[i].ID}, bson.M{"$set": postbs})

		}

		// posts, err := s.store.FindAllPosts({"ca"})

		s.respond(w, r, http.StatusOK, "Category successfully deleted")
	}
}

/*
* Materials Handlers for CRUD operations
 */
// Handle POST materials on /materials
func (s *Server) handleCreateMaterial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nm := &models.Material{
			ID:   primitive.NewObjectID(),
			Time: time.Now(),
		}

		nm.TimeString = nm.Time.Format("02.01.2006")

		err := json.NewDecoder(r.Body).Decode(nm)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = nm.Validate(); err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := s.store.InsertOne("materials", nm)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.store.UpdateOne("matcategories", bson.M{"title": nm.Category}, bson.M{"$push": bson.M{"materials": nm}})

		js, err := json.Marshal(res.InsertedID)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// Handel GET materials on /materials
func (s *Server) handleGetMaterials() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mats, err := s.store.FindMaterials(bson.M{"deleted": false})
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, mats)
	}
}

/*
* MatCategories Handlers for CRUD operations
 */
// Handle GET matcategories on /api/matcategories
func (s *Server) handleGetMatcat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mats, err := s.store.FindMatcategories(bson.M{})
		if err != nil {
			s.logger.Error(err)
			s.respond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, mats)
	}
}

func (s *Server) handleGetOneMatcat() http.HandlerFunc {
	type matcat struct {
		models.MatCategory
		Materials []models.Material `bson:"materials,omitempty" json:"materials,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		sl := r.URL.Query().Get("slug")
		if sl == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("Query param not set. Try again."))
			return
		}

		bs, err := s.store.FindOne("matcategories", bson.M{"slug": sl})
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		materials, err := s.store.FindMaterials(bson.M{"category_slug": sl, "deleted": false})
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		bsbytes, err := bson.Marshal(bs)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		var matcategory models.MatCategory
		err = bson.Unmarshal(bsbytes, &matcategory)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, matcat{matcategory, materials})
	}
}

// Handle POST matcategories on /api/matcategories
func (s *Server) handleCreateMatcat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		np := &models.MatCategory{
			ID: primitive.NewObjectID(),
		}

		err := json.NewDecoder(r.Body).Decode(np)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		np.Slug = helpers.GenerateSlug(np.Title)

		if err = np.Validate(); err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := s.store.InsertOne("matcategories", np)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(res.InsertedID)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// Handle PUT matcategories on /api/matcategories
func (s *Server) handleUpdateMatcat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nc := &models.MatCategory{}

		err := json.NewDecoder(r.Body).Decode(nc)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bsbytes, err := bson.Marshal(nc)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var bs bson.M
		err = bson.Unmarshal(bsbytes, &bs)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, exist := bs["_id"]; exist {
			delete(bs, "_id")
		}

		res, err := s.store.UpdateOne("matcategories", bson.M{"_id": nc.ID}, bson.M{"$set": bs})
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if res.UpsertedID == nil {
			bsbytes, err = json.Marshal(nc.ID)
			if err != nil {
				s.logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bsbytes)
	}
}

/* Upload routes */
// Handle POST on /upload
func (s *Server) handleUploadFile() http.HandlerFunc {
	type res struct {
		OK   bool   `json:"ok"`
		Name string `json:"name"`
		Size int64  `json:"size"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Set max upload size at ~32MB
		r.ParseMultipartForm(32 << 20)

		tf, h, err := r.FormFile("acg_upload")
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tf.Close()

		// s.logger.Infof("uploaded file: %+v\n", h.Filename)
		// s.logger.Infof("file size: %+v bytes\n", h.Size)
		// s.logger.Infof("MIME type: %+v\n", h.Header)

		suf := time.Now().Format("02-01-2006_15-04-05")
		upload_path := "uploads/" + suf + "_" + h.Filename

		f, err := os.OpenFile(upload_path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		wrtn, err := io.Copy(f, tf)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.logger.Infof("File %s uploaded in %s. Size %v\n", h.Filename, f.Name(), wrtn)

		resp := &res{
			OK:   true,
			Name: upload_path,
			Size: wrtn,
		}

		js, err := json.Marshal(resp)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

/* Auth Routes */
func (s *Server) handleRegistration() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if req.Username == "" || req.Password == "" {
			s.error(w, r, http.StatusBadRequest, errWrongPasswordOrLogin)
			return
		}

		s.logger.Info(req)

		// If Username already taken than return custom error
		if _, err := s.store.FindOne("users", bson.M{"username": req.Username}); err == nil {
			s.error(w, r, http.StatusNotAcceptable, errAlreadyExists)
			return
		}

		u := &models.User{
			ID:       primitive.NewObjectID(),
			Username: req.Username,
		}

		if err := u.HashPassword(req.Password); err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := u.Validate(); err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		_, err := s.store.InsertOne("users", u)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, u)
	}
}

type Claims struct {
	Username string
	jwt.StandardClaims
}

func (s *Server) handleLogin() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if req.Password == "" || req.Username == "" {
			s.logger.Debug("Empty username or password")
			s.error(w, r, http.StatusBadRequest, errWrongPasswordOrLogin)
			return
		}

		bs, err := s.store.FindOne("users", bson.M{"username": req.Username})
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, errWrongPasswordOrLogin)
			return
		}

		u := &models.User{}

		bs_bytes, err := bson.Marshal(bs)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := bson.Unmarshal(bs_bytes, u); err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if !u.CheckPassword(req.Password) {
			s.logger.Debug("Wrong email or password")
			s.error(w, r, http.StatusUnauthorized, errWrongPasswordOrLogin)
			return
		}

		expTime := time.Now().Add(12 * time.Hour)
		// expTime := time.Now().Add(2 * time.Minute)
		claims := &Claims{
			Username: u.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(s.config.SecretKey))
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "TKN",
			Value:    tokenStr,
			Expires:  expTime,
			HttpOnly: true,
		})

		// w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
		s.respond(w, r, http.StatusOK, struct {
			Token string       `json:"token"`
			User  *models.User `json:"user"`
		}{tokenStr, u})
	}
}
