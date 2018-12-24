package person_data_fields

import (
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

