package csci152

import (
	"io"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func pageHome(res http.ResponseWriter, req *http.Request) {
	readCookie(res, req)
	if req.Method == "POST" {
		fn := req.FormValue("cmdbutton")
		switch fn {
		case "OK":
			if WriteNewUserInformation(res, req) {
				http.Redirect(res, req, "/home", http.StatusFound)
			}
		case "Cancel":
			http.Redirect(res, req, "/home", http.StatusFound)
		case "Login":
			if checkUserLogin(res, req) {
				initSpotify(res, req)
				webInformation.Radar.Data = []int{55, 45, 11, 46, 44}
				webInformation.Radar.Labels = []string{"Soft Rook", "Heavy Metal", "Rap", "Classical", "Adult"}
			}
		}
	}
	webInformation.User.Username = spotifyUser()
	if clientOK() {
		ctx := appengine.NewContext(req)
		webInformation.User.UserPlaylistID = loadPlayLists(res, req)
		tgl, tgs, ta, err := generateUserGenreStatistics(&spotClient, 10, "short_term")
		log.Infof(ctx, "tgl: %v", tgl)
		log.Infof(ctx, "tgs: %v", tgs)
		log.Infof(ctx, "ta: %v", ta)
		log.Infof(ctx, "err: %v", err)
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

// check if username exists.
func pageRegisterUsernameCheck(res http.ResponseWriter, req *http.Request) {
	if ex, _ := UsernameExists(req); ex {
		io.WriteString(res, "false")
		return
	}
	io.WriteString(res, "true")
}
