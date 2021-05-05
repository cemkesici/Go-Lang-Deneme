package main

import (
	"database/sql"
	"fmt"
	_ "go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type users struct {
	iduser  int
	eMail   string
	lName   string
	fName   string
	Message string
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := sql.Open("mysql", "root:123456@/users")
	checkError(err)
	defer db.Close()

	createStatement := `
	    user (
		iduser int(11) NOT NULL AUTO_INCREMENT,
		lName varchar(45) NOT NULL,
		fName varchar(45) NOT NULL,
		eMail varchar(255) NOT NULL,
		Message longtext DEFAULT NULL,
		PRIMARY KEY (iduser)
	  ) ;`

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS" + createStatement); err != nil {
		log.Fatal(err)
	}
	var (
		iduser  int
		eMail   string
		lName   string
		fName   string
		Message string
	)

	// Eklenen kayıtları getir
	rows, err := db.Query("SELECT * FROM user")
	checkError(err)

	for rows.Next() {
		err = rows.Scan(&iduser, &lName, &fName, &eMail, &Message)
		checkError(err)
		log.Printf("Bulunan satır içeriği : %q ", strconv.Itoa(iduser)+" "+lName+" "+fName+" "+eMail+" "+Message)
	}

	tmpl := template.Must(template.ParseFiles("forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		lName := r.FormValue("isim")
		fName := r.FormValue("soyisim")
		eMail := r.FormValue("email")
		Message := r.FormValue("message")

		res, err := db.Exec(`INSERT INTO user (fName, lName, eMail,Message) VALUES (?, ?, ?, ?)`, fName, lName, eMail, Message)
		checkError(err)

		id, err := res.LastInsertId()
		fmt.Println(id)
		checkError(err)
		rowCount, err := res.RowsAffected()
		checkError(err)
		log.Printf("İnserted  %d rows ", rowCount)

		tmpl.Execute(w, struct{ Success bool }{true})
	})

	http.ListenAndServe(":9000", nil)
}
