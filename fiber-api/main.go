package main

import (
	"database/sql"
	"fiber-api/contact"
	"fiber-api/database"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	app := fiber.New()
	log.Printf("Staring web server...\n")
	var contactDao = contact.ContactDaoS{
		Connection: database.Connection,
	}
	contact.RegisterRouterHook(app)(&contactDao)
	defer func(Connection *sql.DB) {
		err := Connection.Close()
		if err != nil {
			log.Fatal("Error during db connection closing")
		}
	}(database.Connection)
	log.Fatal(app.Listen(":3000"))
}
