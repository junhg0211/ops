package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/noirbizarre/gonja"
	"net/http"
)

func index(w http.ResponseWriter, _ *http.Request) {
	tpl, err := gonja.FromFile("templates/index.html")
	if err != nil {
		fmt.Println("importing index templates:", err)
		return
	}

	out, err := tpl.Execute(nil)
	if err != nil {
		fmt.Println("executing index templates:", err)
		return
	}

	_, err = w.Write([]byte(out))
	if err != nil {
		fmt.Println("writing index:", err)
		return
	}
}

func profile(w http.ResponseWriter, req *http.Request) {
	tpl, err := gonja.FromFile("templates/profile.html")
	if err != nil {
		fmt.Println("importing profile templates:", err)
		return
	}

	param := mux.Vars(req)
	var email = param["email"]
	var isUser = email == "junhg0211@gmail.com"
	var username string

	if isUser {
		username = "스치"
	} else {
		username = "정보 없음"
	}
	out, err := tpl.Execute(gonja.Context{
		"email":    email,
		"is_user":  isUser,
		"username": username,
	})

	if err != nil {
		fmt.Println("executing profile templates:", err)
		return
	}
	_, err = w.Write([]byte(out))
	if err != nil {
		fmt.Println("writing profile:", err)
		return
	}
}

func upload(w http.ResponseWriter, _ *http.Request) {
	tpl, err := gonja.FromFile("templates/upload.html")
	if err != nil {
		fmt.Println("importing upload templates:", err)
		return
	}

	out, err := tpl.Execute(nil)

	if err != nil {
		fmt.Println("executing upload templates:", err)
		return
	}
	_, err = w.Write([]byte(out))
	if err != nil {
		fmt.Println("writing upload:", err)
		return
	}
}

func problem(w http.ResponseWriter, req *http.Request) {
	tpl, err := gonja.FromFile("templates/problem.html")
	if err != nil {
		fmt.Println("importing problem templates:", err)
		return
	}

	param := mux.Vars(req)
	out, err := tpl.Execute(gonja.Context{
		"problem_code": param["problem_code"],
	})

	if err != nil {
		fmt.Println("executing problem templates:", err)
		return
	}
	_, err = w.Write([]byte(out))
	if err != nil {
		fmt.Println("writing problem:", err)
		return
	}
}

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	router.HandleFunc("/", index)
	router.HandleFunc("/index", index)
	router.HandleFunc("/profile/{email}", profile)
	router.HandleFunc("/upload", upload)
	router.HandleFunc("/problem/{problem_code:[A-Z]+\\d+}", problem)

	err := http.ListenAndServe(":80", router)
	if err != nil {
		fmt.Println("listening and serving:", err)
	}
}
