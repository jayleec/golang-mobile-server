package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func main(){
	router := httprouter.New()
	router.GET("/", indexHandler)
	router.GET("/login", loginHandler)
	router.GET("/auth", authHandler)
	http.ListenAndServe(":8080", router)
}
