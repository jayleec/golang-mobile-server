package main

import (
	"github.com/gorilla/sessions"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"io/ioutil"
	"golang.org/x/oauth2"
	"context"
	"encoding/json"
	"golang.org/x/oauth2/google"
)

type Credentials struct {
	Cid 	string	`json:"cid"`
	Csecret	string	`json:"csecret"`
}

var(
	state string
	OAuthConfig *oauth2.Config
	SessionStore sessions.Store
	cookieStore = sessions.NewCookieStore([]byte("something-very-secret"))
)

func init(){
	//read oAuth Client info from json file
	var c Credentials
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil{
		fmt.Printf("File err %v\n", err)
		return
	}
	json.Unmarshal(file, &c)
	//oAuth config init
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

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	fmt.Fprintln(w, "This is index")
}

func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionID := uuid.NewV4().String()
	//use session ID for state params, protects against CSRF
	redirectURL := OAuthConfig.AuthCodeURL(sessionID)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

//oAuth callback handler
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
}
