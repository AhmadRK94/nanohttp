package main

import (
	"fmt"
	"net/url"
	"slices"
)

var ValidHttpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

func main() {
	s := "/search?q=golang&page=2"
	parsedURL, _ := url.ParseRequestURI(s)
	fmt.Println(parsedURL)
	fmt.Println(parsedURL.Path)
	fmt.Println(parsedURL.RawPath)
	fmt.Println(parsedURL.Query())
	fmt.Println(slices.Contains(ValidHttpMethods, "GET"))
	fmt.Println(len("سلام"))

}
