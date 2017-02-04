package main

import (
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"log"
	"fmt"
	"golang.org/x/oauth2/google"
	"context"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"encoding/json"
)

type Credentials struct {
	Cid 	string	`json:"cid"`
	Csecret	string	`json:"csecret"`
}

var(
	state string
	//cred Credentials
	OAuthConfig *oauth2.Config
	SessionStore sessions.Store
	cookieStore = sessions.NewCookieStore([]byte("something-very-secret"))
)

func init(){
	var c Credentials
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil{
		fmt.Printf("File err %v\n", err)
		return
	}
	json.Unmarshal(file, &c)

	OAuthConfig  = &oauth2.Config{
		ClientID: c.Cid,
		ClientSecret: c.Csecret,
		RedirectURL: "http://localhost:8080/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	//session init
	cookieStore.Options = &sessions.Options{
		HttpOnly:true,
	}
	SessionStore = cookieStore
}


func main(){
	router := httprouter.New()
	router.GET("/", indexHandler)
	router.GET("/login", loginHandler)
	router.GET("/auth", authHandler)
	http.ListenAndServe(":8080", router)
}

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	fmt.Fprintln(w, "This is index")
}

func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionID := uuid.NewV4().String()
	redirectURL := OAuthConfig.AuthCodeURL(sessionID)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func authHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	oauthFlowSession, err := SessionStore.Get(r, r.FormValue("state"))
	fmt.Println("oauthFlow session : ")
	if err != nil {
		fmt.Println("invalid state: ", oauthFlowSession)
		log.Fatalln(err)
	}

	code := r.FormValue("code")
	tok, err := OAuthConfig.Exchange(context.Background(), code)
	if err != nil{
		fmt.Println("invalid token: ")
		log.Fatalln(err)
	}

	client := OAuthConfig.Client(oauth2.NoContext, tok)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		fmt.Println("Can not retrieve email data")
		log.Fatalln(err)
	}
	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)
	log.Println("Email body: ", string(data))

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	//http.Redirect(w, r, redirectURL, http.StatusFound)
}