package csci152

import (
	"io"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func pageHome(res http.ResponseWriter, req *http.Request) {
	readCookie(res, req)
	ctx := appengine.NewContext(req)
	log.Infof(ctx, "Received path = %s", req.URL.Path)
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
				// webInformation.Radar.Data = []int{55, 45, 11, 46, 44}
				// webInformation.Radar.Labels = []string{"Soft Rook", "Heavy Metal", "Rap", "Classical", "Adult"}
			}
		}
	}
	term := req.FormValue("term")
	switch term {
	case "long":
		term = "long_term"
	case "medium":
		term = "medium_term"
	default:
		term = "short_term"
	}
	log.Infof(ctx, "Term = %s", term)
	webInformation.User.Username = spotifyUser()
	initSpotify(res, req)
	ok := clientOK()
	if ok {
		tgl, tgs, ta, err := generateUserGenreStatistics(&spotClient, 7, term)
		if err != nil {
			log.Infof(ctx, "err: %v", err)
		} else {
			playlist, err := generateUserPlaylist(&spotClient, playlistSizeConst, term, tgl, tgs, ta)
			if err != nil {
				log.Infof(ctx, "err: %v", err)
			} else {
				webInformation.User.UserPlaylistID = string(playlist.URI)
				log.Infof(ctx, "succeed: should")
				log.Infof(ctx, string(playlist.URI))
			}
		}
		// ctx := appengine.NewContext(req)
		log.Infof(ctx, "tgl: %v", tgl)
		log.Infof(ctx, "tgs: %v", tgs)
		// log.Infof(ctx, "ta: %v", ta)
		log.Infof(ctx, "err: %v", err)
		if err == nil {
			shuffleListsInParallel(tgs, tgl)
			webInformation.User.SpotifyUsername = webInformation.User.Username
			webInformation.Radar.Data = tgs
			webInformation.Radar.Labels = tgl
		} else {
			webInformation.User.UserPlaylistID = getLoggedInUsersPlaylist(res, req)
			webInformation.User.SpotifyUsername = "Sample"
			webInformation.Radar.Data = []int{55, 45, 11, 46, 44}
			webInformation.Radar.Labels = []string{"Soft Rock", "Heavy Metal", "Rap", "Classical", "Adult"}
		}
	}
	log.Infof(ctx, "ok: %v", ok)
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
