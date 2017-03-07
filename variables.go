package csci152

import (
	"html/template"

	"github.com/zmb3/spotify"
)

var tpl *template.Template              // html web page processing / parsing object
var userInformation userInformationType // logged in user's information and preferences.
// var movieAPI *tmdbgae.TMDb              	// movie / tv database access object instance.
var baseURL string // base url to get images.
var auth spotify.Authenticator
var webInformation = webInformationType{
	User: &userInformation,
}