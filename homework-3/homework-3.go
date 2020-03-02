package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/MySQL"
)

//TPost - post struct
type TPost struct {
	ID       string
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
var post = template.Must(template.New("MyPost").ParseFiles("./homework-3/post.html"))
var edit = template.Must(template.New("MyPost").ParseFiles("./homework-3/edit.html"))

//MyBlog - my blog variable
var MyBlog = TBlog{
	Name:  "Blog",
	Title: "My blog",
	PostList: []TPost{
		TPost{"0", "1st subj", "01.01.2020", "1st text"},
		TPost{"1", "2nd subj", "02.01.2020", "2nd text"},
		TPost{"2", "3rd subj", "03.01.2020", "3rd text"},
	},
}

//DSN - MySQL Data Source Name
var DSN = "root:12345@tcp(localhost:3306)/shop?charset=utf8"

func main() {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Printf("db= %v", db)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("db pinged!")
	}

	router := http.NewServeMux()

	router.HandleFunc("/", viewList)
	router.HandleFunc("/post/", viewPost)
	router.HandleFunc("/edit/", editPost)
	router.HandleFunc("/save/", savePost)
	router.HandleFunc("/new/", newPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func viewList(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.ExecuteTemplate(w, "blog", MyBlog); err != nil {
		log.Println(err)
	}
}

func viewPost(w http.ResponseWriter, r *http.Request) {
	indp, err := strconv.ParseInt(r.URL.Path[len("/post/"):], 10, 16)
	if err != nil {
		log.Println(err)
		return
	}

	if err := post.ExecuteTemplate(w, "post", MyBlog.PostList[indp]); err != nil {
		log.Println(err)
	}
}

func editPost(w http.ResponseWriter, r *http.Request) {
	indp, err := strconv.ParseInt(r.URL.Path[len("/edit/"):], 10, 16)
	if err != nil {
		log.Println(err)
		return
	}

	if err := edit.ExecuteTemplate(w, "edit", MyBlog.PostList[indp]); err != nil {
		log.Println(err)
	}
}

func savePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	indp, err := strconv.ParseInt(string(r.FormValue("id")), 10, 16)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", 500)
	}
	MyBlog.PostList[indp].PostTime = r.FormValue("fpt")
	MyBlog.PostList[indp].Subj = r.FormValue("fsubj")
	MyBlog.PostList[indp].Text = r.FormValue("body")

	http.Redirect(w, r, "/", 303)
}

func newPost(w http.ResponseWriter, r *http.Request) {
	indp := len(MyBlog.PostList)
	newp := TPost{strconv.Itoa(indp), "", "", ""}

	MyBlog.PostList = append(MyBlog.PostList, newp)

	if err := edit.ExecuteTemplate(w, "edit", MyBlog.PostList[indp]); err != nil {
		log.Println(err)
	}

}
