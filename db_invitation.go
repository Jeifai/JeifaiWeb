package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func InsertInvitation(email string, whyjoin string, whoareyou string, whichcompanies string, anythingelse string) (err error) {
	fmt.Println(Gray(8-1, "Starting InsertInvitation..."))
	statement := `INSERT INTO invitations (uuid, email, whyjoin, whoareyou, whichcompanies, anythingelse, createdat)
                  VALUES ($1, $2, $3, $4, $5, $6, current_timestamp);`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(createUUID(), email, whyjoin, whoareyou, whichcompanies, anythingelse).Scan()
	return
}

func SelectInvitationIdByUuidAndEmail(email string, uuid string) (invitation_id int) {
	fmt.Println(Gray(8-1, "Starting SelectInvitationIdByUuidAndEmail..."))
	_ = Db.QueryRow(`
					SELECT id FROM invitations
                   	WHERE email=$1 AND uuid=$2
                   	AND usedat IS NULL;`, email, uuid).Scan(&invitation_id)
	return
}

func UpdateInvitation(email string) {
	fmt.Println(Gray(8-1, "Starting UpdateInvitation..."))
	statement := `UPDATE invitations SET usedat = current_timestamp WHERE email=$1;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(email)
	if err != nil {
		panic(err.Error())
	}
}

func InsertSubscriberByEmail(email string) {
	fmt.Println(Gray(8-1, "Starting InsertSubscriberByEmail..."))
	statement := `
		INSERT INTO subscribers (email, createdat) 
		VALUES ($1, current_timestamp)
		ON CONFLICT (email) DO NOTHING;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(
		email,
	)
	if err != nil {
		panic(err.Error())
	}
}
