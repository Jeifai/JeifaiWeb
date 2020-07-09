package main

import (
	"fmt"
    "time"
    
	. "github.com/logrusorgru/aurora"
)

type Invitation struct {
	Id             int
	Uuid           string
	Email          string
	Whyjoin        string
	Whichcompanies string
	Anythingelse   string
}

func (invitation *Invitation) InvitationIdByEmail() (err error) {
	fmt.Println(Gray(8-1, "Starting InvitationIdByEmail..."))
	err = Db.QueryRow(`SELECT
                         i.id
                       FROM invitations i
                       WHERE i.email=$1`, invitation.Email).Scan(&invitation.Id)
	return
}

func (user *User) InvitationIdByUuidAndEmail() (invitation Invitation) {
	fmt.Println(Gray(8-1, "Starting InvitationIdByUuidAndEmail..."))
	_ = Db.QueryRow(`SELECT i.id
                       FROM invitations i
                       WHERE i.uuid=$1
                       AND i.email=$2 
                       AND i.usedat IS NULL`,
		user.InvitationCode,
		user.Email).Scan(&invitation.Id)
	return
}

func (invitation *Invitation) CreateInvitation() {
	fmt.Println(Gray(8-1, "Starting CreateInvitation..."))
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
	fmt.Println(Gray(8-1, "Starting UpdateInvitation..."))
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
