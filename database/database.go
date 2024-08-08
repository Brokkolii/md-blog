package database

import (
	"database/sql"
	"fmt"
	"log"
	mdtohtml "md-blog-api/mdToHtml"
)

var Db *sql.DB

type Post struct {
	Content string
}

func newPost(content string) Post {
	return Post{
		Content: content,
	}
}

func InitDB() {
	var err error
	connStr := "user=mdblog password=mdblogdb dbname=mdblog sslmode=disable"
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = Db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")
}

func GetPost(id string) (Post, error) {
	var md string
	err := Db.QueryRow("SELECT content FROM posts WHERE id = $1", id).Scan(&md)
	if err != nil {
		return Post{}, err
	}
	html := mdtohtml.Parse(md)
	post := newPost(html)
	return post, nil
}

func GetPosts() ([]Post, error) {

	rows, err := Db.Query("SELECT content FROM posts")
	if err != nil {
		return []Post{}, err
	}
	defer rows.Close()

	var posts []Post = []Post{}
	for rows.Next() {
		var markdownContent string
		err := rows.Scan(&markdownContent)
		if err != nil {
			return []Post{}, err
		}

		html := mdtohtml.Parse(markdownContent)
		post := newPost(html)
		posts = append(posts, post)
	}

	return posts, nil
}

func CreatePost(content string) (Post, error) {
	Db.Exec("INSERT INTO posts (content) VALUES ($1)", content)
	post := newPost(content)
	return post, nil
}
