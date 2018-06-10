package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Tomoka64/RECIPE_Api/internal/postgres"
	"github.com/Tomoka64/RECIPE_Api/internal/router"
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
	r := router.New(s.Logger)
	r.GET("/", s.Index())
	r.GET("/recipes", s.List())              // not restricted
	r.POST("/recipes", s.Create())           // restricted
	r.GET("/recipes/{id}", s.Get())          // not restricted
	r.PATCH("/recipes/{id}", s.Update())     // restricted
	r.DELETE("/recipes/{id}", s.Delete())    // restricted
	r.POST("/recipes/{id}/rating", s.Rate()) // not restricted

	r.GET("/form/login", s.Login())   //just to serve template
	r.GET("/form/signup", s.Signup()) //just to serve template

	r.POST("/api/createuser", s.CreateUser())
	r.POST("/api/login", s.LoginProcess())
	r.GET("/api/logout", s.Logout())

	// router.GET("/post/done", Done)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	return server.ListenAndServe()
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
			s.Logger.Info(err.Error())
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
			sess.Options.Secure = true
			sess.Options.HttpOnly = true
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
		w.Header().Set("Content-Type", "application/json")
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

func (s *Server) Rate() Handler {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		recipiId := ps.ByName("id")
		if err := postgres.UpdateRate(s.DB, recipiId); err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/recipes/%s", recipiId), http.StatusSeeOther)
	}
}

func (s *Server) Login() Handler {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		sess, err := s.Store.Get(r, "session")
		if err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		if Loggedin(sess) {
			http.Redirect(w, r, "/recipes", http.StatusSeeOther)
		}
		s.Tpl.ExecuteTemplate(w, "login", nil)
	}
}

func (s *Server) Logout() Handler {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		sess, err := s.Store.Get(r, "session")
		if err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		sess.Options.MaxAge = -1
		sess.Save(r, w)
		http.Redirect(w, r, "/recipes", http.StatusSeeOther)
	}
}

func (s *Server) Signup() Handler {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		sess, err := s.Store.Get(r, "session")
		if err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
		if !Loggedin(sess) {
			s.Tpl.ExecuteTemplate(w, "signup", nil)
		}
		http.Redirect(w, r, "/recipes", http.StatusSeeOther)
	}
}

func (s *Server) CreateUser() Handler {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var user *postgres.User
		name, password, email := r.FormValue("name"), r.FormValue("password"), r.FormValue("email")
		if !IsPasswordValid(password) {
			return
		}
		if s.UserExists(name, email) {

			return
		}
		user.Name, user.Password, user.Email = name, password, email
		if err := user.CreateUser(s.DB); err != nil {
			fmt.Fprintln(w, http.StatusInternalServerError)
			return
		}
	}
}
