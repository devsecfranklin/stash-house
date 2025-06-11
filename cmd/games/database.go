/*
sudo apt-get install sqlite3 libsqlite3-dev gcc
cd sites/games
go get github.com/mattn/go-sqlite3
*/
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // If a package is imported with a blank identifier, the package's init function is called. The driver is registered using this function.
)


type Table struct {
	name string
	db string
	exist bool
}

type User struct {
	id     int
	name   string
	email  string
}

func ConnectDB() {
	db, err := sql.Open("sqlite3", ":memory:") // open a database specified by database driver name and a driver-specific data source name

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // close the connection

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version) // executes a query that is expected to return at most one row. The column from the matched row is copied into the version variable by the Scan function

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version)
}

/*
func ReadDB() {
	// Method implementation
}

func WriteDB() {
	// Method implementation
}
*/

/*
Perform CRUD Operation using golang sqlite driver
Create a database table and Read from the Table
*/

func CrudDB(t *Table) bool {
	db, err := sql.Open("sqlite3", t.db)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Create table
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS " + t.name + " (id INTEGER PRIMARY KEY, name VARCHAR(64), email VARCHAR(64), password VARCHAR(32) NULL)")
	if err != nil {
		log.Println("Error in creating table")
		return false
	} else {
		log.Println("Successfully created table " + t.name + "  in db games")
	}
	statement.Exec()

	// Create
	statement, _ = db.Prepare("INSERT INTO " + t.name + " (name, email, password) VALUES (?, ?, ?)")
	statement.Exec("franklin", "frank378@gmail.com", "66pickUPstickSS")
	log.Println("Inserted initial user into database")

	// Read
	rows, _ := db.Query("SELECT id, name, email FROM games")
	var tempUser User
	for rows.Next() {
		rows.Scan(&tempUser.id, &tempUser.name, &tempUser.email)
		log.Printf("ID:%d, User:%s, email:%s\n", tempUser.id,
			tempUser.name, tempUser.email)
	}
	return true
}

func check_table_exist(t *Table) bool {

	db, err := sql.Open("sqlite3", t.db)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	table := "users"
	_, table_check := db.Query("select * from " + table + ";")

    if table_check == nil {
        fmt.Println("table " + t.name + " is there")
		log.Println("table " + t.name + " is there")
		return true
    } else {
        fmt.Println("table " + t.name + " is not there")
		log.Println("table " + t.name + " is not there")
		return false
	}
}

/*
	db, err := sql.Open("sqlite3", "books.db")
	if err != nil {
		log.Println(err)
	}

	// Create table
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS books (id INTEGER PRIMARY KEY, isbn INTEGER, author VARCHAR(64), name VARCHAR(64) NULL)")
	if err != nil {
		log.Println("Error in creating table")
	} else {
		log.Println("Successfully created table books!")
	}
	statement.Exec()

	// Create
	statement, _ = db.Prepare("INSERT INTO books (name, author, isbn) VALUES (?, ?, ?)")
	statement.Exec("A Tale of Two Cities", "Charles Dickens", 140430547)
	log.Println("Inserted the book into database!")

	// Read before Update
	rows, _ := db.Query("SELECT id, name, author FROM books")
	var tempBook Book
	for rows.Next() {
		rows.Scan(&tempBook.id, &tempBook.name, &tempBook.author)
		log.Printf("ID:%d, Book:%s, Author:%s\n", tempBook.id,
			tempBook.name, tempBook.author)
	}

	// Update
	statement, _ = db.Prepare("update books set name=? where id=?")
	statement.Exec("A Tale of Three Cities", 1)
	log.Println("Successfully updated the book in database!")

	// Read after Update
	rows, _ = db.Query("SELECT id, name, author FROM books")

	for rows.Next() {
		rows.Scan(&tempBook.id, &tempBook.name, &tempBook.author)
		log.Printf("ID:%d, Book:%s, Author:%s\n", tempBook.id,
			tempBook.name, tempBook.author)
	}

	// Delete
	statement, _ = db.Prepare("delete from books where id=?")
	statement.Exec(1)
	log.Println("Successfully deleted the book in database!")
*/

/*

Different go sqlite driver functions
Go sqlite3 Exec
The Exec function executes a query without returning any rows. First of all, we run a query to create 'students' table:

go

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	sts := `
DROP TABLE IF EXISTS students;
CREATE TABLE students(id INTEGER PRIMARY KEY, name TEXT, score REAL);
INSERT INTO students(name, score) VALUES('Anna',8.5);
INSERT INTO students(name, score) VALUES('Bob',7.5);
INSERT INTO students(name, score) VALUES('Claire',9.5);
INSERT INTO students(name, score) VALUES('Charlie',6.5);
INSERT INTO students(name, score) VALUES('Daniel',8.0);
INSERT INTO students(name, score) VALUES('Hellen',7.0);
INSERT INTO students(name, score) VALUES('Hummer',7.5);
INSERT INTO students(name, score) VALUES('John',10);
`
	// run the query
	_, err = db.Exec(sts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("table created")
}
*/

/*
Select rows with Query
The Query method runs a SELECT query that returns rows. The optional arguments are for any query placeholder parameters. Here's an example of query all students who have score > 8:

go

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM students where score > 8")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
        
        // iterate through all the records
	for rows.Next() {
		var id int
		var name string
		var score float64
		err = rows.Scan(&id, &name, &score)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v %v %v\n", id, name, score)
	}
}
*/