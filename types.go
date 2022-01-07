package main

type username struct {
	email string
	name  string
}

type problem struct {
	authorEmail string
	code        int
	content     string
	title       string
	subject     string
}

func getProblem(problemCode int) (p problem, err error) {
	p.code = problemCode
	row := db.QueryRow("select email, content, title, subject from problem where problem_code = ?", p.code)
	err = row.Scan(&p.authorEmail, &p.content, &p.title, &p.subject)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (p problem) shortenedContent() string {
	return p.content[:30]
}
