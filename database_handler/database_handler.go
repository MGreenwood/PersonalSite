package database_handler

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	Title     string
	Timestamp string
	Body      string
	Image     string
}

func Add(title string, timestamp string, body string, image string) error {
	pass := os.Getenv("DB_PASS")
	connection_string := "root:" + pass + "@/main"
	conn, err := sql.Open("mysql", connection_string)
	if err != nil {
		return err
	}
	defer conn.Close()

	statement := fmt.Sprintf(`insert into posts (Title, Timestamp, Body, Image) values ('%s', '%s', '%s', '%s')`,
		title, timestamp, body, image)
	_, err = conn.Exec(statement)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func RetrieveTopFive() []Post {
	pass := os.Getenv("DB_PASS")
	connection_string := "root:" + pass + "@/main"
	conn, err := sql.Open("mysql", connection_string)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	statement := "select * from posts order by id desc limit 5"
	query, err := conn.Query(statement)
	if err != nil {
		fmt.Println(err)
	}

	posts := []Post{}
	for query.Next() {
		var r Post
		var id int
		err = query.Scan(&r.Title, &r.Timestamp, &r.Body, &r.Image, &id)
		if err != nil {
			fmt.Println(err)
		}
		decodedBody, _ := base64.URLEncoding.DecodeString(r.Body)
		r.Body = string(decodedBody)
		posts = append(posts, r)
	}

	return posts
}
