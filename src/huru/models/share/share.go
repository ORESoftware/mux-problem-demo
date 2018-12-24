package share

import (
	"huru/dbs"
	"sync"
)

// Model The share type
type Model struct {
	ID         int    `json:"id"`
	Me         int    `json:"me"`
	You        int    `json:"you"`
	FieldName  string `json:"fieldName"`
	FieldValue bool   `json:"fieldValue"`
}

const schema = `
	DROP TABLE share;

	CREATE TABLE share (
		id SERIAL,
		me int,
		you int,
		fieldName text,
		fieldValue boolean
	) PARTITION BY LIST(me);

	CREATE TABLE share_0 PARTITION OF share FOR VALUES IN (0);
	CREATE TABLE share_1 PARTITION OF share FOR VALUES IN (1);
	CREATE TABLE share_2 PARTITION OF share FOR VALUES IN (2);
`

// Map mappy
type Map map[string]Model

var (
	mtx    sync.Mutex
	shares Map
)

func init(){
	shares = make(Map)
}

// Init create collection
func Init() Map {
	mtx.Lock()
	shares["1"] = Model{ID: 1, Me: 1, You: 2, FieldName: "sharePhone", FieldValue: false}
	shares["2"] = Model{ID: 2, Me: 2, You: 1, FieldName: "shareEmail", FieldValue: true}
	shares["3"] = Model{ID: 3, Me: 1, You: 2, FieldName: "sharePhone", FieldValue: true}
	shares["4"] = Model{ID: 4, Me: 2, You: 1, FieldName: "shareEmail", FieldValue: false}
	mtx.Unlock()
	return shares
}

// CreateTable whatever
func CreateTable() {

	s1 := shares["1"]
	s2 := shares["2"]
	s3 := shares["3"]

	db := dbs.GetDatabaseConnection()
	db.Exec(schema)

	tx, err := db.Begin()
	if err != nil {
		panic("could not begin transaction")
	}

	tx.Exec("INSERT INTO share (me, you, fieldName, fieldValue) VALUES ($1, $2, $3, $4)", s1.Me, s1.You, s1.FieldName, s1.FieldValue)
	tx.Exec("INSERT INTO share (me, you, fieldName, fieldValue) VALUES ($1, $2, $3, $4)", s2.Me, s2.You, s2.FieldName, s2.FieldValue)
	tx.Exec("INSERT INTO share (me, you, fieldName, fieldValue) VALUES ($1, $2, $3, $4)", s3.Me, s3.You, s3.FieldName, s3.FieldValue)

	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	// tx.NamedExec("INSERT INTO share (me, you, sharePhone, shareEmail) VALUES (:me, :you, :sharePhone, :shareEmail)", s1)
	// tx.NamedExec("INSERT INTO share (me, you, sharePhone, shareEmail) VALUES (:me, :you, :sharePhone, :shareEmail)", s2)
	tx.Commit()

}
