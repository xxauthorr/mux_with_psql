package database

import (
	"Login_Admin_Using_Postgres/models"
	"fmt"
	"log"
)

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

func GetUser(userId string) (string, bool) {
	var userEmail string
	getStmt := `SELECT email FROM clientuser WHERE "id" = $1;`
	rows, err := Db.Query(getStmt, userId)
	if err != nil {
		log.Fatal("GetUser err :", err)
		return "", false
	}
	defer rows.Close()
	for rows.Next() {
		if err1 := rows.Scan(&userEmail); err1 != nil {
			log.Fatal("hashpass retrieve error", err1.Error())
			return "", false
		}
	}
	return userEmail, true
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
	var ClientData models.Sample
	var Data = models.ClientUser{}
	rows, err := Db.Query(`SELECT id,email FROM clientuser ORDER BY id asc`)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&tableId, &tableEmail); err != nil {
			log.Fatal(err.Error())
		}
		Data = models.ClientUser{
			Id:    tableId,
			Email: tableEmail,
		}
		ClientData.Data = append(ClientData.Data, Data)
	}
	return ClientData
}

func UpdateUserdata(userId, newEmail string) bool {
	insertStmt := `UPDATE clientuser SET email = $1 WHERE id = $2`
	_, updateErr := Db.Exec(insertStmt, newEmail, userId)
	if updateErr != nil {
		log.Fatal(updateErr.Error())
		return false
	}
	return true
}
