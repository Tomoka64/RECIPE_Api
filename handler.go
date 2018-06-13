package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Tomoka64/RECIPE_Api/internal/middleware/logger"
	"github.com/Tomoka64/RECIPE_Api/internal/middleware/session"
	"github.com/Tomoka64/RECIPE_Api/internal/router"
	"github.com/Tomoka64/RECIPE_Api/internal/status"
	"github.com/gorilla/sessions"

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

type User struct {
	Name     string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handler() error {
	r := router.New(s.Logger)
	r.UseMiddleWare(
		logger.MiddleWare(s.Logger),
		session.MiddleWare("session", s.Store),
	)
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

	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Index() router.HandlerFunc {
	return func(c router.Context) error {
		return c.Redirect(status.SeeOther, "/recipes")
	}
}

func (s *Server) List() router.HandlerFunc {
	return func(c router.Context) error {
		recipes, err := postgres.GetAllRecipes(s.DB)
		if err != nil {
			return router.NewHTTPError(status.NotFound, err.Error())
		}
		return c.JSON(status.OK, &Data{Recipes: recipes})
	}
}

func (s *Server) LoginProcess() router.HandlerFunc {
	return func(c router.Context) error {
		v, ok := c.Get("session")
		if ok {
			sess := v.(*sessions.Session)
			_, ok := sess.Values["username"]
			if ok {
				return c.Redirect(status.SeeOther, "/")
			}
		}
		username := c.Request().FormValue("username")
		email := c.Request().FormValue("email")
		password := c.Request().FormValue("password")
		if !s.UserExists(username, email) {
			return router.NewHTTPError(status.NotFound, "your username or email does not exist")
		}
		if s.UserAndPasswordMatch(username, password) {
			c.Set("session", username)
		}
		return nil
	}
}

func (s *Server) Get() router.HandlerFunc {
	return func(c router.Context) error {
		recipeId := c.Params().ByName("id")
		recipe, err := postgres.GetRecipeByID(s.DB, recipeId)
		if err != nil {
			return router.NewHTTPError(status.NotFound, err.Error())
		}
		return nil
	}
}

func (s *Server) Create() router.HandlerFunc {
	return func(c router.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")
		v, ok := c.Get("session")
		if !ok {
			return c.Redirect(status.SeeOther, "/form/login")
		}
		var recipe postgres.Recipe
		recipeId, err := strconv.Atoi(c.Params().ByName("id"))
		if err != nil {
			return router.NewHTTPError(status.Unauthorized, err.Error())
		}
		recipe.ID = recipeId
		if err := json.NewDecoder(c.Request().Body).Decode(&recipe); err != nil {
			return router.NewHTTPError(status.InternalServerError, err.Error())
		}
		if err := recipe.CreateRecipe(s.DB); err != nil {
			return router.NewHTTPError(status.InternalServerError, err.Error())

		}
		return nil
	}
}

func (s *Server) Update() router.HandlerFunc {
	return func(c router.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")
		v, ok := c.Get("session")
		if !ok {
			return c.Redirect(status.SeeOther, "/form/login")
		}
		recipeId, err := strconv.Atoi(c.Params().ByName("id"))
		if err != nil {
			return router.NewHTTPError(status.Unauthorized, err.Error())
		}
		var recipe postgres.Recipe
		sess := v.(*sessions.Session)
		u, err := postgres.GetAuthorByRecipeID(s.DB, c.Params().ByName("id"))
		if err != nil {
			return router.NewHTTPError(status.NotFound, err.Error())

		}
		if sess.Values["username"] == u.Name {
			if err := recipe.DeleteRecipe(s.DB); err != nil {
				return router.NewHTTPError(status.InternalServerError, err.Error())
			}
		}
		recipe.ID = recipeId
		if err := recipe.UpdateRecipe(s.DB); err != nil {
			return router.NewHTTPError(status.InternalServerError, err.Error())
		}
		return nil
	}
}

func (s *Server) Delete() router.HandlerFunc {
	return func(c router.Context) error {
		v, ok := c.Get("session")
		if !ok {
			return c.Redirect(status.SeeOther, "/form/login")
		}
		var recipe postgres.Recipe
		sess := v.(*sessions.Session)
		u, err := postgres.GetAuthorByRecipeID(s.DB, c.Params().ByName("id"))
		if err != nil {
			return router.NewHTTPError(status.NotFound, err.Error())

		}
		if sess.Values["username"] == u.Name {
			if err := recipe.DeleteRecipe(s.DB); err != nil {
				return router.NewHTTPError(status.InternalServerError, err.Error())
			}
		}
		return nil
	}
}

func (s *Server) Rate() router.HandlerFunc {
	return func(c router.Context) error {
		recipiId := c.Params().ByName("id")
		if err := postgres.UpdateRate(s.DB, recipiId); err != nil {
			return router.NewHTTPError(status.InternalServerError, err.Error())
		}
		return router.NewHTTPError(status.SeeOther, fmt.Sprintf("/recipes/%s", recipiId))
	}
	return nil
}

func (s *Server) Login() router.HandlerFunc {
	return func(c router.Context) error {
		_, ok := c.Get("session")
		if ok {
			return c.Redirect(status.SeeOther, "/")
		}
		return c.JSON(status.OK, &User{})
	}
}

func (s *Server) Logout() router.HandlerFunc {
	return func(c router.Context) error {
		v, ok := c.Get("session")
		if ok {
			return c.Redirect(status.SeeOther, "/")
		}
		sess := v.(*sessions.Session)
		sess.Options.MaxAge = -1
		sess.Save(c.Request(), c.Response())
		return c.Redirect(status.SeeOther, "/recipes")
	}
	return nil
}

func (s *Server) Signup() router.HandlerFunc {
	return func(c router.Context) error {
		_, ok := c.Get("session")
		if ok {
			return c.Redirect(status.SeeOther, "/")
		}
		return c.JSON(status.OK, &User{})
	}
	return nil
}

func (s *Server) CreateUser() router.HandlerFunc {
	return func(c router.Context) error {
		var user *User
		if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
			return router.NewHTTPError(status.InternalServerError, err.Error())
		}
		if !IsPasswordValid(user.Password) {
			return c.JSON(status.BadRequest, &user)
		}
		if s.UserExists(user.Name, user.Email) {
			return c.JSON(status.Found, &User{})
		}
		var pUser postgres.User
		pUser.Name, pUser.Password, pUser.Email = user.Name, user.Password, user.Email
		if err := pUser.CreateUser(s.DB); err != nil {
			return router.NewHTTPError(status.InternalServerError, err.Error())
		}
		return nil
	}
}
