package csci152

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/nu7hatch/gouuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

/**************************************************  public functions  **************************************************/

// UsernameExists checks the user name dictionary if an
func UsernameExists(req *http.Request) (bool, error) {
	var names dictionaryUserName

	ctx := appengine.NewContext(req)    // generate an appengine context.
	bs, err := ioutil.ReadAll(req.Body) // read username as it is typed in.
	if err != nil {
		return false, err // exit if error.
	}
	names.Name = string(bs) // convert.
	return readDataStore(ctx, nameDict, names.Name, &names) != nil, nil
}

// SearchUser searches if a username has been registered.
// return true and the user's uuid if found otherwise returns false and the uuid is empty.
func SearchUser(ctx context.Context, user string) (ukey string) {
	var ui dictionaryUserName

	user = strings.ToLower(user)
	dsQuery := datastore.NewQuery(nameDict).Run(ctx)
	for {
		_, err := dsQuery.Next(&ui)
		if err == datastore.Done {
			break
		}
		if found := strings.ToLower(ui.Name) == user; found {
			ukey = ui.UUID
			break
		}
	}
	return
}

// WriteNewUserInformation writes newly register user information and preferences to datastore / memcache.
// req contains all received information.
func WriteNewUserInformation(res http.ResponseWriter, req *http.Request) (registered bool) {
	var names dictionaryUserName
	var err error

	ctx := appengine.NewContext(req)

	pass := req.FormValue("newpassword")
	conf := req.FormValue("confirm")
	spot := "X1"
	un := req.FormValue("newusername")

	names.Name = un

	if pass == conf && spot != "" && un != "" {
		uid, _ := generateUUID()
		names = dictionaryUserName{
			Name: un,
			UUID: uid,
		}
		userInformation = userInformationType{
			UserID:         uid,
			SpotifyAccount: spot,
			Password:       EncryptPassword(pass),
			Username:       un,
			LoggedIn:       true,
		}
		if err = writeDataStore(ctx, nameDict, un, &names); err == nil {
			err = WriteUserInformation(ctx, req)
		}
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		registered = true
	}
	return
}

// WriteUserInformation write user data to both memcache and datastore.
func WriteUserInformation(ctx context.Context, req *http.Request) error {
	var err error

	err = writeDataStore(ctx, dataKind, userInformation.UserID, &userInformation)
	if err == nil {
		err = writeMemcache(ctx, req)
	}
	return err
}

// ReadUserInformation reads user information first from memcache,
// and if not present read from datastore and write that data back into memcache.
func ReadUserInformation(ctx context.Context, req *http.Request, userID string) error {
	var err error
	if !readMemcache(ctx, req, userID) {
		err := readDataStore(ctx, dataKind, userID, &userInformation)
		if err == nil {
			err = writeMemcache(ctx, req)
		}
	}
	return err
}

// EncryptPassword encrypts user password using prefix and suffix salt value and using sha256 hashing.
func EncryptPassword(pass string) string {
	h := sha256.New()
	io.WriteString(h, passwordPrefix)
	io.WriteString(h, pass)
	io.WriteString(h, passwordSuffix)
	return fmt.Sprintf("%x", h.Sum(nil))
}

/**************************************************  private functions  **************************************************/

// read designated data from the desired datastore.
// returns if there was an error.
// data is returned within the data interface parameter.
func readDataStore(ctx context.Context, kind, key string, data interface{}) error {
	dsKey := datastore.NewKey(ctx, kind, key, 0, nil)
	err := datastore.Get(ctx, dsKey, data)
	return err
}

// write data to the datastore.
//returns if there was an error, reports an error 500.
func writeDataStore(ctx context.Context, kind, key string, data interface{}) error {
	dsKey := datastore.NewKey(ctx, kind, key, 0, nil)
	_, err := datastore.Put(ctx, dsKey, data)
	return err
}

// read information based on logged in user's uuid.
// returns true if data read successfully.
func readMemcache(ctx context.Context, req *http.Request, userID string) bool {
	item, err := memcache.Get(ctx, userID)
	if err == nil {
		err = json.Unmarshal(item.Value, &userInformation)
	}
	return err == nil
}

// write user data to memcache.
// data is to be defined within the userInformation variable.
func writeMemcache(ctx context.Context, req *http.Request) error {

	bs, err := json.Marshal(userInformation)
	if err != nil {
		return err
	}
	memData := memcache.Item{
		Key:   userInformation.UserID,
		Value: bs,
	}
	err = memcache.Set(ctx, &memData)
	if err != nil {
		return err
	}
	return nil
}

// get an UUID for user.
func generateUUID() (string, error) {
	uuid, err := uuid.NewV4()
	return uuid.String(), err
}

// ToInt converts returned form value data to integer
// req: http request containing data to be converted.
// key: field key / name of data control.
// returns converted value, if error, returns 0
func toInt(req *http.Request, key string) (val int) {
	tv, err := strconv.Atoi(req.FormValue(key))
	if err == nil {
		val = int(tv)
	}
	return
}

