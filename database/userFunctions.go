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

func CheckUserEmail(formEmail string) (string, bool) {
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
func CheckAdminEmail(userName string) (string, bool) {
	if userName == "" {
		fmt.Println("formEmail is nil : func CheckEmail(.. string)(string, bool)")
		return "", false
	}
	tbl := models.ClientUser{}
	getStmt := `SELECT username,hashpass FROM adminuser WHERE "username" = $1;`
	rows, err := Db.Query(getStmt, userName)
	if err != nil {
		log.Fatal("CheckEmail err :", err)
		return "", false
	}
	defer rows.Close()
	for rows.Next() {
		if err1 := rows.Scan(&tbl.Email, &tbl.Hashpass); err1 != nil {
			log.Fatal("hashpass retrieve error", err1.Error())
			return "", false
		}
	}

	if tbl.Email == "" {
		fmt.Println("Email,(", userName, ") is not in clientUser Database")
		return "", false
	}

	fmt.Println("Hash-pass of ", tbl.Email, ":", tbl.Hashpass)
	return tbl.Hashpass, true
}

func DeleteUser(id string) bool {
	delStmt := `DELETE FROM clientuser WHERE "id" = $1;`
	_, err := Db.Exec(delStmt, id)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	return true
}

func EditUser(id, newEmail string) bool {
	_, err := Db.Exec(`UPDATE clientuser SET email = $1`, newEmail)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	return true
}

func FetchUserData() models.Sample {
	var (
		tableId, tableEmail string
	)

	rows, err := Db.Query(`SELECT id,email FROM clientuser`)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()
	// var Data = models.ClientUser{}
	var data models.Sample
	for rows.Next() {
		if err := rows.Scan(&tableId, &tableEmail); err != nil {
			log.Fatal(err.Error())
		}
		Data := models.ClientUser{
			Id:    tableId,
			Email: tableEmail,
		}
		data.Data = append(data.Data, Data)
	}
	return data
}
