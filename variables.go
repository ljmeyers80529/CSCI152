package csci152

import (
	"html/template"
)

var tpl *template.Template              // html web page processing / parsing object
var userInformation userInformationType // logged in user's information and preferences.
// var movieAPI *tmdbgae.TMDb              // movie / tv database access object instance.
var baseURL string						// base url to get images.
var webInformation = webInformationType {
	User: &userInformation,
}
