package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book struct
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Status *Status `json:"status"`
}

// Status struct
type Status struct {
	Onshelf bool  `json:"onshelf"`
	Loan    *Loan `json:"loan"`
}

// Loan struct
type Loan struct {
	Borrower string `json:"borrower"`
	Phone    string `json:"phone"`
	Duedate  string `json:"duedate"`
}

// Init books var as a slice Book struct
var books []Book

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if book.Status.Onshelf == true && book.Status.Loan != nil {
		msg := fmt.Sprintf("Error. When book on shelf, loan must be null")
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	book.ID = createID()
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			var book Book
			err := json.NewDecoder(r.Body).Decode(&book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if book.Status.Onshelf == true && book.Status.Loan != nil {
				msg := fmt.Sprintf("Error. When book on shelf, loan must be null")
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			books = append(books[:index], books[index+1:]...)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

// Delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	// Init router
	r := mux.NewRouter()

	// Mock data
	books = append(books, Book{ID: "14215", Isbn: "0679405429", Title: "Pride and Prejudice", Author: "Jane Austen", Status: &Status{Onshelf: true}})
	books = append(books, Book{ID: "25743", Isbn: "0141439564", Title: "Great Expectations", Author: "Charles Dickens", Status: &Status{Onshelf: false, Loan: &Loan{Borrower: "John Doe", Phone: "860-482-8532", Duedate: "2020-05-20"}}})

	// Router handlers / endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}

// Helper function to create unique ID
func createID() string {
	tryID := ""
	unique := false
	for unique == false {
		unique = true
		tryID = strconv.Itoa(rand.Intn(10000000))
		for _, item := range books {
			if item.ID == tryID {
				unique = false
				break
			}
		}
	}
	return tryID
}
