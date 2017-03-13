package csci152

import (
	"net/http"

	"github.com/zmb3/spotify"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func initSpotify(res http.ResponseWriter, req *http.Request) {
	// if baseURL == "" {
	ctx := appengine.NewContext(req)
	sr, err := spotify.DefaultClient.Search("Rumors", spotify.SearchTypeAlbum)
	// 	auth = spotify.NewAuthenticator(retrieveURI, spotify.ScopeUserReadPrivate)
	// 	// auth.SetAuthInfo(clientID, spotKey)
	// 	// baseURL = auth.AuthURL(spotStateValue)
	log.Infof(ctx, "Here ==> %v", "XXX")
	log.Infof(ctx, "Error %v", err)
	log.Infof(ctx, "Search: ID=%v, URI=%v", sr.Albums.Albums[0].ID, sr.Albums.Albums[0].URI)

	albumID := sr.Albums.Albums[0].ID

	trks, err := spotify.GetAlbumTracks(albumID)

	if err == nil {
		for _, val := range trks.Tracks {
			log.Infof(ctx, "Tracks=%v", val.Name)
		}
	} else {
		log.Infof(ctx, "Error %v", err)
	}
}
