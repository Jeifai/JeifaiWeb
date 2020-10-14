package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	. "github.com/logrusorgru/aurora"
)

func ServeMetabaseJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting ServeMetabaseJobs..."))
	claims := &jwt.MapClaims{
		"exp":    time.Now().Add(time.Hour).Unix(),
		"params": map[string]string{},
		"resource": map[string]int{
			"dashboard": 2,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := "b8077945abd4d7bf06f4dcf3544dbb0a91fda42502a412ee627e2edb18a58d44"
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		panic(err.Error())
	}
	iframeUrl := "https://34.120.21.103/embed/dashboard/" + tokenString
	infos := struct {
		Metabase string
	}{
		iframeUrl,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func ServeMetabaseCompany(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting ServeMetabaseCompany..."))
	claims := &jwt.MapClaims{
		"exp":    time.Now().Add(time.Hour).Unix(),
		"params": map[string]string{},
		"resource": map[string]int{
			"dashboard": 3,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := "b8077945abd4d7bf06f4dcf3544dbb0a91fda42502a412ee627e2edb18a58d44"
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		panic(err.Error())
	}
	iframeUrl := "https://34.120.21.103/embed/dashboard/" + tokenString
	infos := struct {
		Metabase string
	}{
		iframeUrl,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
