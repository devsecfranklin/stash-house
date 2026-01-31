package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
)
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"unique"`
    Email    string
}

func main() {
    db, err := connectToMariaDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Perform database migration
    err = db.AutoMigrate(&User{})
    if err != nil {
        log.Fatal(err)
    }

    // Create a user
    newUser := &User{Username: "john_doe", Email: "john.doe@example.com"}
    err = createUser(db, newUser)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Created User:", newUser)

    // Query user by ID
    userID := newUser.ID
    user, err := getUserByID(db, userID)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("User by ID:", user)

    // Update user
    user.Email = "updated_email@example.com"
    err = updateUser(db, user)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Updated User:", user)

    // Delete user
    err = deleteUser(db, user)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Deleted User:", user)
}

func connectToMariaDB() (*gorm.DB, error) {
    dsn := "username:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    return db, nil
}

func deleteUser(db *gorm.DB, user *User) error {
    result := db.Delete(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func createUser(db *gorm.DB, user *User) error {
    result := db.Create(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func getUserByID(db *gorm.DB, userID uint) (*User, error) {
    var user User
    result := db.First(&user, userID)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

func updateUser(db *gorm.DB, user *User) error {
    result := db.Save(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func main() {
    db, err := connectToMariaDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Perform database migration
    err = db.AutoMigrate(&User{})
    if err != nil {
        log.Fatal(err)
    }

    // Create a user
    newUser := &User{Username: "john_doe", Email: "john.doe@example.com"}
    err = createUser(db, newUser)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Created User:", newUser)

    // Query user by ID
    userID := newUser.ID
    user, err := getUserByID(db, userID)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("User by ID:", user)

    // Update user
    user.Email = "updated_email@example.com"
    err = updateUser(db, user)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Updated User:", user)

    // Delete user
    err = deleteUser(db, user)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Deleted User:", user)
}