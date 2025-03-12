package main

import (
	"log"
	"net/http"
)

func main() {
	db := SetupDatabase()
	defer db.Close()

	http.HandleFunc("/jmap", JMAPHandler(db))

	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}