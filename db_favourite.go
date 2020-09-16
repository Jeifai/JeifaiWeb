package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func (user *User) InsertFavourite(resultid int) {
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
		resultid,
	)
}

func (user *User) DeleteFavourite(resultid int) {
	fmt.Println(Gray(8-1, "Starting InsertFavourite..."))
	statement := `
                    UPDATE favouriteresults
                    SET deletedat = current_timestamp
                    WHERE userid = $1
                    AND resultid = $2
                    AND deletedat IS NULL;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(
		&user.Id,
		resultid,
	)
}
