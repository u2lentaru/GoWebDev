package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/MySQL"
)

// Server - server struct
type Server struct {
	database *sql.DB
	currBlog string
}

func (server *Server) viewList(w http.ResponseWriter, r *http.Request) {
	MyBlog, err := GetBlog(server.database, server.currBlog)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.ExecuteTemplate(w, "blog", MyBlog); err != nil {
		log.Fatal(err)
	}
}

func (server *Server) viewPost(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	dbpost, err := GetPost(server.database, url[len(url)-1])
	if err != nil {
		log.Fatal(err)
	}
	if err := post.ExecuteTemplate(w, "post", dbpost); err != nil {
		log.Fatal(err)
	}
}

func (server *Server) editPost(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	dbpost, err := GetPost(server.database, url[len(url)-1])
	if err != nil {
		log.Fatal(err)
	}
	if err := edit.ExecuteTemplate(w, "edit", dbpost); err != nil {
		log.Fatal(err)
	}
}

func (server *Server) savePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res, err := server.database.Exec("update myblog.posts set subj = ?, posttime = ?, posttext = ? where id = ?",
		r.FormValue("fsubj"), r.FormValue("fpt"), r.FormValue("body"), r.FormValue("id"))

	if err != nil {
		log.Fatalf("res %v, err %v", res, err)
	}
	http.Redirect(w, r, "/", 303)
}

func (server *Server) newPost(w http.ResponseWriter, r *http.Request) {
	var indp int

	res, err := server.database.Exec("insert into myblog.posts (blogid, subj, posttime, posttext) VALUES (?,'',NOW(),'')", server.currBlog)
	if err != nil {
		log.Fatalf("err %v, res %v", err, res)
	}
	row := server.database.QueryRow("select LAST_INSERT_ID() from myblog.posts")
	err = row.Scan(&indp)
	if err != nil {
		log.Fatal(err)
	}
	newp := TPost{
		ID:       strconv.Itoa(indp),
		Subj:     "",
		PostTime: time.Now().Format("2006-01-02 15:04:05"),
		Text:     ""}
	if err := edit.ExecuteTemplate(w, "edit", newp); err != nil {
		log.Fatal(err)
	}
}

func (server *Server) delPost(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	res, err := server.database.Exec("delete from myblog.posts where id = ?", url[len(url)-1])
	if err != nil {
		log.Fatalf("err %v, res %v", err, res)
	}
	http.Redirect(w, r, "/", 303)
}

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

// MyBlog - my blog variable
var MyBlog TBlog

/* mysql> INSERT INTO blogs (name, title) VALUES (
   -> "Blog", "My Blog");
mysql> INSERT INTO posts (blogid, subj, posttime, posttext) VALUES (
   -> 1,"1st subj",NOW()-100000,"1st text"),
   -> (1,"2nd subj",NOW()-10000,"2st text"),
   -> (1,"3rd subj",NOW(),"3rd text");
*/

func main() {
	DSN := "root:qw12345@tcp(localhost:3306)/myblog?charset=utf8"
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}

	serv := Server{database: db, currBlog: "1"}
	defer serv.database.Close()

	if err := serv.database.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("db pinged!")

	router := http.NewServeMux()

	router.HandleFunc("/", serv.viewList)
	router.HandleFunc("/post/", serv.viewPost)
	router.HandleFunc("/edit/", serv.editPost)
	router.HandleFunc("/save/", serv.savePost)
	router.HandleFunc("/new/", serv.newPost)
	router.HandleFunc("/del/", serv.delPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// GetBlog - get blog from database
func GetBlog(db *sql.DB, id string) (TBlog, error) {
	blog := TBlog{}

	row := db.QueryRow("select * from myblog.blogs where blogs.id = ?", id)
	err := row.Scan(&blog.ID, &blog.Name, &blog.Title)
	if err != nil {
		return TBlog{}, err
	}

	rows, err := db.Query("select * from posts where blogid = ?", id)
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
func GetPost(db *sql.DB, id string) (TPost, error) {
	post := TPost{}
	row := db.QueryRow("select * from myblog.posts where posts.id = ?", id)
	err := row.Scan(&post.ID, new(int), &post.Subj, &post.PostTime, &post.Text)
	if err != nil {
		return TPost{}, err
	}
	return post, nil
}
