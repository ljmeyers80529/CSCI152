package csci152

type dictionaryUserName struct { // name type to be read from dictionary of usernames from datastore
	Name, UUID string
}

// // countdown timer implementation definition.
// type upComming struct {
// 	Title                                string
// 	ID, Year, Month, Day, Hours, Minutes int
// }

// // list of top rated movies.
// type topRatedPop struct {
// 	Title   string
// 	ID      int
// 	Rating  float32
// 	Release string
// }

// // definition for maovie / tv / game detail.
// type movieTvGameInformation struct {
// 	ID                              int
// 	Image, Description, ReleaseDate string
// 	TVSeasons, TVEpisodes           int
// 	Genres                          []string
// 	UID                             int
// 	Search                          string
// 	Youtube                         string
// 	mtgType							int
// }

// // user's favorites / watch list.
// type watch struct {
// 	ID      int32
// 	MTGType int // 0 = movie, 1 = tv, 2 = game
// }

// // all upcomming movies list.
// type cdUpcomming []upComming

// // all top or popular rated movies.
// type topPopRated []topRatedPop

// type to contain user information and preferences.
type userInformationType struct {
	UserID   string
	Name     string
	Password string
	Username string
	LoggedIn bool
	// Watched  []watch
}

// // watched definition.
// type watchedType struct {
// 	ID                               int
// 	Title                            string
// 	Rating                           float32
// 	Movie, TV, Game, Future          bool
// 	Release                          string
// 	Year, Month, Day, Hours, Minutes int
// }

// type definition to rendering information to the website.
type webInformationType struct {
	User        *userInformationType
	// Counters    cdUpcomming
	// Top, Pop, Game    topPopRated
	// MovieTvGame movieTvGameInformation
	// Watched     []watchedType
}
