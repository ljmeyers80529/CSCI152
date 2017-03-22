package csci152

import (
	"html/template"

	"golang.org/x/oauth2"

	"github.com/ljmeyers80529/spot-go-gae"
)

var tpl *template.Template              // html web page processing / parsing object
var userInformation userInformationType // logged in user's information and preferences.
// var baseURL string = "" 				// base url to get images.
var auth spotify.Authenticator
var spotToken *oauth2.Token
var spotClient spotify.Client

var webInformation = webInformationType{
	User: &userInformation,
}
