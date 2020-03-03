package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/MySQL"
)

// TBlog - blog struct
type TBlog struct {
	ID       string
	Name     string
	Title    string
	PostList []TPost
}

// mysql> CREATE TABLE blogs (
//    -> id SERIAL PRIMARY KEY,
//    -> name VARCHAR(255) NOT NULL,
//    -> title VARCHAR(255) NOT NULL);

// TPost - post struct
type TPost struct {
	ID       string
	Subj     string
	PostTime string
	Text     string
}

// mysql> CREATE TABLE posts (
//    -> id SERIAL PRIMARY KEY,
//    -> blogid INT NOT NULL,
//   -> subj VARCHAR(255) NOT NULL,
//    -> posttime TIMESTAMP NOT NULL,
//    -> posttext TEXT NOT NULL);

var tmpl = template.Must(template.New("MyTemplate").ParseFiles("./homework-3/tmpl.html"))
var post = template.Must(template.New("MyPost").ParseFiles("./homework-3/post.html"))
var edit = template.Must(template.New("MyPost").ParseFiles("./homework-3/edit.html"))
var database *sql.DB

// MyBlog - my blog variable
var MyBlog TBlog

/*var MyBlog = TBlog{
	Name:  "Blog",
	Title: "My blog",
	PostList: []TPost{
		TPost{"0", "1st subj", "01.01.2020", "1st text"},
		TPost{"1", "2nd subj", "02.01.2020", "2nd text"},
		TPost{"2", "3rd subj", "03.01.2020", "3rd text"},
	},
}*/

// mysql> INSERT INTO blogs (name, title) VALUES (
//    -> "Blog", "My Blog");
// mysql> INSERT INTO posts (blogid, subj, posttime, posttext) VALUES (
//    -> 1,"1st subj",NOW()-100000,"1st text"),
//    -> (1,"2nd subj",NOW()-10000,"2st text"),
//    -> (1,"3rd subj",NOW(),"3rd text");

// DSN - MySQL Data Source Name
var DSN = "root:qw12345@tcp(localhost:3306)/myblog?charset=utf8"

func main() {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	database = db
	defer database.Close()
	log.Printf("db= %v", database)

	if err := database.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("db pinged!")
	//log.Println(blog)

	router := http.NewServeMux()

	router.HandleFunc("/", viewList)
	router.HandleFunc("/post/", viewPost)
	router.HandleFunc("/edit/", editPost)
	router.HandleFunc("/save/", savePost)
	router.HandleFunc("/new/", newPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func viewList(w http.ResponseWriter, r *http.Request) {
	MyBlog, _ := GetBlog(strconv.Itoa(1))
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

// GetBlog - get blog from database
func GetBlog(id string) (TBlog, error) {
	blog := TBlog{}

	row := database.QueryRow(fmt.Sprintf("select * from myblog.blogs where blogs.id = %v", id))
	err := row.Scan(&blog.ID, &blog.Name, &blog.Title)
	if err != nil {
		return blog, err
	}

	rows, err := database.Query(fmt.Sprintf("select * from posts where blogid = %v", id))
	if err != nil {
		return blog, err
	}
	defer rows.Close()

	for rows.Next() {
		post := TPost{}
		err := rows.Scan(&post.ID, new(int), &post.Subj, &post.PostTime, &post.Text)
		if err != nil {
			log.Println(err)
			continue
		}
		blog.PostList = append(blog.PostList, post)
	}
	return blog, nil
}
