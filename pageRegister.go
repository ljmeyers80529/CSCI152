package csci152

import (
	"io"
	"net/http"
)

func pageRegister(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		switch req.FormValue("cmdbutton") {
		case "OK": // new user registration.
			if WriteNewUserInformation(res, req) {
				http.Redirect(res, req, "/login", http.StatusFound)
			} else {
				http.Redirect(res, req, "/register", http.StatusFound)
			}
		case "Cancel":
			http.Redirect(res, req, "/login", http.StatusFound)
		}
	}
	tpl.ExecuteTemplate(res, "register.html", nil)
}

// check if username exists.
func pageRegisterUsernameCheck(res http.ResponseWriter, req *http.Request) {
	if ex, _ := UsernameExists(req); ex {
		io.WriteString(res, "false")
		return
	}
	io.WriteString(res, "true")
}
