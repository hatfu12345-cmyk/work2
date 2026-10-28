package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	dsn := "host=" + env("DB_HOST", "db") +
		" port=" + env("DB_PORT", "5432") +
		" user=" + env("DB_USER", "app") +
		" password=" + env("DB_PASSWORD", "secret123") +
		" dbname=" + env("DB_NAME", "appdb") +
		" sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Waiting for DB... (%v)", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS items (id SERIAL PRIMARY KEY, name TEXT NOT NULL)"); err != nil {
		log.Fatal(err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	if count == 0 {
		db.Exec("INSERT INTO items (name) VALUES ('first item'), ('second item')")
	}

	staticDir := env("STATIC_DIR", "/app/static")
	fs := http.FileServer(http.Dir(staticDir))

	http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM items ORDER BY id")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		items := []Item{}
		for rows.Next() {
			var it Item
			if err := rows.Scan(&it.ID, &it.Name); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			items = append(items, it)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	http.Handle("/", fs)

	addr := ":" + env("HTTP_PORT", "8080")
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
