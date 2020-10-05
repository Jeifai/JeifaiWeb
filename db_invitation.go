package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func CreateInvitation(email string, whyjoin string, whoareyou string, whichcompanies string, anythingelse string) (err error) {
	fmt.Println(Gray(8-1, "Starting CreateInvitation..."))
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

func InvitationIdByUuidAndEmail(email string, uuid string) (err error) {
	fmt.Println(Gray(8-1, "Starting InvitationIdByUuidAndEmail..."))
	err = Db.QueryRow(`
					SELECT id FROM invitations
                   	WHERE email=$2 AND uuid=$1
                   	AND usedat IS NULL`, email, uuid).Scan()
	return
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