package csci152

import (
	"html/template"
	// "log"
	"net/http"

	"github.com/zmb3/spotify"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	configureResourceLocation("images", "img")
	configureResourceLocation("css", "css")
	configureResourceLocation("images", "js/images")
	configureResourceLocation("js", "js")
	setUserDefault()
	// http.Handle("/favicon.ico", http.NotFoundHandler()) // ignore favicon request (error 404)

	auth = spotify.NewAuthenticator(retrieveURI, spotify.ScopeUserReadPrivate)
	auth.SetAuthInfo(clientID, spotKey)
	baseURL = auth.AuthURL(spotStateValue)

	http.HandleFunc("/", pageMain) // main page.
	http.HandleFunc("/callback", completeAuth)

	tpl = template.Must(template.ParseGlob("html/*.html"))
}

// map resource physical location to href relative location.
// phyDir : resource files physical location relative to html file.
// hrefDir: resource location as defined withing the href tag.
func configureResourceLocation(phyDir, hrefDir string) {
	fs := http.FileServer(http.Dir(phyDir))
	fs = http.StripPrefix("/"+hrefDir, fs)
	http.Handle("/"+hrefDir+"/", fs)
}

func completeAuth(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	_, err := auth.Token(spotStateValue, req)
	if err != nil {
		http.Error(res, "Couldn't get token", http.StatusForbidden)
		log.Errorf(ctx, "Error %v", err)
	}
	log.Infof(ctx, "Created")
}
