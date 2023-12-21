package helpers

import (
	"os"
	"strings"
)

func EnforceHTTP(url string) string {
	if url[:4] != "http" { //means if url contains text like goolg.com -> it attaches http://google.com ok!
		return "http://" + url
	}
	return url

}
func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") { //if it is localhost:3000 then refuse/false
		return false
	}
	newURL := strings.Replace(url, "http://", " ", 1) //if url contains these-> http/https etc then remove all these and return false
	newURL = strings.Replace(newURL, "https://", " ", 1)
	newURL = strings.Replace(newURL, "www.", " ", 1)
	newURL = strings.Split(newURL, "/")[0] //select the first element

	if newURL == os.Getenv("Domain") {
		return false
	}

	return true //if doesn't contains these all then return true

}
