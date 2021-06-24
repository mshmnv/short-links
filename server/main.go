package main

// originally generated link
// post url which is already in db

// bd
// - Create - check url in bd - if yes -> give link from bd
//                        if no  -> generate -> check its uniqueness
// - Get - check link - if yes -> give url from bd
//                    - if no -> say "none"

// database short_db
// table short
// post 5432

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	ps "proto"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

var db *sql.DB

func GetDB() *sql.DB {
	var err error
	if db == nil {
		db, err = sql.Open("postgres",
			"host=127.0.0.1 port=5432 dbname=short_db sslmode=disable user=postgres password=postgres")
		if err != nil {
			log.Fatalf("Unable to connect to database: %s", err)
		}
	}
	return db
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateLink() string {
	var unique = false
	link := make([]byte, 10)
	for !unique {
		for i := range link {
			link[i] = charset[seededRand.Intn(len(charset))]
		}
		unique = checkUniqueness(string(link))
	}
	return string(link)
}

func checkUniqueness(link string) bool {
	var unique bool = true

	db = GetDB()
	rows, err := db.Query("SELECT link FROM short WHERE link = $1;", link)
	if err != nil {
		log.Fatalf("Unable to select from table: %s", err)
	}
	for rows.Next() {
		unique = false
		rows.Close()
	}
	return unique
}

type ShortServer struct {
}

func (s *ShortServer) Create(ctx context.Context, req *ps.UrlRequest) (*ps.LinkResponse, error) {
	var err error
	response := new(ps.LinkResponse)
	link := GenerateLink()
	db = GetDB()
	var toInsert bool = true

	rows, err := db.Query("SELECT link FROM short WHERE url = $1;", req.Url)
	if err != nil {
		log.Fatalf("Unable to select from table: %s", err)
	}
	for rows.Next() {
		toInsert = false
		err = rows.Scan(&link)
		if err != nil {
			log.Fatalf("Unable to scan url after select: %s", err)
		}
	}
	if toInsert {
		_, err = db.Exec("INSERT INTO short VALUES ($1, $2);", req.Url, link)
		if err != nil {
			log.Fatalf("Unable to insert into table: %s", err)
		}
	}
	response.Link = "short.com/" + link
	return response, err
}

func (s *ShortServer) Get(ctx context.Context, req *ps.LinkRequest) (*ps.UrlResponse, error) {
	var err error
	response := new(ps.UrlResponse)

	if strings.Contains(req.Link, "short.com/") {
		req.Link = strings.Replace(req.Link, "short.com/", "", 1)
	}
	response.Url = "none"

	db = GetDB()
	rows, err := db.Query("SELECT url FROM short WHERE link = $1;", req.Link)
	if err != nil {
		log.Fatalf("Unable to select from table: %s", err)
	}
	for rows.Next() {
		err = rows.Scan(&response.Url)
		if err != nil {
			log.Fatalf("Unable to scan url after select: %s", err)
		}
	}
	return response, err
}

func main() {
	// ...   connecting to database   ...
	db, err := sql.Open("postgres",
		"host=127.0.0.1 port=5432 dbname=short_db user=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " +
		`short("url" varchar(255),` +
		`"link" varchar(10));`)
	if err != nil {
		log.Fatalf("Unable to create table: %s", err)
	}

	// ...   initializing server   ...
	server := grpc.NewServer()
	instance := new(ShortServer)
	ps.RegisterShortServer(server, instance)

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Fatalf("Unable to create grpc listener: %s", err)
	}
	// ...   starting the server   ...
	if err = server.Serve(listener); err != nil {
		log.Fatalf("Unable to start server: %s", err)
	}
}
