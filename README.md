# JaifaiWeb

Here everything about the web application of Jaifai.

As any webapp JeifaiWeb has two sides:

* *Backend*: using [Golang](https://golang.org/)
    * *Main*, the starting point.
    * *Route*, coordinator between the main and the database.
    * *DB*, all the queries to manage and manipulate data.
* *Frontend*: using [Vuejs](https://vuejs.org/)
    * *Templates*, all the pure HTML files.
    * *Public*, the Vue code is stored here, as well as the CSS and images.

## Backend

A typical scenario:

* User is inside the webapp and opens the page dedicated to his Profile.
* *main.go* receives the request.
* *main.go* calls the function *Profile()* in the file *route_profile.go*.
* *route_profile.go* receives the request.
* *route_profile.go* call the function *UserById()* in the file *db_user.go*.
* *db_user.go* receives the request.
* *db_user.go* execute the query, get the data.
* *db_user.go* returns the data to *route_profile.go*.
* *route_profile.go* receives the data.
* *route_profile.go* returns them to the user.
* User visualises the data in the webapp.

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

### DB

Any interaction with the database happens with a file called *db...go*, for example *db_user.go*.
The purpose of these files is to keep ordered all the queries written to dialogue with the database.

For example:

```golang
func (keyword *Keyword) KeywordByText() (err error) {
    
    // Here a simple query returning a unique value.
    // Given a keyword's text, return its unique id.
    err = Db.QueryRow(`SELECT
                         k.id
                       FROM keywords k
                       WHERE k.text=$1`, keyword.Text).Scan(&keyword.Id)
    return
}
```

## Fontend

The frontend is divided in two parts:
* If a user is logged out: render the HTML files from the folder *templates*.
* If a user is logged in: render the Vue code fro the folder *public/js*.

### Vuejs

The great advantage of using a frontend framework like Vuejs is that it does not need to render the page multiple times.
The server returns only one all the code necessary, then Javascript renders the page.
All the upcoming interaction between client and server are then managed using specific API endpoints.

The Vuejs applications is structured in this way:
* *main.js* file which acts as router.
* *components* folder which contains all the different sections of the webapp:
    * *Home.js*
    * *Analytics.js*
    * *Keywords.js*
    * *Matches.js*
    * *Profile.js*
    * *Targets.js*

An example of how the *main.js* file is structured:

```javascript
import Home from './components/Home.js';

const router = new VueRouter({
    mode: 'hash',
    routes: [
        {
            path: '/',
            component: Home
        },
    ]
})

var app = new Vue({
    router,
    delimiters: ["[[","]]"],
    data() {
        return {
            data1: "example",
            data2: 25,
            data3: [1, 2, 3],
        }
    },
    watch: {
        data1: function(val) {
            data2 = 30;
        }
    },
    methods: {
        csvExport(arrData, fileName) {}
    },
    mounted() {
        let styleElem = document.createElement('style');
        styleElem.textContent = `
            .body {
                font-size: 16px;
            }`
        document.head.appendChild(styleElem);
    }
})
app.$mount('#app');
```

An example of component:

```javascript
export default {
    name: 'home',
    delimiters: ["[[","]]"],
    data: function () {
        return {
            data1: 0,
        }
    },
    mounted() {
        this.data1 = 1;
    },
    created () {
        this.fetchData()
    },
    methods: {
        fetchData: function() {
            this.$http.get('/fetchData').then(function(response) {
                this.data1 = response.data.Test;
            }).catch(function(error) {
                console.log(error)
            });
        }
    },
    template: `
        <div>
            <span v-html="data1"></span>
        </div>
    `,
};
```

## TODO
* How to run the webapp locally
* Write architecture_decision_record
* How to deploy
* Discuss emailing
* Dependencies (which libraries and why)