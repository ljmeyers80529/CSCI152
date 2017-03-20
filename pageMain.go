package csci152

import (
	// "fmt"
	"net/http"
	// "strings"
	// "strconv"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// main (top) web page.
func pageMain(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	readCookie(res, req) // maintain user login / out state.apikey
	log.Infof(ctx, "Cookie = %v", webInformation.User)
	if webInformation.User.LoggedIn {
		log.Infof(ctx,"%s", "Loggedin")
		// http.Redirect(res, req, "/home", http.StatusSeeOther)
	} else {
		log.Infof(ctx,"%s", "Needs to Loggedin")
		// userLogin(res, req)
	}

	// initSpotify(res, req);
	// popWatch(ctx)

	// if req.Method == "POST" {
	// 	infoID := toInt(req, "cmdID")
	// 	removeID := toInt(req, "cmdRM")
	// 	if infoID > 0 {
	// 		switch itemType(infoID) {
	// 		case 0:
	// 			moviePost(ctx, res, req)
	// 			if webInformation.MovieTvGame.ID != 0 { // no detail, search.
	// 				http.Redirect(res, req, fmt.Sprintf("%s#moviemodal", req.URL.Path), http.StatusFound)
	// 			}
	// 		case 1:
	// 			tvPost(ctx, res, req)
	// 			if webInformation.MovieTvGame.ID != 0 { // no detail, search.
	// 				http.Redirect(res, req, fmt.Sprintf("%s#tvmodal", req.URL.Path), http.StatusFound)
	// 			}
	// 		case 2:
	// 			gamePost(ctx, res, req)
	// 			if webInformation.MovieTvGame.ID != 0 {
	// 				http.Redirect(res, req, fmt.Sprintf("%s#gamemodal", req.URL.Path), http.StatusFound)
	// 			}
	// 		}
	// 	}
	// 	if removeID > 0 {
	// 		if removeID == 4096 {
	// 			removeID = 0
	// 		}
	// 		userInformation.Watched = append(userInformation.Watched[:removeID], userInformation.Watched[removeID+1:]...)
	// 		updateCookie(res, req)
	// 		WriteUserInformation(ctx, req) // write added item to datastore / memcache
	// 		popWatch(ctx)
	// 	}
	// 	executeSearch(res, req)
	// }
	tpl.ExecuteTemplate(res, "index.html", webInformation)
}

// // get watch items.
// func popWatch(ctx context.Context) {
// 	var wi watchedType
// 	var wis []watchedType

// 	for _, wat := range userInformation.Watched {
// 		wi.ID = int(wat.ID)
// 		wi.Movie = false
// 		wi.TV = false
// 		wi.Game = false
// 		switch wat.MTGType {
// 		case 0:
// 			mvi, _ := movieAPI.GetMovieInfo(ctx, wi.ID, nil)
// 			wi.Movie = true
// 			wi.Title = mvi.Title
// 			wi.Rating = float32(setPrecision(float64(mvi.VoteAverage), 1))
// 			wi.Release = formatDate(mvi.ReleaseDate)
// 			dr, ok := movieRelease(ctx, wi.ID)
// 			wi.Future = ok
// 			if wi.Future {
// 				s := strings.Split(dr, "-")
// 				wi.Year, _ = strconv.Atoi(s[0])
// 				wi.Month, _ = strconv.Atoi(s[1])
// 				wi.Day, _ = strconv.Atoi(s[2])
// 				wi.Hours = 0
// 				wi.Minutes = 0
// 			}
// 		case 1:
// 			tvi, _ := movieAPI.GetTvInfo(ctx, wi.ID, nil)
// 			wi.TV = true
// 			wi.Future = false
// 			wi.Title = tvi.Name
// 			wi.Rating = float32(setPrecision(float64(tvi.VoteAverage), 1))
// 			wi.Release = formatDate(tvi.FirstAirDate)
// 		case 2:
// 			gms, _ := igdbgo.GetGames(ctx, "", 1, 0, 0, strconv.Itoa(wi.ID))
// 			gmi := gms[0]
// 			wi.Game = true
// 			wi.Title = gmi.Name
// 			wi.Rating = round(gmi.Rating)
// 			y, m, d := gmi.GetDate()
// 			date := strconv.Itoa(m) + "/" + strconv.Itoa(d) + "/" + strconv.Itoa(y)
// 			wi.Release = date
// 			wi.Future = gmi.CheckFuture()
// 			if wi.Future {
// 				wi.Year = y
// 				wi.Month = m
// 				wi.Day = d
// 				wi.Hours = 0
// 				wi.Minutes = 0
// 			}
// 		}
// 		wis = append(wis, wi)
// 	}
// 	webInformation.Watched = wis
// }

// func itemType(ID int) int {
// 	var t = -1

// 	for _, val := range webInformation.Watched {
// 		if val.ID == ID {
// 			if val.Movie {
// 				t = 0
// 			} else if val.TV {
// 				t = 1
// 			} else {
// 				t = 2
// 			}
// 		}
// 	}
// 	return t
// }

// // remove an item from the watch list.
// func removeItem(ID int) int {
// 	var location = 0

// 	w := webInformation.Watched
// 	for _, val := range w {
// 		if val.ID == ID {
// 			break
// 		}
// 		location++
// 	}
// 	// w = append(w[:location], w[location+1:]...)
// 	// webInformation.Watched = w
// 	return location
// }
