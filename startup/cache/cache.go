package cache

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"log"
)

//FTC - First Time Cache.  Used to check if this app has run previously
type FTC struct {
	id      string
	cacheDB *sql.DB
}

func (f *FTC) New(db string) (*FTC, error) {

	firstRunTable := "CREATE TABLE FIRSTRUN (id varchar(50) not null primary key);"

	cachedb, err := sql.Open("sqlite3", db)
	if err != nil {
		log.Fatalln("could not load database ", err)
	}

	if _, err = cachedb.Exec(firstRunTable); err != nil {
		log.Println("could not create new table")
		return &FTC{}, err
	}

	id := uuid.NewV4().String()

	return &FTC{
		id:      id,
		cacheDB: cachedb,
	}, nil

}

func (f *FTC) Write() error {

	insertStatement := `insert into FIRSTRUN(id) VALUES (?);`

	stmt, err := f.cacheDB.Prepare(insertStatement)
	if err != nil {
		log.Println("failed to write uuid to table ", err)
		return err
	}

	res, err := stmt.Exec(f.id)
	if err != nil {
		log.Println("issue happened when trying to execute statement ", err)
	}

	fmt.Print(res.LastInsertId())

	return nil

}
