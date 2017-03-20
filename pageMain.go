package csci152

import (
	// "fmt"
	"net/http"
	// "strings"
	// "strconv"

	// "google.golang.org/appengine"
	// "google.golang.org/appengine/log"
)

// main (top) web page.
func pageMain(res http.ResponseWriter, req *http.Request) {
	// ctx := appengine.NewContext(req)
	readCookie(res, req) // maintain user login / out state.
	if webInformation.User.LoggedIn {
		http.Redirect(res, req, "/home", http.StatusSeeOther)
	}
}
