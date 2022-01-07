package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/noirbizarre/gonja"
	"github.com/noirbizarre/gonja/exec"
	"net/http"
	"strconv"
)

var (
	db        *sql.DB
	templates = make(map[string]*exec.Template)
)

func indexRoute(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("select * from problem order by problem_code desc limit 10")
	if err != nil {
		fmt.Println("executing querying db index:", err)
		return
	}

	var problems []problem

	for rows.Next() {
		p := problem{}
		err = rows.Scan(&p.authorEmail, &p.code, &p.content, &p.title, &p.subject)
		if err != nil {
			fmt.Println("scanning row index:", err)
			return
		}
		problems = append(problems, p)
	}

	out, err := templates["index"].Execute(gonja.Context{
		"problems": problems,
	})
	if err != nil {
		fmt.Println("executing indexRoute templates:", err)
		return
	}

	_, err = w.Write([]byte(out))
	if err != nil {
		fmt.Println("writing indexRoute:", err)
		return
	}
}

func profileRoute(w http.ResponseWriter, req *http.Request) {
	var (
		param  = mux.Vars(req)
		u      = username{email: param["authorEmail"]}
		isUser = true
	)

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
		fmt.Println("executing profileRoute templates:", err)
		return
	}
	_, err = w.Write([]byte(out))
	if err != nil {
		fmt.Println("writing profileRoute:", err)
		return
	}
}

func uploadRoute(w http.ResponseWriter, _ *http.Request) {
	out, err := templates["upload"].Execute(nil)

	if err != nil {
		fmt.Println("executing uploadRoute templates:", err)
		return
	}
	_, err = w.Write([]byte(out))
	if err != nil {
		fmt.Println("writing uploadRoute:", err)
		return
	}
}

func problemRoute(w http.ResponseWriter, req *http.Request) {
	param := mux.Vars(req)
	isProblem := true

	problemCode, err := strconv.Atoi(param["problem_code"])
	if err != nil {
		isProblem = false
	}

	p, err := getProblem(problemCode)
	if err != nil {
		isProblem = false
	}

	out, err := templates["problem"].Execute(gonja.Context{
		"is_problem":   isProblem,
		"problem_code": p.code,
		"content":      p.content,
		"title":        p.title,
		"author_email": p.authorEmail,
	})

	if err != nil {
		fmt.Println("executing problemRoute templates:", err)
		return
	}
	_, err = w.Write([]byte(out))
	if err != nil {
		fmt.Println("writing problemRoute:", err)
		return
	}
}

func reloadRoute(w http.ResponseWriter, req *http.Request) {
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
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/index", indexRoute)
	router.HandleFunc("/profile/{authorEmail}", profileRoute)
	router.HandleFunc("/upload", uploadRoute)
	router.HandleFunc("/problem/{problem_code:\\d+}", problemRoute)
	router.HandleFunc("/reload", reloadRoute)

	err = http.ListenAndServe(":80", router)
	if err != nil {
		fmt.Println("listening and serving:", err)
	}
}
