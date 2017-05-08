package csci152

import (
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
	// ctx := appengine.NewContext(req)
	// log.Infof(ctx, "Spot login state: %v", webInformation.User.SpotLogged)
	if !webInformation.User.SpotLogged {
		auth = spotify.NewAuthenticator(retrieveURI, spotify.ScopePlaylistModifyPublic, spotify.ScopeUserTopRead)
		auth.SetAuthInfo(clientID, spotKey)
		http.Redirect(res, req, auth.AuthURL(spotStateValue), http.StatusFound)
	}
}

// check if client object has been defined; returns true if Ok..
func clientOK() bool {
	return !(spotClient == spotify.Client{})
}

// get spotify username.
func spotifyUser() string {
	var user = "ZZZZZ"

	if clientOK() {
		cu, err := spotClient.CurrentUser()
		if err == nil {
			user = cu.ID
		} else {
			user = ""
		}
	}
	return user
}

func getLoggedInUsersPlaylist(res http.ResponseWriter, req *http.Request) string {
	var pID string

	// ctx := appengine.NewContext(req)
	// log.Infof(ctx, "In playlist")
	if clientOK() {
		// log.Infof(ctx, "Client Not Empty")

		// user, _ := spotClient.CurrentUser()
		// log.Infof(ctx, "User ID: %v", user.ID)
		// //
		// sr, err := spotClient.CurrentUsersPlaylists()
		sr, _ := spotClient.CurrentUsersPlaylists()
		// srx, err := spotClient.CurrentUserRecentTracks(5)
		//
		// pID = sr.Playlists[0].ID.String()
		// log.Infof(ctx, "Playlist (full) ==> %v", sr.Playlists[0])
		// log.Infof(ctx, "Playlist (full) ==> %v", sr.Playlists[0].URI)
		pID = string(sr.Playlists[0].URI)
		// 	if err == nil {
		// 		for _, val := range sr.Playlists {
		// 			trk, err := spotClient.GetPlaylistTracks(user.ID, val.ID)
		// 			if err == nil {
		// 				for _, tVal := range trk.Tracks {
		// 					log.Infof(ctx, "\nPlaylist (track) ==> %v\n", tVal.Track.Name)
		// 				}
		// 			}
		// 			if err != nil {
		// 				log.Infof(ctx, "Error")
		// 			}
		// 		}
		// 		log.Infof(ctx, "Items ==> %v", srx.Items)
		// 		// for _, val := range srx.Items {
		// 		// 	trk, err := spotClient.GetPlaylistTracks(user.ID, val.ID)
		// 		// 	if err == nil {
		// 		// 		for _, tVal := range trk.Tracks {
		// 		// 			log.Infof(ctx, "\nPlaylist (track) ==> %v\n", tVal.Track.Name)
		// 		// 		}
		// 		// 	}
		// 		if err != nil {
		// 			log.Infof(ctx, "Error")
		// 			// }
		// 		}
		// 	} else {
		// 		log.Infof(ctx, "Error %v", err)
		// 	}
		// } else {
		// 	log.Infof(ctx, "Client Empty")
	}
	return pID
	// http.Redirect(res, req, "/home", http.StatusSeeOther)

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
