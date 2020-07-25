package main

import (
	"fmt"
	"html/template"
    "net/http"
    "encoding/json"

	. "github.com/logrusorgru/aurora"
)

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Home..."))
	sess, err := GetSession(r)
	if err != nil {
		fmt.Println(Yellow("User not logged in..."))
		templates := template.Must(template.ParseFiles(
			"templates/OUT_layout.html",
			"templates/OUT_home.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		fmt.Println(Blue("User logged in..."))
		user := User{
			Id: sess.UserId,
		}
		user.UserById()

		home, err := user.GetHomeData()
		if err != nil {
			fmt.Println(Gray(8-1, "User has no data..."))
		}
		type TempStruct struct {
			User User
			Home HomeData
		}
		infos := TempStruct{user, home}
		templates := template.Must(
			template.ParseFiles(
				"templates/IN_layout.html",
				"templates/IN_topbar.html",
				"templates/IN_sidebar.html",
				"templates/IN_home.html"))
		templates.ExecuteTemplate(w, "layout", infos)
	}
}

func How(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting How..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_how.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func Pricing(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting How..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_pricing.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func Features(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Features..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_features.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func Faq(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Faq..."))
	templates := template.Must(template.ParseFiles(
		"templates/OUT_layout.html",
		"templates/OUT_faq.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}


func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Test..."))
	templates := template.Must(template.ParseFiles(
		"templates/test.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func TestMatch(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting TestMatches..."))

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	matches, err := user.MatchesByUser()

	type PublicMatch struct {
		CreatedDate string
		Target      string
		Title       string
		Url         string
	}

	var public_matches []PublicMatch
	for _, match := range matches {
		var public_match PublicMatch
		public_match.CreatedDate = match.CreatedDate
		public_match.Target = match.Target
		public_match.Title = match.Title
		public_match.Url = match.Url
		public_matches = append(public_matches, public_match)
	}

	type TempStruct struct {
		Data []PublicMatch
	}

	infos := TempStruct{public_matches}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func TestHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Home..."))
	sess, err := GetSession(r)
	if err != nil {
		fmt.Println(Yellow("User not logged in..."))
		templates := template.Must(template.ParseFiles(
			"templates/OUT_layout.html",
			"templates/OUT_home.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		fmt.Println(Blue("User logged in..."))
		user := User{
			Id: sess.UserId,
		}
		user.UserById()

		home, err := user.GetHomeData()
		if err != nil {
			fmt.Println(Gray(8-1, "User has no data..."))
        }
        home.UserName = user.UserName
		type TempStruct struct {
			Home HomeData
        }
		infos := TempStruct{home}
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(infos)
	}
}

func TestTargets(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting TestTargets..."))

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	user_targets_info, err := user.InfoUsersTargetsByUser()
	if err != nil {
		panic(err.Error())
	}

	var name_targets []string
	not_selected_targets, err := user.NotSelectedUsersTargetsByUser()
	if err != nil {
		panic(err.Error())
	}
	for _, v := range not_selected_targets {
		name_targets = append(name_targets, v.Name)
	}

	// vue-taggable-select does not work if name_targets is empty
	if len(name_targets) == 0 {
		name_targets = []string{" "}
	}

	type TempStruct struct {
		Targets     []TargetInfo
		NameTargets []string
	}

    infos := TempStruct{user_targets_info, name_targets}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func TestKeywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting TestKeywords..."))

	sess, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	user := User{
		Id: sess.UserId,
	}
	user.UserById()

	struct_targets, err := user.UsersTargetsByUser()
	if err != nil {
		panic(err.Error())
	}

	var arr_targets []string
	for _, v := range struct_targets {
		arr_targets = append(arr_targets, v.Name)
	}

	utks, err := user.GetUserTargetKeyword()
	if err != nil {
		panic(err.Error())
	}

	type PublicTargetKeyword struct {
		CreatedDate string
		KeywordText string
		TargetName  string
	}

	var public_utks []PublicTargetKeyword
	for _, utk := range utks {
		var public_utk PublicTargetKeyword
		public_utk.CreatedDate = utk.CreatedDate
		public_utk.KeywordText = utk.KeywordText
		public_utk.TargetName = utk.TargetName
		public_utks = append(public_utks, public_utk)
	}

	type TempStruct struct {
		Targets []string
		Utks    []PublicTargetKeyword
	}

	infos := TempStruct{arr_targets, public_utks}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}