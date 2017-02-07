package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
	"html/template"
)

var chatTemplate = template.Must(template.ParseFiles("chat.html"))

func serveChat(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	log.Println(r.URL)
	if r.URL.Path != "/chat" {
		http.Error(w, "Not Found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset-utf-8")
	chatTemplate.Execute(w, r.Host)
}

func main(){
	hub := newHub()
	go hub.run()
	router := httprouter.New()
	router.GET("/", indexHandler)
	router.GET("/login", loginHandler)
	router.GET("/auth", authHandler)
	router.GET("/chat", serveChat)
	router.GET("/ws", func (w http.ResponseWriter, r *http.Request, _ httprouter.Params){
		serveWs(hub, w, r)
	})
	http.ListenAndServe(":8080", router)
}
