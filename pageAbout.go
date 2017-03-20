package csci152

import (
	"net/http"
)

func pageAbout(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "about.html", nil)
}
