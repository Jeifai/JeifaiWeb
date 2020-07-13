package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

type HomeData struct {
    Targets     int
    Keywords    int
}

func (user *User) GetHomeData() (home HomeData, err error) {
	fmt.Println(Gray(8-1, "Starting GetHomeData..."))
	err = Db.QueryRow(`SELECT
                        COUNT(DISTINCT utk.targetid),
                        COUNT(DISTINCT utk.keywordid)
                      FROM userstargetskeywords utk
                      WHERE utk.userid=$1
                      AND utk.deletedat IS NULL`,
		user.Id,
	).
		Scan(
			&home.Targets,
			&home.Keywords,
		)
	return
}