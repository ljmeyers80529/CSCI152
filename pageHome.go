package csci152

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func pageHome(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	readCookie(res, req)
	if req.Method == "POST" {
	}
	initSpotify(res, req)
	log.Infof(ctx, "Client code => %v", spotClient)
	tpl.ExecuteTemplate(res, "homepage.html", webInformation)
}
