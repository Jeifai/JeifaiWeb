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
	user := UserById(sess.UserId)

	urlshort, _ := mux.Vars(r)["urlshort"]
	resultid, url := SelectUrlByShortUrl(urlshort)
	
	InsertUserResultVisit(user.Id, resultid) 

	http.Redirect(w, r, url, 301)
}
