package database

import (
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) {
	for range 10 {

		m, err := migrate.New(
			"file:///app/migrations",
			dbURL,
		)
		if err != nil {
			log.Printf("migration init failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = m.Up()

		if err == nil || err == migrate.ErrNoChange {
			log.Println("migrations applied successfully")
			return
		}

		log.Printf("migration attempt failed: %v", err)

		time.Sleep(2 * time.Second)
	}

	log.Fatal("migration failed after retries")
}
