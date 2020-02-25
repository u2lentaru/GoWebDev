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
	Name     string
	Title    string
	PostList []TPost
}

var tmpl = template.Must(template.New("MyTemplate").ParseFiles("./homework-3/tmpl.html"))
var post = template.Must(template.New("MyPost").ParseFiles("./homework-3/edit.html"))
var editpost = template.Must(template.New("MyEditPost").ParseFiles("./homework-3/edit.html"))

//MyBlog - my blog variable
var MyBlog = TBlog{
	Name:  "Blog",
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
	router.HandleFunc("/post", viewPost)
	router.HandleFunc("/edit", editPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func viewList(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.ExecuteTemplate(w, "blog", MyBlog); err != nil {
		log.Println(err)
	}
}

func viewPost(w http.ResponseWriter, r *http.Request) {
	if err := post.ExecuteTemplate(w, "post", MyBlog.PostList[0]); err != nil {
		log.Println(err)
	}
}

func editPost(w http.ResponseWriter, r *http.Request) {
	if err := post.ExecuteTemplate(w, "editpost", MyBlog.PostList[0]); err != nil {
		log.Println(err)
	}
}
