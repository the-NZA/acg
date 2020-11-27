package acg

import (
	"encoding/json"
	"net/http"

	"github.com/the-NZA/acg/internal/app/store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
 * Pages Handlers for GET, POST and UPDATE
 */
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

		res, err := s.store.UpdateOne("pages", bson.M{"_id": np.ID}, bs)
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
