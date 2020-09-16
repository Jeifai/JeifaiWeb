package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

type Result struct {
    Url string
    Id int
}

func (user *User) InsertFavourite(result Result) {
	fmt.Println(Gray(8-1, "Starting InsertFavourite..."))
    statement := `INSERT INTO favouriteresults (userid, resultid, createdat)
                  VALUES ($1, $2, current_timestamp);`
    stmt, err := Db.Prepare(statement)
    if err != nil {
        panic(err.Error())
    }
    defer stmt.Close()
    stmt.QueryRow(
        &user.Id,
        result.Id,
    )
}
