package main

import (
	"fmt"
	"time"
)

type Invitation struct {
	Id             int
	Email          string
	Whyjoin        string
	Whichcompanies string
	Anythingelse   string
}

func (invitation *Invitation) InvitationByEmail() (err error) {
	fmt.Println("Starting InvitationByEmail...")
	err = Db.QueryRow(`SELECT
                         i.id
                       FROM invitations i
                       WHERE i.email=$1`, invitation.Email).Scan(&invitation.Id)
	return
}

func (invitation *Invitation) CreateInvitation() {
	fmt.Println("Starting CreateInvitation...")
	statement := `INSERT INTO invitations (uuid, email, whyjoin, whichcompanies, anythingelse, createdat)
                  VALUES ($1, $2, $3, $4, $5, $6)`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(
		createUUID(),
		invitation.Email,
		invitation.Whyjoin,
		invitation.Whichcompanies,
		invitation.Anythingelse,
		time.Now(),
	)
}
