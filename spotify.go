package csci152

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/ljmeyers80529/spot-go-gae"
)

func completeAuth(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	log.Infof(ctx, "Callback executed")

	_, err := auth.Token(spotStateValue, req)
	if err != nil {
		// 	http.Error(res, "Couldn't get token", http.StatusForbidden)
		// 	log.Errorf(ctx, "Error %v", err)
		// } else {
		// 	spotClient = auth.NewClient(spotToken)
	}
}

func initSpotify(res http.ResponseWriter, req *http.Request) {
	// ctx := appengine.NewContext(req)
	// log.Infof(ctx, "Check ==> %v", spotToken)
	if spotToken == nil {
		auth := spotify.NewAuthenticator(retrieveURI, spotify.ScopeUserReadPrivate)
		auth.SetAuthInfo(clientID, spotKey)
		http.Redirect(res, req, auth.AuthURL(spotStateValue), http.StatusFound)
	}
	// sr, err := spotify.DefaultClient.Search("Rumors", spotify.SearchTypeAlbum)
	// baseURL = auth.AuthURL(spotStateValue)
	// log.Infof(ctx, "Here ==> %v", "XXX")
	// log.Infof(ctx, "Auth => %v\tCfg => %v", auth)
	// log.Infof(ctx, "Error %v", err)
	// log.Infof(ctx, "Search: ID=%v, URI=%v", sr.Albums.Albums[0].ID, sr.Albums.Albums[0].URI)

	// albumID := sr.Albums.Albums[0].ID

	// trks, err := spotify.GetAlbumTracks(albumID)

	// if err == nil {
	// 	for _, val := range trks.Tracks {
	// 		log.Infof(ctx, "Tracks=%v", val.Name)
	// 	}
	// } else {
	// 	log.Infof(ctx, "Error %v", err)
	// }
}
