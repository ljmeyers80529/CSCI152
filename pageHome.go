package csci152

import (
	"net/http"

	// "google.golang.org/appengine"
)

func pageHome(res http.ResponseWriter, req *http.Request) {
	readCookie(res, req)
	if req.Method == "POST" {
	}
	tpl.ExecuteTemplate(res, "homepage.html", webInformation)
}