// set defaults.
func setUserDefault() {
	rdr := radarType{nil, nil}
	userInformation = userInformationType{"", "", "", "", false, false, "", ""}
	webInformation = webInformationType{&userInformation, rdr}
}

// // tv information and search results.
// func tvPost(ctx context.Context, res http.ResponseWriter, req *http.Request) {

// 	info := toInt(req, "cmdID")     // get possible movie id to show detail.
// 	watchID := toInt(req, "cmdAdd") // add to favorites / watch list.

// 	webInformation.MovieTvGame.ID = 0 // no detail, search.
// 	if info != 0 {
// 		tvInfo(ctx, info)
// 	}
// 	if watchID != 0 && !duplicate(int32(watchID), 1) {
// 		w := watch{int32(watchID), 1}
// 		userInformation.Watched = append(userInformation.Watched, w)
// 		updateCookie(res, req)
// 		WriteUserInformation(ctx, req) // write added item to datastore / memcache
// 	}
// 	executeSearch(res, req)
// }

// // movie information and search results.
// func moviePost(ctx context.Context, res http.ResponseWriter, req *http.Request) {

// 	info := toInt(req, "cmdID")     // get possible movie id to show detail.
// 	watchID := toInt(req, "cmdAdd") // add to favorites / watch list.

// 	webInformation.MovieTvGame.ID = 0 // no detail, search.
// 	if info != 0 {
// 		movieInfo(ctx, info)
// 	}
// 	if watchID != 0 && !duplicate(int32(watchID), 0) {
// 		w := watch{int32(watchID), 0}
// 		userInformation.Watched = append(userInformation.Watched, w)
// 		updateCookie(res, req)
// 		WriteUserInformation(ctx, req) // write added item to datastore / memcache
// 	}
// 	executeSearch(res, req)
// }

// func duplicate(id int32, mtgType int) bool {
// 	var dup = false

// 	for _, wid := range userInformation.Watched {
// 		dup = (wid.ID == id) && (wid.MTGType == mtgType)
// 		if dup {
// 			break
// 		}
// 	}
// 	return dup // if duplicate return true, otherwise return false.
// }

// // movie / tv information results.
// func movieTvPost(ctx context.Context, res http.ResponseWriter, req *http.Request) {

// 	watchID := toInt(req, "cmdMAdd") // add to favorites / watch list.
// 	if watchID != 0 && !duplicate(int32(watchID), 0) {
// 		w := watch{int32(watchID), 0}
// 		userInformation.Watched = append(userInformation.Watched, w)
// 		updateCookie(res, req)
// 		WriteUserInformation(ctx, req) // write added item to datastore / memcache
// 	}

// 	watchID = toInt(req, "cmdTAdd") // add to favorites / watch list.
// 	if watchID != 0 && !duplicate(int32(watchID), 1) {
// 		w := watch{int32(watchID), 1}
// 		userInformation.Watched = append(userInformation.Watched, w)
// 		updateCookie(res, req)
// 		WriteUserInformation(ctx, req) // write added item to datastore / memcache
// 	}

// 	watchID = toInt(req, "cmdGAdd") // add to favorites / watch list.
// 	if watchID != 0 && !duplicate(int32(watchID), 2) {
// 		w := watch{int32(watchID), 2}
// 		userInformation.Watched = append(userInformation.Watched, w)
// 		updateCookie(res, req)
// 		WriteUserInformation(ctx, req) // write added item to datastore / memcache
// 	}

// 	infoID := toInt(req, "cmdMID")
// 	if infoID != 0 {
// 		movieInfo(ctx, infoID)
// 	}
// 	infoID = toInt(req, "cmdTID")
// 	if infoID != 0 {
// 		tvInfo(ctx, infoID)
// 	}
// 	i := req.FormValue("cmdGID")
// 	infoID = toInt(req, "cmdGID")
// 	if infoID != 0 {
// 		gameInfo(ctx, infoID, i)
// 	}
// 	executeSearch(res, req)
// }

// // gamePost retrieves and formats individual game information
// func gamePost(ctx context.Context, res http.ResponseWriter, req *http.Request) {
// 	i := req.FormValue("cmdID")
// 	info := toInt(req, "cmdID")
// 	watchID := toInt(req, "cmdAdd")

// 	webInformation.MovieTvGame.ID = 0
// 	if info != 0 {
// 		gameInfo(ctx, info, i)
// 	}
// 	if watchID != 0 && !duplicate(int32(watchID), 2) {
// 		w := watch{int32(watchID), 2}
// 		userInformation.Watched = append(userInformation.Watched, w)
// 		updateCookie(res, req)
// 		WriteUserInformation(ctx, req)
// 	}
// 	executeSearch(res, req)
// }

// func gameInfo(ctx context.Context, info int, i string) {
// 	gm, err := igdbgo.GetGames(ctx, "", 1, 0, 0, i)
// 	if err != nil {
// 		return
// 	}
// 	game := gm[0]

