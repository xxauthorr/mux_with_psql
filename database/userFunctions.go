package database

import (
	"Login_Admin_Using_Postgres/models"
	"Login_Admin_Using_Postgres/utils"
	"fmt"
	"log"
)

func RegisterUser(email, pass string) {
	encryptdPass, err := utils.HashEncrypt(pass)
	if err != nil {
		log.Fatal("Encryption err - ", err)
	}
	insertStmt := `INSERT INTO "clientuser"("email", "hashpass") VALUES($1, $2)`
	_, insertErr := Db.Exec(insertStmt, email, encryptdPass)
	if insertErr != nil {
		log.Fatal("register user db insert err : ", insertErr)
	}
}

func CheckEmail(formEmail string) (string, bool) {
	if formEmail == "" {
		fmt.Println("formEmail is nil : func CheckEmail(.. string)(string, bool)")
		return "", false
	}
	tbl := models.ClientUser{}
	getStmt := `SELECT email,hashpass FROM clientuser WHERE "email" = $1;`
	rows, err := Db.Query(getStmt, formEmail)
	if err != nil {
		log.Fatal("CheckEmail err :", err)
		return "", false
	}
	defer rows.Close()
	for rows.Next() {
		if err1 := rows.Scan(&tbl.Email, &tbl.Hashpass); err1 != nil {
			log.Fatal("hashpass retrieve error", err1)
			return "", false
		}
	}
	if tbl.Email == "" {
		fmt.Println("Email,(", formEmail, ") is not in clientUser Database")
		return "", false
	}

	fmt.Println("Hash-pass of ", tbl.Email, ":", tbl.Hashpass)
	return tbl.Hashpass, true
}
