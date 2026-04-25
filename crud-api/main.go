package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	FirstName string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var film []Movie

func init() {
	//using apend methods
	film = append(film, Movie{Id: "0", Isbn: "98456", Title: "Avatar", Author: &Author{FirstName: "Jone", Lastname: "Hero"}})

	film = append(film, Movie{Id: "1", Isbn: "94456", Title: "Ava", Author: &Author{FirstName: "Captain", Lastname: "USA"}})

}

func main() {

	// fmt.Print(film)

	fmt.Println("Create, Read, Update, and Delete")

	r := mux.NewRouter()

	r.HandleFunc("/", Hello).Methods("GET")

	r.HandleFunc("/movies", getMovies).Methods("GET")

	r.HandleFunc("/movie/{id}", getMovie).Methods("GET")

	r.HandleFunc("/movie", createMovie).Methods("POST")

	r.HandleFunc("/movie/{id}", updateMovie).Methods("PUT")

	r.HandleFunc("/movie/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Serving at http://localhost:8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}

// simple function
func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
	w.Write([]byte("\nHellno"))
}

func getMovies(w http.ResponseWriter, r *http.Request) {

	// fmt.Println("serving Movies")
	w.Header().Set("Content-Type", "application/json")

	// seeding with fake value to test
	json.NewEncoder(w).Encode(film)

}

func getMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving a movie")

	// set the header to json only accepts json
	w.Header().Set("Content-Type", "application/json")

	// htt....com/user?name=".."
	prams := mux.Vars(r)

	// turning the str parms in int ,don't need when used name to serve movie
	index, _ := strconv.Atoi(prams["id"])

	// Access by index (with a safety check)
	if index >= 0 && index < len(film) {

		json.NewEncoder(w).Encode(film[index])
		return

	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Movie not found",
		"info":  "Index out of range",
	})
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Deleted !")

	w.Header().Set("Content-Type", "application/json")

	prams := mux.Vars(r)
	// id := prams["id"]

	index, _ := strconv.Atoi(prams["id"])

	// for i, Movie_data := range film {
	// 	if Movie_data.Id == id {

	// 		// Remove
	// 		film = append(film[:i], film[i+1:]...)

	// 		// Send back the updated data (or the whole film)
	// 		w.WriteHeader(http.StatusOK)
	// 		json.NewEncoder(w).Encode(film)

	// 		return

	// 	}

	// }

	if index >= 0 && index < len(film) {

		film = append(film[:index], film[index+1:]...)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(film)
		return

	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "movie not found"})
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Adding a new movie to DB!")

	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("Invalid movie_data")
	}

	// {}

	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid JSON data")
		return
	}

	if movie.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Movie title is required")
		return
	}

	// rand.Seed(time.Now().UnixNano()) //no seeding required

	movie.Id = strconv.Itoa(rand.Intn(100))

	// film is a slice
	// var film []movie

	film = append(film, movie)

	fmt.Printf("Adding a new movie: %s! Total: %d\n", movie.Title, len(film))
	w.WriteHeader(http.StatusCreated) // Use 201 Created for new resources
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updating !")

	w.Header().Set("Content-Type", "application/json")

	//updating book through book_id ,grabs id for request
	params := mux.Vars(r)
	id := params["id"]

	// looping, id ,remove ,add

	//can use if else too
	for i, movie_data := range film {
		if movie_data.Id == id {

			// Remove the old book
			film = append(film[:i], film[i+1:]...)

			// a variable to hold the NEW data
			var movie Movie

			err1 := json.NewDecoder(r.Body).Decode(&movie)

			if err1 != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// ID remains the same as the URL param
			movie.Id = id

			// Append the updated book
			film = append(film, movie)

			// Send back the updated book (or the whole film)
			json.NewEncoder(w).Encode(film)

			return

		}

	}
}
