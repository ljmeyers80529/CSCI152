package csci152

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

/**************************************************  constants, types and variables  **************************************************/

const cookieSessionName string = "cookieTastTest"

// readCookie reads current state
// create a new cookie if it does not exists or expired.
func readCookie(res http.ResponseWriter, req *http.Request) {
	cookie := readCreateCookie(req)
	http.SetCookie(res, cookie)                               // set cookie into browser.
	userInformation = cookieInformationDecoding(cookie.Value) // decode and set user state into page variable.
	if !userInformation.LoggedIn {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
	}
}

// read an existing cookie or create a new one.
// returns the cookie.
func readCreateCookie(req *http.Request) (cookie *http.Cookie) {
	cookie, err := req.Cookie(cookieSessionName) // get if a cookie already exists (had not expired)
	if err == http.ErrNoCookie {
		cookie = newCookie() // need a new cookie.
	}
	return
}

// create a new cookie, set value fields to default values, JSON / base 64 processed.
func newCookie() (cookie *http.Cookie) {
	cookie = &http.Cookie{
		Name:     cookieSessionName,
		Value:    cookieInformationEncoding(),
		HttpOnly: true,
		//Secure: false,
	}
	return
}

// UpdateCookie updates cookie information.
// read existing cookie or re-sreate the cookie if expired.
func updateCookie(res http.ResponseWriter, req *http.Request) {
	cookie := readCreateCookie(req)
	cookie.Value = cookieInformationEncoding()
	http.SetCookie(res, cookie) // set cookie into browser.
}

// encode user state information marshalled using JSON and then converted into base 64.
// returns JSON / base 64 encoded string.
func cookieInformationEncoding() (encoded string) {
	j, err := json.Marshal(userInformation) // encode using JSON and base 64.
	if err == nil {
		encoded = base64.URLEncoding.EncodeToString(j)
	}
	return
}

// retrieve information from the cookie and decode from base 64 then unmarshal JSON.
// returns user information data.
func cookieInformationDecoding(userInfo string) (ci userInformationType) {
	decode, _ := base64.URLEncoding.DecodeString(userInfo)
	json.Unmarshal(decode, &ci)
	return
}
