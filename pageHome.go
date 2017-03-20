package csci152

import (
	"net/http"

	// "google.golang.org/appengine"
)

func pageHome(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
	}
	tpl.ExecuteTemplate(res, "homepage.html", nil)
}