package main

import (
	"fmt"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	// 2. The form submission route
	http.HandleFunc("/form", formHandler)

	// 3. The redirect route
	http.HandleFunc("/r/", redirectHandler)

	fmt.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", nil)
}
