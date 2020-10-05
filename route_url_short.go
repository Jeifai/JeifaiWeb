package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/logrusorgru/aurora"
)

func ManageUrlShort(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting ManageUrlShort..."))
	sess := GetSession(r)
	urlshort, _ := mux.Vars(r)["urlshort"]
	resultid, url := SelectUrlByShortUrl(urlshort)
	InsertUserResultVisit(sess.UserId, resultid)
	http.Redirect(w, r, url, 301)
}