// 	webInformation.MovieTvGame.ID = info
// 	webInformation.MovieTvGame.mtgType = 2
// 	webInformation.MovieTvGame.Description = game.Summary
// 	if game.FirstRelease != 0 {
// 		y, m, d := game.GetDate()
// 		date := strconv.Itoa(m) + "/" + strconv.Itoa(d) + "/" + strconv.Itoa(y)
// 		webInformation.MovieTvGame.ReleaseDate = date
// 	} else {
// 		webInformation.MovieTvGame.ReleaseDate = "Future"
// 	}
// 	webInformation.MovieTvGame.Image = game.GetImageURL()
// 	webInformation.MovieTvGame.Genres = game.GetGenres()
// 	webInformation.MovieTvGame.Youtube, err = game.GetVideoURL()
// 	if err != nil {
// 		webInformation.MovieTvGame.Youtube = ""
// 	}
// }

// // execute movie / tv / game search.
// func executeSearch(res http.ResponseWriter, req *http.Request) {
// 	searchCmd := req.FormValue("cmdSearch") // get possible search type.
// 	search := req.FormValue("search")       // get possible title to seach for.

// 	if searchCmd == "movies_tv" {
// 		http.Redirect(res, req, fmt.Sprintf("/results?srch=%s", search), http.StatusFound)
// 	}
// }

// func movieInfo(ctx context.Context, movieID int) {
// 	var g []string

// 	burl, _ := movieAPI.GetConfiguration(ctx)
// 	mvi, _ := movieAPI.GetMovieInfo(ctx, movieID, nil)
// 	// trail, _ := movieAPI.GetMovieVideos(ctx, info, nil)

// 	webInformation.MovieTvGame.ID = movieID
// 	webInformation.MovieTvGame.mtgType = 0
// 	webInformation.MovieTvGame.Description = mvi.Overview
// 	if mvi.ReleaseDate != "" {
// 		webInformation.MovieTvGame.ReleaseDate = formatDate(mvi.ReleaseDate)
// 	} else {
// 		webInformation.MovieTvGame.ReleaseDate = "Future"
// 	}
// 	webInformation.MovieTvGame.Image = fmt.Sprintf("%s%s%s", burl.Images.BaseURL, burl.Images.PosterSizes[1], mvi.PosterPath)
// 	for _, gn := range mvi.Genres {
// 		g = append(g, gn.Name)
// 	}
// 	webInformation.MovieTvGame.Genres = g

// 	mv, _ := movieAPI.GetMovieVideos(ctx, movieID, nil)
// 	log.Infof(ctx, "Trailer: %v", mv)
// 	// log.Infof(ctx, "Trailer: %s", mv.Results[0].Key)

// 	if len(mv.Results) > 0 {
// 		webInformation.MovieTvGame.Youtube = fmt.Sprintf(youTubeBase, mv.Results[0].Key)
// 	} else {
// 		webInformation.MovieTvGame.Youtube = ""
// 	}
// }

// func tvInfo(ctx context.Context, tvID int) {
// 	var g []string
// 	burl, _ := movieAPI.GetConfiguration(ctx)
// 	tvi, _ := movieAPI.GetTvInfo(ctx, tvID, nil)
// 	webInformation.MovieTvGame.ID = tvID
// 	webInformation.MovieTvGame.mtgType = 1
// 	webInformation.MovieTvGame.Image = fmt.Sprintf("%s%s%s", burl.Images.BaseURL, burl.Images.PosterSizes[1], tvi.PosterPath)
// 	webInformation.MovieTvGame.Description = tvi.Overview
// 	webInformation.MovieTvGame.TVSeasons = tvi.NumberOfSeasons
// 	webInformation.MovieTvGame.TVEpisodes = tvi.NumberOfEpisodes
// 	webInformation.MovieTvGame.ReleaseDate = formatDate(tvi.FirstAirDate)
// 	log.Infof(ctx, "Air Date: %s", tvi.FirstAirDate)
// 	for _, gn := range tvi.Genres {
// 		g = append(g, gn.Name)
// 	}
// 	webInformation.MovieTvGame.Genres = g
// }

// // roundFloat is an intermediate function that rounds a float64 into an int
// func roundFloat(num float64) int {
// 	return int(num + math.Copysign(0.5, num))
// }

// // setPrecision is an intermediate function that formats a float64 to a specific number of significat digits
// func setPrecision(num float64, prec int) float64 {
// 	output := math.Pow(10, float64(prec))
// 	return float64(roundFloat(num*output)) / output
// }

// // round is a function that takes a float64 and rounds it down to a number out of 10 with only 1 significant digit
// func round(num float64) float32 {
// 	return float32(setPrecision((num / 10), 1))
// }

// // formatDate takes a date string that has been previously formatted by the Movie and TV API and returns
// // a date formatted in the American format
// func formatDate(date string) string {
// 	s := strings.Split(date, "-")
// 	m, _ := strconv.Atoi(s[1])
// 	d, _ := strconv.Atoi(s[2])
// 	return strconv.Itoa(m) + "/" + strconv.Itoa(d) + "/" + s[0]
// }
