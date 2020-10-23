package main

import (
	"fmt"
	"time"
	"html/template"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	. "github.com/logrusorgru/aurora"
)

func Data(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Data..."))
	claims := &jwt.MapClaims{
		"exp":    time.Now().Add(time.Hour).Unix(),
		"params": map[string]string{},
		"resource": map[string]int{
			"dashboard": 1,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := "b8077945abd4d7bf06f4dcf3544dbb0a91fda42502a412ee627e2edb18a58d44"
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		panic(err.Error())
	}
	iframeUrl := "https://metabase.jeifai.com/embed/dashboard/" + tokenString
	infos := struct {
		Metabase string
	}{
		iframeUrl,
	}
	templates := template.Must(
		template.ParseFiles(
			"templates/OUT_navbar.html",
			"templates/OUT_head.html",
			"templates/OUT_data.html",
			"templates/OUT_footer.html",))
	templates.ExecuteTemplate(w, "layout", infos)
}