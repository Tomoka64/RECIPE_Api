package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Tomoka64/RECIPE_Api/internal/postgres"

	"github.com/julienschmidt/httprouter"
)

type Handler = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params)

type Recipe struct {
	ID          int       `json:"recipe_id"`
	Topic       string    `json:"topic"`
	Description string    `json:"description"`
	UserId      int       `json:"user_id"`
	PrepTime    int       `json:"prep_time`
	Difficulty  int       `json:"difficulty"`
	Vegetarian  bool      `json:"vegetarian"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *Server) handler() error {
	router := httprouter.New()
	http.Handle("/", router)
	router.GET("/", s.Index())
	router.GET("/recipes", s.List())           // not restricted
	router.POST("/recipes", s.Create())        // restricted
	router.GET("/recipes/{id}", s.Get())       // not restricted
	router.PATCH("/recipes/{id}", s.Update())  // restricted
	router.DELETE("/recipes/{id}", s.Delete()) // restricted
	router.POST("/recipes/{id}/rating", Rate)  // not restricted

	router.GET("/form/login", Login)
	router.GET("/form/signup", Signup)

	// router.POST("/api/createuser", createUser)
	// router.POST("/api/login", LoginProcess)
	// router.GET("/api/logout", logout)
	// router.GET("/post/done", Done)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Index() Handler {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.Redirect(w, r, "/recipes", http.StatusSeeOther)
	}
}

func (s *Server) List() Handler {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		recipes, err := postgres.GetAllRecipes(s.DB)
		if err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		s.GenerateHTML(w, r, "layout.html", &Data{Recipes: recipes})
	}
}

func (s *Server) LoginProcess() Handler {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		username := r.FormValue("username")
		email := r.FormValue("email")
		if s.UserAndPasswordMatch(username, email) {
			sess, err := s.Store.Get(r, "session")
			if err != nil {
				fmt.Fprintln(w, http.StatusInternalServerError)
				return
			}
			sess.Values["user_name"] = username
			sess.Save(r, w)
		}
		return
	}
}

func (s *Server) Get() Handler {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		recipeId := ps.ByName("id")
		recipe, err := postgres.GetRecipeByID(s.DB, recipeId)
		if err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		s.GenerateHTML(w, r, "layout.html", &Data{Recipe: recipe})
	}
}

func (s *Server) Create() Handler {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		sess, err := s.Store.Get(r, "session")
		if err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		if !Loggedin(sess) {
			fmt.Fprintln(w, http.StatusUnauthorized)
			return
		}
		var recipe postgres.Recipe
		recipeId, err := strconv.Atoi(ps.ByName("id"))
		if err != nil {
			fmt.Fprintln(w, http.StatusUnauthorized)
			return
		}
		recipe.ID = recipeId
		if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		if err := recipe.CreateRecipe(s.DB); err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) Update() Handler {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		sess, err := s.Store.Get(r, "session")
		if err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		if !Loggedin(sess) {
			fmt.Fprintln(w, http.StatusUnauthorized)
			return
		}
		var recipe postgres.Recipe
		recipeId, err := strconv.Atoi(ps.ByName("id"))
		if err != nil {
			fmt.Fprintln(w, http.StatusUnauthorized)
			return
		}
		recipe.ID = recipeId
		if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		if err := recipe.UpdateRecipe(s.DB); err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) Delete() Handler {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		sess, err := s.Store.Get(r, "session")
		if err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		if !Loggedin(sess) {
			fmt.Fprintln(w, http.StatusUnauthorized)
			return
		}
		var recipe postgres.Recipe
		recipeId, err := strconv.Atoi(ps.ByName("id"))
		if err != nil {
			fmt.Fprintln(w, http.StatusUnauthorized)
			return
		}
		recipe.ID = recipeId
		if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		if err := recipe.DeleteRecipe(s.DB); err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
	}
}

func Rate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
