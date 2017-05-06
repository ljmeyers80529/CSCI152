package csci152

import (
	"net/http"

	"google.golang.org/appengine"
)

func pageHome(res http.ResponseWriter, req *http.Request) {
	readCookie(res, req)
	if req.Method == "POST" {
		fn := req.FormValue("cmdbutton")
		switch fn {
		// case "Register":
		// 	http.Redirect(res, req, "/register", http.StatusSeeOther)
		case "Login":
			if checkUserLogin(res, req) {
				initSpotify(res, req)
			}
		}
	}
	webInformation.User.Username = spotifyUser()
	if clientOK() {
		webInformation.User.UserPlaylistID = loadPlayLists(res, req)
	}
	tpl.ExecuteTemplate(res, "homepage.html", webInformation)
}

// check if user has successfully logged in.
// returns true if success.
func checkUserLogin(res http.ResponseWriter, req *http.Request) bool {
	var uuidKey string

	ctx := appengine.NewContext(req)
	user := req.FormValue("username")
	pass := EncryptPassword(req.FormValue("password"))

	if uuidKey = SearchUser(ctx, user); uuidKey != "" {
		ReadUserInformation(ctx, req, uuidKey)
		userInformation.LoggedIn = userInformation.Password == pass
		updateCookie(res, req)
	}
	return userInformation.LoggedIn
}
