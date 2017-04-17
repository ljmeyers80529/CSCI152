package csci152

import (
	"net/http"
)

func pageHome(res http.ResponseWriter, req *http.Request) {
	// ctx := appengine.NewContext(req)
	readCookie(res, req)
	if req.Method == "POST" {
	}
	initSpotify(res, req)
	// log.Infof(ctx, "Client code in home => %v", spotClient)
	// loadPlayLists(res, req)
	tpl.ExecuteTemplate(res, "homepage.html", webInformation)
}
