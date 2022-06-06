package database_handler

import (
	"database/sql"
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

	statement := "select * from posts limit 5"
	query, err := conn.Query(statement)
	if err != nil {
		fmt.Println(err)
	}

	posts := []Post{}
	for query.Next() {
		var r Post
		err = query.Scan(&r.Title, &r.Timestamp, &r.Body, &r.Image)
		if err != nil {
			fmt.Println(err)
		}
		posts = append(posts, r)
	}

	return posts
}
