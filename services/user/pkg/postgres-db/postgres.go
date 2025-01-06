package postgresdb

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func InitPostgresClient(connectionString string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Printf("PostgreSQL bağlantı denemesi %d başarısız: %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("PostgreSQL'e başarıyla bağlanıldı!")
			return db, nil
		}

		log.Printf("PostgreSQL ping denemesi %d başarısız: %v", i+1, err)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("PostgreSQL'e bağlanılamadı: %v", err)
}
