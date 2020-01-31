package main

import (
	"encoding/gob"
	"fmt"
	"github-oauth/app/base"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	ClientId   string
	RequestUri string
}

var M map[string]interface{}

func init() {
	gob.Register(&base.OAuthAccessResponse{})
	gob.Register(&M)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/oauth/redirect", redirectHandler)
	http.HandleFunc("/welcome", welcomeHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, isLogin := base.CheckLogin(r)
	if isLogin {
		http.Redirect(w, r, "/welcome", 302)
		return
	}
	data := PageData{
		ClientId:   base.ClientId,
		RequestUri: base.RequestUri,
	}
	t := template.Must(template.ParseFiles("app/templates/index.html"))
	_ = t.Execute(w, data)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]
	if ok && len(keys[0]) > 1 {
		code := keys[0]
		res, err := base.Auth(code)
		if err != nil {
			fmt.Println("err")
		} else {
			sessionStore := base.SessionStore
			session, err := sessionStore.New(r, "user")
			if err != nil {
				log.Fatalln("new session err-->", err)
			}
			session.Values["auth"] = res
			err = sessionStore.Save(r, w, session)
			if err != nil {
				http.Error(w,"session 存储失败",500)
				return
			}
			http.Redirect(w, r, "/welcome", 302)
		}
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	_,isLogin := base.CheckLogin(r)
	if !isLogin {
		http.Redirect(w,r,"/",302)
		return
	}
	t, err := template.ParseFiles("app/templates/welcome.html")
	if err != nil {
		log.Fatalln("render template err", err)
		return
	}
	user,err := base.FetchUserInfo(r,w)
	if err != nil {
		http.Redirect(w,r,"/",302)
		return
	}
	_ = t.Execute(w, user)
}
