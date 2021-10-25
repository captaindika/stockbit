package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
	"database/sql"
	"context"
    "time"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
)

const (  
    username = "user"
    password = "password"
    hostname = "127.0.0.1:3306"
    dbname   = "db"
)

type search struct {  
    search      string
}

type Movie struct {
	Title string `json:"Title"`
	Year string `json:"Year"`
	imdbID string `json:"imdbID"`
	Type string `json:"Type"`
	Poster string `json:"Poster"`
}

var Movies []Movie

func omdbAPI(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("http://www.omdbapi.com/?apikey=faf7e5bb&s=Batman&page=2")
    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
	json.NewEncoder(w).Encode(responseData)
    // fmt.Println(string(responseData))
}
func dsn(dbName string) string {  
    return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func dbConnection() (*sql.DB, error) {  
    db, err := sql.Open("mysql", dsn(""))
    if err != nil {
        log.Printf("Error %s when opening DB\n", err)
        return nil, err
    }
    //defer db.Close()

    ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelfunc()
    res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
    if err != nil {
        log.Printf("Error %s when creating DB\n", err)
        return nil, err
    }
    no, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error %s when fetching rows", err)
        return nil, err
    }
    log.Printf("rows affected %d\n", no)

    db.Close()
    db, err = sql.Open("mysql", dsn(dbname))
    if err != nil {
        log.Printf("Error %s when opening DB", err)
        return nil, err
    }
    //defer db.Close()

    db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(20)
    db.SetConnMaxLifetime(time.Minute * 5)

    ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelfunc()
    err = db.PingContext(ctx)
    if err != nil {
        log.Printf("Errors %s pinging DB", err)
        return nil, err
    }
    log.Printf("Connected to DB %s successfully\n", dbname)
    return db, nil
}

func createLogTable(db *sql.DB) error {  
    query := `CREATE TABLE IF NOT EXISTS logs(id int primary key auto_increment, search text, 
	created_at datetime default CURRENT_TIMESTAMP)`
    ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelfunc()
    res, err := db.ExecContext(ctx, query)
    if err != nil {
        log.Printf("Error %s when creating product table", err)
        return err
    }
    rows, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error %s when getting rows affected", err)
        return err
    }
    log.Printf("Rows affected when creating table: %d", rows)
    return nil
}

func insertDB(db *sql.DB, s search) error {
	query := "INSERT INTO logs(search) VALUES(?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statemtent", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.search)
    if err != nil {
        log.Printf("Error %s when inserting row into products table", err)
        return err
    }
    rows, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error %s when finding rows affected", err)
        return err
    }
    log.Printf("%d products created ", rows)
    return nil
}



func main() {
	s := search{
		search: "Spiderman",
	}  
    db, err := dbConnection()
    if err != nil {
        log.Printf("Error %s when getting db connection", err)
        return
    }
    defer db.Close()
    log.Printf("Successfully connected to database")
    err = createLogTable(db)
    if err != nil {
        log.Printf("Create log table failed with error %s", err)
        return
    }
	err = insertDB(db, s)
	if err != nil {
		log.Printf("INSERT logs failed with error %s", err)
		return
	}
	http.HandleFunc("/movie", omdbAPI)
	log.Fatal(http.ListenAndServe(":9000",nil))
}