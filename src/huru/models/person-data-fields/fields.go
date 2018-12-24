package person_data_fields

import (
	"huru/dbs"
	"sync"
)

//personalEmail text,
//businessEmail text,
//facebook text,
//instagram text

// Model The person Type (more like an object)
type Model struct {
	ID       int    `json:"id, omitempty"`
	PersonId int    `json:"person_id, omitempty"`
	Key      string `json:"key, omitempty"`
	Value    string `json:"value, omitempty"`
}

const schema = `
	DROP TABLE person_data_fields;
	DROP INDEX IF EXISTS person_data_fields_name;

	CREATE TABLE person_data_fields (
		id SERIAL,
		person_id integer NOT NULL, 
		key text NOT NULL,
		value text NOT NULL,
	);

	CREATE UNIQUE INDEX person_data_fields_key ON person_data_fields ('key');
    CREATE UNIQUE INDEX person_data_fields_person_id ON person_data_fields ('person_id');
`

// Map muh person map duh
type Map map[string]Model

//var Map = map[string]Model{}

var (
	mtx              sync.Mutex
	personDataFields Map
)

func init() {
	personDataFields = make(Map)
}

// Init create collection
func Init() Map {
	mtx.Lock()
	personDataFields["1"] = Model{ID: 1, PersonId: 1, Key: "email", Value: "jason@example.com"}
	personDataFields["2"] = Model{ID: 2, PersonId: 1, Key: "firstName", Value: "jason"}
	mtx.Unlock()
	return personDataFields
}

// CreateTable whatever
func CreateTable() {

	db := dbs.GetDatabaseConnection()
	db.Exec(schema)

	tx, err := db.Begin()

	if err != nil {
		panic("could not begin transaction")
	}

	s1 := personDataFields["1"]
	s2 := personDataFields["2"]

	tx.Exec("INSERT INTO person (person_id, name, value) VALUES ($1, $2, $3)", s1.PersonId, s1.Key, s1.Value)
	tx.Exec("INSERT INTO person (firstname, name, value) VALUES ($1, $2, $3)", s2.PersonId, s2.Key, s2.Value)

	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person

	// tx.NamedExec("INSERT INTO person (firstname, lastname, email) VALUES (:Firstname, :Lastname, :Email)", s1)
	// tx.NamedExec("INSERT INTO person (firstname, lastname, email) VALUES (:Firstname, :Lastname, :Email)", s2)
	tx.Commit()
}
