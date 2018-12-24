package person

import (
	"huru/dbs"
	"sync"
)

// Model The person Type (more like an object)
type Model struct {
	ID        int    `json:"id, omitempty"`
	Handle    string `json:"handle, omitempty"`
	Role      string `json:"role, omitempty"`
	Password  string `json:"password, omitempty"`
	Work      string `json:"work, omitempty"`
	Image     string `json:"image, omitempty"`
	Firstname string `json:"firstname, omitempty"`
	Lastname  string `json:"lastname, omitempty"`
	Phone     string `json:"phone, omitempty"`
	Email     string `json:"email, omitempty"`
}

const schema = `
	DROP TABLE person;
	DROP INDEX IF EXISTS person_handle;
	DROP INDEX IF EXISTS person_email;

	CREATE TABLE person (
		id SERIAL,
		handle text,
		firstname text,
		lastname text,
		email text,
		phone text,
		role text,
		work text,
		image text
	);

	CREATE UNIQUE INDEX person_handle ON person (handle);
	CREATE UNIQUE INDEX person_email ON person (email);
`

// Map muh person map duh
type Map map[string]Model

//var Map = map[string]Model{}

var (
	mtx            sync.Mutex
	peopleById     Map
	peopleByHandle Map
	peopleByEmail Map
)

func init() {
	peopleById = make(Map)
	peopleByHandle = make(Map)
	peopleByEmail = make(Map)
}

// Init create collection
func Init() (Map,Map, Map) {
	mtx.Lock()
	peopleById["1"] = Model{ID: 1, Handle: "alex", Firstname: "Alex", Lastname: "Chaz", Email: "alex@example.com", Password:"foo"}
	peopleById["2"] = Model{ID: 2, Handle: "jason",Firstname: "Jason", Lastname: "Statham", Email: "jason@example.com", Password:"foo"}
	peopleByHandle["alex"] = peopleById["1"]
	peopleByHandle["jason"] = peopleById["2"]

	mtx.Unlock()
	return peopleById, peopleByHandle, peopleByEmail
}

// CreateTable whatever
func CreateTable() {

	db := dbs.GetDatabaseConnection()
	db.Exec(schema)

	tx, err := db.Begin()

	if err != nil {
		panic("could not begin transaction")
	}

	s1 := peopleById["1"]
	s2 := peopleById["2"]

	tx.Exec("INSERT INTO person (firstname, lastname, email) VALUES ($1, $2, $3)", s1.Firstname, s1.Lastname, s1.Email)
	tx.Exec("INSERT INTO person (firstname, lastname, email) VALUES ($1, $2, $3)", s2.Firstname, s2.Lastname, s2.Email)

	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person

	// tx.NamedExec("INSERT INTO person (firstname, lastname, email) VALUES (:Firstname, :Lastname, :Email)", s1)
	// tx.NamedExec("INSERT INTO person (firstname, lastname, email) VALUES (:Firstname, :Lastname, :Email)", s2)
	tx.Commit()
}
