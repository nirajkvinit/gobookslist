package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
)

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:title`
	Year   string `json:year`
}

var books []Book
var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	logFatal(err)
	// log.Fatal(err)
	log.Println(pgURL)

	db, err = sql.Open("postgres", pgURL)
	logFatal(err)
	// log.Fatal(err)

	err = db.Ping()
	logFatal(err)

	router := mux.NewRouter()

	// books = append(books,
	// 	Book{ID: 1, Title: "Golang Pointers", Author: "Mr. Golang", Year: "2010"},
	// 	Book{ID: 2, Title: "Goroutines", Author: "Mr. GoRouting", Year: "2011"},
	// 	Book{ID: 3, Title: "Golang Routers", Author: "Mr. Router", Year: "2012"},
	// 	Book{ID: 4, Title: "Golang concurrency", Author: "Mr. Currency", Year: "2013"},
	// 	Book{ID: 5, Title: "Go good parts", Author: "Mr. Good", Year: "2014"},
	// )

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var book Book
	books = []Book{}

	rows, err := db.Query("select * from books")
	logFatal(err)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)

		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	params := mux.Vars(r)
	// book_ID, _ := strconv.Atoi(params["id"])
	rows := db.QueryRow("select * from books where id=$1", params["id"])

	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)

	logFatal(err)

	json.NewEncoder(w).Encode(bookID)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	// var bookID int

	json.NewDecoder(r.Body).Decode(&book)

	result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id;",
		&book.Title, &book.Author, &book.Year, &book.ID)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	result, err := db.Exec("delete from books where id = $1;", params["id"])

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}

/**
func getBooks(w http.ResponseWriter, r *http.Request) {
	// log.Println("Gets all books")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	// log.Println("Get one book")
	params := mux.Vars(r)
	// log.Println(params)
	// log.Println(reflect.TypeOf(params["id"]))
	bookID, _ := strconv.Atoi(params["id"])

	for _, book := range books {
		if book.ID == bookID {
			json.NewEncoder(w).Encode(&book)
		}
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	// log.Println("Add a book")
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	books = append(books, book)

	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	// log.Println("Update a book")
	var book Book

	json.NewDecoder(r.Body).Decode(&book)

	for i, item := range books {
		if item.ID == book.ID {
			books[i] = book
		}
	}

	json.NewEncoder(w).Encode(books)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	// log.Println("Remove a book")
	params := mux.Vars(r)

	bookID, _ := strconv.Atoi(params["id"])

	for i, item := range books {
		if item.ID == bookID {
			books = append(books[:i], books[i+1:]...)
		}
	}

	json.NewEncoder(w).Encode(books)
}
*/
