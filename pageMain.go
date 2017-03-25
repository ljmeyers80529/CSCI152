package csci152

import (
	"net/http"
)

// main (top) web page.
func pageMain(res http.ResponseWriter, req *http.Request) {
	readCookie(res, req) // maintain user login / out state.
	if webInformation.User.LoggedIn {
		http.Redirect(res, req, "/home", http.StatusSeeOther)
	}
}
