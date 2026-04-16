package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloFunc)

	fmt.Println("Server starting localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		return
	}
	fmt.Fprintf(w, "hello")

}

func formHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "ERROR :%v", err)
		return
	}
	fmt.Fprintf(w, "POST send")

	first_name := r.FormValue("firstname")
	last_name := r.FormValue("lastname")
	addr := r.FormValue("address")

	fmt.Fprintf(w, "Name : %s %s \n", first_name, last_name)
	fmt.Fprintf(w, "Addr : %s", addr)

}
