package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Link struct {
	OriginalURL string
	Clicks      int
}

type URLShortener struct {
	mu    sync.Mutex
	store map[string]*Link // Storing pointers to our Link struct
}

// Encapsulation
func (s *URLShortener) generateCode(n int) string {

	//Create a "slice" (list) of bytes with length 'n'
	b := make([]byte, n)

	// write completely random numbers
	rand.Read(b)

	//Turn those random numbers into a "Base64" string
	// Base64 turns numbers into letters and symbols (A-Z, a-z, 0-9)
	// We then cut it to the length 'n'
	return base64.URLEncoding.EncodeToString(b)[:n]

}

var shortener = &URLShortener{store: make(map[string]*Link)}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ERROR: %v", err)
		return
	}

	// Get the URL from your HTML form input (name="url")
	urlGet := r.FormValue("url")
	if urlGet == "" {
		fmt.Fprint(w, "Please enter a URL")
		return
	}

	//

	// Logic: Shorten and Store
	shortener.mu.Lock() //wait -> cocurrency

	code := shortener.generateCode(10)
	shortener.store[code] = &Link{OriginalURL: urlGet, Clicks: 0}

	shortener.mu.Unlock() //I'm done! ->cocurrency

	// Show the result to the user
	fmt.Fprintf(w, "Post sent successfully!\n")
	fmt.Fprintf(w, "To send another post : http://localhost:8080/form.html\n")
	fmt.Fprintf(w, "Original URL: %s\n", urlGet)
	fmt.Fprintf(w, "Your short link: http://localhost:8080/r/%s", code)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract code from URL (e.g., /r/abc123 -> abc123)
	code := strings.TrimPrefix(r.URL.Path, "/r/")

	shortener.mu.Lock() //wait -> cocurrency
	link, exists := shortener.store[code]
	if exists {
		link.Clicks++ // Increment the counter!
	}
	shortener.mu.Unlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	fmt.Printf("Redirecting to %s. Total clicks: %d\n", link.OriginalURL, link.Clicks)
	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}
