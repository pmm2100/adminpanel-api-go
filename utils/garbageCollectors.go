package utils

import (
	"fmt"
	"time"

	"adminPanel.com/db"
)

func StartMailGarbageCollector() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)

		for range ticker.C {

			result, err := db.DB.Exec("DELETE FROM verify WHERE date<=datetime('now')")
			if err != nil {
				fmt.Printf("Could not delete expired codes, %v", err)
			} else {
				rows, _ := result.RowsAffected()
				if rows > 0 {
					fmt.Printf("Mail garbage collector affected %v rows", rows)
				}
			}
		}
	}()
}

func StartUnverifiedUsersGarbagerCollector() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)

		for range ticker.C {
			result, err := db.DB.Exec("DELETE FROM users WHERE date<=datetime('now') AND verified=0")
			if err != nil {
				fmt.Printf("Could not delete not verified accounts, %v", err)
			} else {
				rows, _ := result.RowsAffected()
				if rows > 0 {
					fmt.Printf("User garbage collector affected %v rows", rows)
				}
			}
		}
	}()
}
