package database

import database "main.go/database/models"

// Database provides an abstraction for the database operations
type Database interface {
	database.PostDatabase
	database.UserDatabase
}
