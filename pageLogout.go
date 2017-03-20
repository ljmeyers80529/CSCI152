package csci152

import (
	"net/http"
)

func pageLogout(res http.ResponseWriter, req *http.Request) {
	setUserDefault()
	updateCookie(res, req)
    http.Redirect(res, req, "/login", http.StatusSeeOther)
}
