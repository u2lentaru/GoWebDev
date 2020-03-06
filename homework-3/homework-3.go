package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/MySQL"
)

// TBlog - blog struct
type TBlog struct {
	ID       string
	Name     string
	Title    string
	PostList []TPost
}

/* mysql> CREATE TABLE blogs (
   -> id SERIAL PRIMARY KEY,
   -> name VARCHAR(255) NOT NULL,
   -> title VARCHAR(255) NOT NULL);
*/

// TPost - post struct
type TPost struct {
	ID       string
	Subj     string
	PostTime string
	Text     string
}

/* mysql> CREATE TABLE posts (
    -> id SERIAL PRIMARY KEY,
    -> blogid INT NOT NULL,
   -> subj VARCHAR(255) NOT NULL,
    -> posttime TIMESTAMP NOT NULL,
    -> posttext TEXT NOT NULL);
*/

var tmpl = template.Must(template.New("MyTemplate").ParseFiles("./homework-3/tmpl.html"))
var post = template.Must(template.New("MyPost").ParseFiles("./homework-3/post.html"))
var edit = template.Must(template.New("MyPost").ParseFiles("./homework-3/edit.html"))
var database *sql.DB

// MyBlog - my blog variable
var MyBlog TBlog

/* mysql> INSERT INTO blogs (name, title) VALUES (
   -> "Blog", "My Blog");
mysql> INSERT INTO posts (blogid, subj, posttime, posttext) VALUES (
   -> 1,"1st subj",NOW()-100000,"1st text"),
   -> (1,"2nd subj",NOW()-10000,"2st text"),
   -> (1,"3rd subj",NOW(),"3rd text");
*/

// DSN - MySQL Data Source Name
var DSN = "root:qw12345@tcp(localhost:3306)/myblog?charset=utf8"

func main() {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	database = db
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("db pinged!")

	router := http.NewServeMux()

	router.HandleFunc("/", viewList)
	router.HandleFunc("/post/", viewPost)
	router.HandleFunc("/edit/", editPost)
	router.HandleFunc("/save/", savePost)
	router.HandleFunc("/new/", newPost)
	router.HandleFunc("/del/", delPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func viewList(w http.ResponseWriter, r *http.Request) {
	MyBlog, err := GetBlog(strconv.Itoa(1))
	if err != nil {
		log.Println(err)
	}
	if err := tmpl.ExecuteTemplate(w, "blog", MyBlog); err != nil {
		log.Println(err)
	}
}

func viewPost(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	dbpost, err := GetPost(url[len(url)-1])
	if err != nil {
		log.Println(err)
	}
	if err := post.ExecuteTemplate(w, "post", dbpost); err != nil {
		log.Println(err)
	}
}

func editPost(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	dbpost, err := GetPost(url[len(url)-1])
	if err != nil {
		log.Println(err)
	}
	if err := edit.ExecuteTemplate(w, "edit", dbpost); err != nil {
		log.Println(err)
	}
}

func savePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res, err := database.Exec(fmt.Sprintf("update myblog.posts set subj='%v', posttime='%v', posttext='%v' where id='%v'",
		r.FormValue("fsubj"), r.FormValue("fpt"), r.FormValue("body"), r.FormValue("id")))

	if err != nil {
		log.Printf("res %v, err %v", res, err)
	}

	http.Redirect(w, r, "/", 303)
}

func newPost(w http.ResponseWriter, r *http.Request) {
	var indp int

	res, err := database.Exec(fmt.Sprintf("insert into myblog.posts (blogid, subj, posttime, posttext) VALUES (1,'',NOW(),'')"))
	if err != nil {
		log.Printf("err %v, res %v", err, res)
	}
	row := database.QueryRow(fmt.Sprintf("select LAST_INSERT_ID() from myblog.posts"))
	err = row.Scan(&indp)
	if err != nil {
		log.Println(err)
	}
	newp := TPost{ID: strconv.Itoa(indp), Subj: "", PostTime: time.Now().Format("2006-01-02 15:04:05"), Text: ""}
	if err := edit.ExecuteTemplate(w, "edit", newp); err != nil {
		log.Println(err)
	}
}

func delPost(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	res, err := database.Exec("delete from myblog.posts where id = ?", url[len(url)-1])
	if err != nil {
		log.Printf("err %v, res %v", err, res)
	}
	http.Redirect(w, r, "/", 303)
}

// GetBlog - get blog from database
func GetBlog(id string) (TBlog, error) {
	blog := TBlog{}

	row := database.QueryRow(fmt.Sprintf("select * from myblog.blogs where blogs.id = %v", id))
	err := row.Scan(&blog.ID, &blog.Name, &blog.Title)
	if err != nil {
		return TBlog{}, err
	}

	rows, err := database.Query(fmt.Sprintf("select * from posts where blogid = %v", id))
	if err != nil {
		return TBlog{}, err
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

// GetPost - get post from database
func GetPost(id string) (TPost, error) {
	post := TPost{}

	row := database.QueryRow(fmt.Sprintf("select * from myblog.posts where posts.id = %v", id))
	err := row.Scan(&post.ID, new(int), &post.Subj, &post.PostTime, &post.Text)
	if err != nil {
		return TPost{}, err
	}
	//log.Println(post)
	return post, nil
}
