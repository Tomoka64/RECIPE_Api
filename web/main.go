package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("running yeahh")
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

// func main() {
// 	fmt.Println("running yeahh")
// 	http.HandleFunc("/", hello)
// 	http.ListenAndServe(":8080", nil)
// }
//
// func hello(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "hello world")
// }

//
// func main() {
// 	router := httprouter.New()
// 	http.Handle("/", router)
// 	router.GET("/", Index)
// 	router.GET("/recipes", List)
// 	router.GET("/recipes/create", Create)
// 	router.GET("/recipes/{id}", Get)
// 	router.GET("/recipes/{id}/update", Update)
// 	router.GET("/recipes/{id}/delete", Delete)
// 	router.POST("/recipes/{id}/rating", Rate)
//
// 	router.GET("/form/login", Login)
// 	router.GET("/form/signup", Signup)
//
// 	// router.POST("/api/checkusername", checkUserName)
// 	// router.POST("/api/createuser", createUser)
// 	// router.POST("/api/login", loginProcess)
// 	// router.GET("/api/logout", logout)
// 	// router.GET("/post/done", Done)
//
// 	server := &http.Server{
// 		Addr:         ":8080",
// 		Handler:      router,
// 		ReadTimeout:  time.Duration(20 * time.Second),
// 		WriteTimeout: time.Duration(20 * time.Second),
// 	}
//
// 	server.ListenAndServe()
// }
