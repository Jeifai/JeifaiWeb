# JaifaiWeb

Here everything about the web application of Jaifai.

As any webapp JeifaiWeb has two sides:

* *Backend*: using [Golang](https://golang.org/)
* *Frontend*: using [Vuejs](https://vuejs.org/)


## Backend
There are three main blocks of code in the backend.

* *Main*, the starting point.
* *Route*, coordinator between the main and the database.
* *Db*, all the queries to manage and manipulate data.

A typical scenario:

* User is inside the webapp and opens the page dedicated to his Profile.
* *main.go* receives the request.
* *main.go* calls the function *Profile()* in the file *route_profile.go*.
* *route_profile.go* receives the request.
* *route_profile.go* call the function *UserById()* in the file *db_user.go*.
* *db_user.go* receives the request.
* *db_user.go* execute the query, get the data and return them to *route_profile.go*.
* *route_profile.go* receives the data and returns them to the user.

### Main

The starting point is the file called *main.go*, which contains all the routes of the webapp.
The purpose of the file is to act as router of any page, as well as for any API endpoint created.

For example:

```golang
func main() {
    r.HandleFunc("/", ShowHomePage).Methods("GET")
    r.HandleFunc("/faq", ShowFaqPage).Methods("GET")
    r.HandleFunc("/api_1", ManageApi1).Methods("POST")
}
```

### Route

Here all the files called *route...go*, for example *route_main.go*.
The purpose of these files is to act as middlelayer between the *main.go* and all the databases files.

For example:

```golang
func GetUserKeywords(w http.ResponseWriter, r *http.Request) {

    sess := GetSession(r) // Authenticate the user
    user := UserById(sess.UserId) // Get user's data

     // Extract data from the db
    keywords := user.SelectKeywordsByUser()
    infos := struct {
        Keywords []Keyword
    }{
        keywords,
    }

    // Return data to the client
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(infos)
}
```


## Scrape

This is Jeifai's core. Here happens the extraction of data from career's pages.

How the program flows:

* Start main function
* Connect to the database
* Get the scraper
* Start scraping session
* Scrape
* Save results to the database

### scrapers.go
All the scrapers are included in a single file called *scrapers.go*. The most important aspect is to produce scrapers as identical as possibile to each other. The two main tools used to scrape are:
* *[Colly](http://go-colly.org/)*: Golang scraping best library
* *[Chromedp](https://github.com/chromedp/chromedp)*: Run an headless Google Chrome instance

Not all the career pages are identical:
* *HTML*: data are stored into the HTML response, Colly is used.
* *API*: data are returned after an API call, Colly is used.
* *Javascript*: data are stored into the HTML, but after Javascript renders the page, Colly and Chromedp are used.
* *API_POST*: data are returned after an initial API call to get the cookies, Colly and Chromedp are used.
* *Pagination*: if any of the category above presents pagination, it is necessary to implement a logic for it.

### **How to create a new scraper?**
The creation of a new scraper is divided in two different part:
* Add in the database the new target and the new scraper (*Scripts/CreateScraper.go*)
    * Name, career's url and host url are necessary
```golang
func main() {
    scraper_name := "Google"
    jobs_url := "https://www.google.com/careers"
    host_url := "https://www.google.com"
    scraper := Scraper{scraper_name, jobs_url, host_url}
    scraper.CreateScraper()
}
```

* Create the algorithm to scrape in *scrapers.go*
    * Often it is good practice to build and test the scraper in a separate folder.
    * Here an example fo scraper which extract the information directly from the HTML.
```golang
func (runtime Runtime) Morressier() (results Results) {
    c := colly.NewCollector()
    start_url := "https://morressier-jobs.personio.de/"
    type Job struct {
        Url      string
        Title    string
        Location string
        Type     string
    }
    c.OnHTML("a", func(e *colly.HTMLElement) {
        if strings.Contains(e.Attr("class"), "job-box-link") {
            result_title := e.ChildText(".jb-title")
            result_url := e.Attr("href")
            result_description := e.ChildTexts("span")[0]
            result_location := e.ChildTexts("span")[2]
            results.Add(
                runtime.Name,
                result_title,
                result_url,
                result_location,
                Job{
                    result_url,
                    result_title,
                    result_location,
                    result_description,
                },
            )
        }
    })
    c.OnRequest(func(r *colly.Request) {
        fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
    })
    c.OnError(func(r *colly.Response, err error) {
        fmt.Println(Red("Request URL:"), Red(r.Request.URL))
    })
    c.Visit(start_url)
    return
}
```

### **How to run a scraper?**
```bash
go build
./JeifaiBack scrape -s=[scraper_name] -r=[true/false]
```
* -s select any scraper name
* -r true if results need to be saved, false otherwise (might be useful for testing purposes)