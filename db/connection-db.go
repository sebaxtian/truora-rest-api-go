package connectiondb

import (
	"database/sql"
	"encoding/json"
	"fmt"

	// "log"
	// User Go pq driver
	_ "github.com/lib/pq"

	infoserver "github.com/sebaxtian/truora-rest-api-go/structs"
)

// DBConnect create a connection db
func DBConnect() *sql.DB {
	// Connect to the "truora" database.
	db, err := sql.Open("postgres", "postgresql://truora@localhost:26257/truora?sslmode=disable")
	if err != nil {
		// log.Fatal(err)
		fmt.Println("ERROR connecting to the database: ", err)
	}
	// Create the "accounts" table.
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS infoservers (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), domain STRING, data_infoserver JSONB, last_updated STRING)"); err != nil {
		// log.Fatal(err)
		fmt.Println("ERROR create table info_server: ", err)
	}
	return db
}

// GetInfoServer select and get InfoServer
func GetInfoServer(domain string, db *sql.DB) infoserver.InfoServer {
	// Print out the infoserver.
	rows, err := db.Query("SELECT * FROM infoservers WHERE domain = '" + domain + "'")
	if err != nil {
		// log.Fatal(err)
		fmt.Println("ERROR select info_server: ", err)
	}
	defer rows.Close()
	// fmt.Println("infoserver:")
	var infoServer infoserver.InfoServer
	var id string
	var domainName, dataInfoserver, lastUpdated string
	for rows.Next() {
		if err := rows.Scan(&id, &domainName, &dataInfoserver, &lastUpdated); err != nil {
			// log.Fatal(err)
			fmt.Println("ERROR get data for result set: ", err)
		}
		// fmt.Printf("%s %s %s %s\n", id, domain, dataInfoserver, lastUpdated)
		bytes := []byte(dataInfoserver)
		err := json.Unmarshal(bytes, &infoServer)
		if err != nil {
			// panic(err)
			fmt.Println("ERROR create data infoServer: ", err)
		}
		infoServer.ID = id
		infoServer.LastUpdated = lastUpdated
	}
	// fmt.Println("Create data infoServer with data table infoserver: ", infoServer)
	return infoServer
}

// CreateInfoServer create infoServer into the infoservers table
func CreateInfoServer(domain string, infoServer infoserver.InfoServer, db *sql.DB) bool {
	infoServerStr, err := json.Marshal(infoServer)
	if err != nil {
		// panic(err)
		fmt.Println("ERROR parse data infoServer to string: ", err)
		return false
	}
	// fmt.Println(string(infoServerStr))
	// Insert one row into the "infoservers" table.
	if _, err := db.Exec(
		"INSERT INTO infoservers (domain, data_infoserver, last_updated) VALUES ('" + domain + "', '" + string(infoServerStr) + "', '" + infoServer.LastUpdated + "')"); err != nil {
		// log.Fatal(err)
		fmt.Println("ERROR Insert one row infoservers table: ", err)
		return false
	}
	return true
}

// UpdateInfoServer update infoServer into the infoservers table
func UpdateInfoServer(domain string, infoServer infoserver.InfoServer, db *sql.DB) bool {
	infoServerStr, err := json.Marshal(infoServer)
	if err != nil {
		// panic(err)
		fmt.Println("ERROR parse data infoServer to string: ", err)
		return false
	}
	// fmt.Println(string(infoServerStr))
	// Insert one row into the "infoservers" table.
	if _, err := db.Exec(
		"UPDATE infoservers SET (data_infoserver, last_updated) = ('" + string(infoServerStr) + "', '" + infoServer.LastUpdated + "') WHERE domain = '" + domain + "'"); err != nil {
		// log.Fatal(err)
		fmt.Println("ERROR Update one row infoservers table: ", err)
		return false
	}
	return true
}
