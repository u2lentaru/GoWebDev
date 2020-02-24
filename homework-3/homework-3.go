package main

import (
	"html/template"
	"log"
	"net/http"
)

//TPost - post struct
type TPost struct {
	Subj     string
	PostTime string
	Text     string
}

//TBlog - blog struct
type TBlog struct {
	Title    string
	PostList []TPost
}

var tmpl = template.Must(template.New("MyTemplate").ParseFiles("tmpl.html"))

//MyBlog - my blog variable
var MyBlog = TBlog{
	Title: "My blog",
	PostList: []TPost{
		TPost{"1st subj", "01.01.2020", "1st text"},
		TPost{"2nd subj", "02.01.2020", "2nd text"},
		TPost{"3rd subj", "03.01.2020", "3rd text"},
	},
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", viewList)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func viewList(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.ExecuteTemplate(w, "blog", MyBlog); err != nil {
		log.Println(err)
	}
}
