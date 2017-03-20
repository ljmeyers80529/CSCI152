package csci152

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	configureResourceLocation("images", "img")
	configureResourceLocation("css", "css")
	// configureResourceLocation("images", "js/images")
	configureResourceLocation("js", "js")
	setUserDefault()
	http.Handle("/favicon.ico", http.NotFoundHandler()) // ignore favicon request (error 404)
	http.HandleFunc("/", pageMain)                      // index page to check if already logged in or need to login.
	http.HandleFunc("/login", pageLogin)
	http.HandleFunc("/logout", pageLogout)
	http.HandleFunc("/register", pageRegister)
	http.HandleFunc("/home", pageHome)
	// http.HandleFunc("/username/check", pageRegisterUsernameCheck) // verify username is unique.
	// http.HandleFunc("/about", pageAbout)                          // about web page.
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
	log.Infof(ctx, "Callback executed")
	// _, err := auth.Token(spotStateValue, req)
	// if err != nil {
	// 	http.Error(res, "Couldn't get token", http.StatusForbidden)
	// 	log.Errorf(ctx, "Error %v", err)
	// }
}
