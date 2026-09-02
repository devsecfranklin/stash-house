package main
import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)
var db *sql.DB
func main() {
    db = connectDB()
    defer db.Close()
    router := mux.NewRouter()
    router.HandleFunc("/users", getUsers).Methods("GET")
    router.HandleFunc("/users/{id}", getUser).Methods("GET")
    router.HandleFunc("/users", createUser).Methods("POST")
    router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
    router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":8000", router))
}
func connectDB() *sql.DB {
    // Replace "username", "password", "dbname" with your database credentials
    connectionString := "username:password@tcp(localhost:3306)/dbname"
    db, err := sql.Open("mysql", connectionString)
    if err != nil {
        log.Fatal(err)
    }
    return db
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()
    // Iterate through the rows and create a JSON response
    // You can use a JSON library like "encoding/json" to marshal the data
}
func getUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["id"]
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", userID)
    // Parse the row data and create a JSON response
}

func createUser(w http.ResponseWriter, r *http.Request) {
    // Parse the request body and extract user data
    // Insert the user data into the "users" table
    // Handle validation and error cases
}
func updateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["id"]
    // Parse the request body and extract user data
    // Update the user data in the "users" table based on the user ID
    // Handle validation and error cases
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["id"]
    _, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Send a success response
}
