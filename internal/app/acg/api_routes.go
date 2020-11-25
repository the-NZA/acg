package acg

import (
	"encoding/json"
	"net/http"

	"github.com/the-NZA/acg/internal/app/store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Pages Handlers for GET, POST and UPDATE
func (s *Server) handleGetPages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := s.store.FindAll("pages")
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pages := []models.Page{}
		curItem := models.Page{}

		for _, v := range data {
			bsonBytes, _ := bson.Marshal(v)
			bson.Unmarshal(bsonBytes, &curItem)
			pages = append(pages, curItem)
		}

		bytes, err := json.Marshal(pages)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
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
