package share

import (
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


