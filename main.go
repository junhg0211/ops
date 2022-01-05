package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/noirbizarre/gonja"
	"github.com/noirbizarre/gonja/exec"
	"net/http"
)

var (
	db        *sql.DB
	templates = make(map[string]*exec.Template)
)

func index(w http.ResponseWriter, _ *http.Request) {
	out, err := templates["index"].Execute(nil)
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
	param := mux.Vars(req)
	var u = username{email: param["email"]}
	var isUser = true

	row := db.QueryRow("select name from username where email = ?", u.email)
	err := row.Scan(&(u.name))
	if err != nil {
		isUser = false
	}

	out, err := templates["profile"].Execute(gonja.Context{
		"name":    u.name,
		"email":   u.email,
		"is_user": isUser,
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
	out, err := templates["upload"].Execute(nil)

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
	param := mux.Vars(req)
	out, err := templates["problem"].Execute(gonja.Context{
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

func reload(w http.ResponseWriter, req *http.Request) {
	reloadTemplates()
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

func reloadTemplates() {
	for _, templateName := range []string{"index", "profile", "upload", "problem"} {
		template, err := gonja.FromFile(fmt.Sprintf("templates/%s.html", templateName))
		templates[templateName] = template
		if err != nil {
			fmt.Println(fmt.Sprintf("importing %s templates:", templateName), err)
			return
		}
	}
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:asdf0211@tcp(sch.shtelo.org:3306)/ops")
	if err != nil {
		fmt.Println("connecting sql db:", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("closing sql db:", err)
			return
		}
	}(db)

	reloadTemplates()

	router := mux.NewRouter()
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	router.HandleFunc("/", index)
	router.HandleFunc("/index", index)
	router.HandleFunc("/profile/{email}", profile)
	router.HandleFunc("/upload", upload)
	router.HandleFunc("/problem/{problem_code:[A-Z]+\\d+}", problem)
	router.HandleFunc("/reload", reload)

	err = http.ListenAndServe(":80", router)
	if err != nil {
		fmt.Println("listening and serving:", err)
	}
}
