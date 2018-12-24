package role

import (
	"sync"
)

// Model The person Type (more like an object)
type Model struct {
	ID        int    `json:"id, omitempty"`
	Name    string `json:"name, omitempty"`
}

const schema = `
	DROP TABLE role;
	DROP INDEX IF EXISTS role_name;

	CREATE TABLE role (
		id SERIAL,
		name text,
	);

	CREATE UNIQUE INDEX role_name ON role (name);
`

// Map muh person map duh
type Map map[string]Model

//var Map = map[string]Model{}

var (
	mtx   sync.Mutex
	roles Map
)

func init() {
	roles = make(Map)
}

// Init create collection
func Init() Map {
	mtx.Lock()
	roles["1"] = Model{ID: 1, Name: "admin"}
	roles["2"] = Model{ID: 2, Name: "user"}
	mtx.Unlock()
	return roles
}

// CreateTable whatever
func CreateTable() {

}
