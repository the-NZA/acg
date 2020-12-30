package acg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/the-NZA/acg/internal/app/helpers"
	"github.com/the-NZA/acg/internal/app/store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bsbytes, err := bson.Marshal(np)
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
		mats, err := s.store.FindMaterials(bson.M{})
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pjs, err := json.Marshal(mats)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(pjs)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pjs, err := json.Marshal(mats)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(pjs)
	}
}

func (s *Server) handleGetOneMatcat() http.HandlerFunc {
	type req_body struct {
		Slug string `bson:"slug" json:"slug"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rb := &req_body{}

		err := json.NewDecoder(r.Body).Decode(rb)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bs, err := s.store.FindOne("matcategories", bson.M{"slug": rb.Slug})
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := bson.MarshalExtJSON(bs, true, true)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
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

		suf := time.Now().Format("02-01-2006_15-04")
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

		s.logger.Debug(req)

		u := &models.User{
			ID:                primitive.NewObjectID(),
			Username:          req.Username,
			EncryptedPassword: req.Password,
		}

		if err := u.Validate(); err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		res, err := s.store.InsertOne("users", u)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.logger.Debug(res.InsertedID)

		s.respond(w, r, http.StatusCreated, u)
	}
}
