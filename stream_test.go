package dbstream_test

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	"github.com/cybernostics/dbstream"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	_ "github.com/mattn/go-sqlite3"
)

// Given a Struct for our row return...
type Person struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IpAddress string `json:"ip_address"`
}

// ...and a function that scans an sql.Rows object into that struct
func scanPerson(rows *sql.Rows, p *Person) error {
	return rows.Scan(&p.Id, &p.FirstName, &p.LastName, &p.Email, &p.IpAddress)
}

func TestStreamCreate(t *testing.T) {
	var db *sql.DB
	db, err := sql.Open("sqlite3", "./names.db")
	then.AssertThat(t, err, is.Nil())

	count := 3
	rows, err := db.Query("SELECT id, first_name, last_name, email, ip_address from people LIMIT " + strconv.Itoa(count))

	then.AssertThat(t, err, is.Nil())

	// create a stream from sql.Rows
	rowStream := dbstream.RowStreamFrom(rows, scanPerson)
	defer rowStream.Close()

	for rowStream.Next() {
		p, err := rowStream.Current()
		then.AssertThat(t, err, is.Nil())
		fmt.Printf("%v", p)
	}

}
