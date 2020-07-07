package main

import (
	"fmt"
	"time"
)

type Invitation struct {
    Id              int
    Uuid            string
	Email           string
	Whyjoin         string
	Whichcompanies  string
	Anythingelse    string
}

func (invitation *Invitation) InvitationIdByEmail() (err error) {
	fmt.Println("Starting InvitationIdByEmail...")
	err = Db.QueryRow(`SELECT
                         i.id
                       FROM invitations i
                       WHERE i.email=$1`, invitation.Email).Scan(&invitation.Id)
	return
}

func (invitation *Invitation) InvitationIdByUuidAndEmail() (err error) {
	fmt.Println("Starting InvitationIdByUuidAndEmail...")
    err = Db.QueryRow(`SELECT i.id
                       FROM invitations i
                       WHERE i.uuid=$1
                       AND i.email=$2 
                       AND i.usedat IS NULL`, 
        invitation.Uuid, 
        invitation.Email).Scan(&invitation.Id)
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

func (invitation *Invitation) UpdateInvitation() {
	statement := `UPDATE invitations SET usedat = current_timestamp WHERE id=$1;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
    _, err = stmt.Exec(invitation.Id)
    if err != nil {
		panic(err.Error())
    }
}
