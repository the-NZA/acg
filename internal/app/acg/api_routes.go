package acg

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/the-NZA/acg/internal/app/store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// js, err := json.Marshal(np)
		// if err != nil {
		// 	s.logger.Error(err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		bsbytes, err := bson.Marshal(np)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var bs bson.M
		err = bson.Unmarshal(bsbytes, &bs)
		// err = bson.UnmarshalExtJSON([]byte(js), true, &bs)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, exist := bs["_id"]; exist {
			delete(bs, "_id")
		}

		// js, err = json.Marshal(bs)
		// if err != nil {
		// 	s.logger.Error(err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// err = bson.UnmarshalExtJSON([]byte(js), true, &bs)
		// if err != nil {
		// 	s.logger.Error(err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		res, err := s.store.UpdateOne("pages", bson.M{"_id": np.ID}, bson.M{"$set": bs})
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if res.UpsertedID == nil {
			bsbytes, err = json.Marshal(np.ID)
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

/*
 * Services Handlers for GET, POST and UPDATE
 */
// Handle POST on /services
func (s *Server) handleCreateService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ns := &models.Service{
			ID: primitive.NewObjectID(),
		}

		err := json.NewDecoder(r.Body).Decode(ns)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

/*
* Posts Handlers for CRUD operations
 */
// Handel GET all posts on /posts
func (s *Server) handleGetPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := s.store.FindAllPosts()
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

		if _, exist := bs["_id"]; exist {
			delete(bs, "_id")
		}

		res, err := s.store.UpdateOne("posts", bson.M{"_id": nc.ID}, bson.M{"$set": bs})
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

/*
* Categories Handlers for CRUD operations
 */
// Handel GET all categories on /categories
func (s *Server) handleGetCatigories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cats, err := s.store.FindAllCategories()
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
