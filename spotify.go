package csci152

import (
	"flag"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/ljmeyers80529/spot-go-gae"
)

func completeAuthentication(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	log.Infof(ctx, "Callback executed")

	token, err := auth.Token(spotStateValue, req)
	if err != nil {
		http.Error(res, "Couldn't get token", http.StatusForbidden)
		log.Errorf(ctx, "Error %v", err)
	} else {
		spotClient = auth.NewClient(token)
		webInformation.User.SpotLogged = true
		updateCookie(res, req)
		http.Redirect(res, req, "/home", http.StatusSeeOther)
	}
}

func initSpotify(res http.ResponseWriter, req *http.Request) {
	if !webInformation.User.SpotLogged {
		auth = spotify.NewAuthenticator(retrieveURI, spotify.ScopeUserReadPrivate)
		auth.SetAuthInfo(clientID, spotKey)
		http.Redirect(res, req, auth.AuthURL(spotStateValue), http.StatusFound)
	}
}

// check if client object has been defined; returns true if Ok..
func clientOK() bool {
	return (spotClient == spotify.Client{})
}

func loadPlayLists(res http.ResponseWriter, req *http.Request) {
	flag.Parse()

	ctx := appengine.NewContext(req)
	log.Infof(ctx, "In playlist")
	if clientOK() {
		log.Infof(ctx, "Client Empty")
	} else {
		log.Infof(ctx, "Client not empty")

		user, _ := spotClient.CurrentUser()
		// user, _ := spotify.GetUsersPublicProfile(spotify.ID(*userID))
		log.Infof(ctx, "User ID: %v", user.ID)
		// log.Infof(ctx, "Display name: %v", user.DisplayName)
		// log.Infof(ctx, "Spotify URI: %v", string(user.URI))
		// log.Infof(ctx, "Endpoint: %v", user.Endpoint)
		// log.Infof(ctx, "Followers: %v", user.Followers.Count)
		sr, err := spotClient.CurrentUsersPlaylists()
		// sr, err := spotClient.GetPlaylistsForUser(user.ID)
		log.Infof(ctx, "Playlist (full) ==> %v", sr)
		if err == nil {
			// trks := sr.Playlists[0].Name
			// log.Infof(ctx, "\nPlaylist (track) ==> %v\n", sr.Playlists[0].ID)
			for _, val := range sr.Playlists {
				trk, err := spotClient.GetPlaylistTracks(user.ID, val.ID)
				if err == nil {
					for _, tVal := range trk.Tracks {
						log.Infof(ctx, "\nPlaylist (track) ==> %v\n", tVal.Track.Name)
					}
				}
				// trk, _ := spotClient.GetTrack(val.ID)
				// log.Infof(ctx, "\nPlaylist (track) ==> %v\n", val.ID)
				// trk, err := spotClient.GetTrack(val.ID)
				if err != nil {
					log.Infof(ctx, "Error")
				}
			}
			// 	// for _, val := range trks {
			// 	// 	log.Infof(ctx, "Tracks=%v", val.Name)
			// 	// }
		} else {
			log.Infof(ctx, "Error %v", err)
		}
	}
	http.Redirect(res, req, "/home", http.StatusSeeOther)

}

// 	// log.Infof(ctx, "Playlist ==> %v", sr)

// 	// sr, err := spotClient.CurrentUsersPlaylists()

// 	// sr, err := spotClient.Search("Rumors", spotify.SearchTypeAlbum)
// 	// log.Infof(ctx, "Error %v", err)
// 	// log.Infof(ctx, "Search: ID=%v, URI=%v", sr.Albums.Albums[0].ID, sr.Albums.Albums[0].URI)

// 	// albumID := sr.Albums.Albums[0].ID

// 	// trks, err := spotify.GetAlbumTracks(albumID)

// 	// if err == nil {
// 	// 	log.Infof(ctx, "Playlist ==> %v", sr)
// 	// 	// for _, val := range trks.Tracks {
// 	// 	// 	log.Infof(ctx, "Tracks=%v", val.Name)
// 	// } else {
// 	// 	log.Infof(ctx, "Error %v", err)
// 	// }

// }
