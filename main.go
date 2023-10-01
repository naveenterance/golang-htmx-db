package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

func main() {

	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "nst",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "recordings",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		var albums []Album

		rows, err := db.Query("SELECT * FROM album")
		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		for rows.Next() {
			var alb Album
			if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
				fmt.Printf("error")
			}

			albums = append(albums, alb)
		}

		tmpl.Execute(w, albums)

	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		title := r.PostFormValue("title")
		artist := r.PostFormValue("artist")
		price := r.PostFormValue("price")
		_, err = db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", title, artist, price)

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", Album{Title: title, Artist: artist, Price: price})

	}

	/*h3 := func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			id := r.Form.Get("id")

			db.Query("DELETE  FROM album WHERE ID = ?", id)

		}

	}*/
	/*http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		// Get the value from the request (assuming it's a GET request)
		inputValue := r.URL.Query().Get("myInput")

		// Process the value or perform any necessary operations
		valueProcessed := "Processed value: " + inputValue

		// Return the processed value as the response
		fmt.Fprint(w, valueProcessed)
	})*/

	h3 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		myVariable := r.FormValue("myVariable")
		fmt.Println("The variable value is:", myVariable)
		// Process the variable value here
		db.Query("DELETE  FROM album WHERE ID = ?", myVariable)

	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-film/", h2)
	//http.HandleFunc("/delete/", h3)
	http.HandleFunc("/process", h3)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
