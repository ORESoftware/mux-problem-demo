package nearby

import (
	"reflect"
	"sync"
	"time"
)

// Model - Nearby model whatever
type Model struct {
	ID          int   `json:"id,omitempty"`
	Me          int   `json:"me,omitempty"`
	You         int   `json:"you,omitempty"`
	ContactTime int64 `json:"contactTime,omitempty"`
}

var schema = `
	DROP TABLE nearby;

	CREATE TABLE nearby (
		id SERIAL,
		me integer,
		you integer,
		contactTime bigint
	) PARTITION BY LIST(me);

	CREATE TABLE nearby_0 PARTITION OF nearby FOR VALUES IN (0);
	CREATE TABLE nearby_1 PARTITION OF nearby FOR VALUES IN (1);
	CREATE TABLE nearby_2 PARTITION OF nearby FOR VALUES IN (2);
	CREATE TABLE nearby_3 PARTITION OF nearby FOR VALUES IN (3);
`

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(1)
}

func getValues(m interface{}) []interface{} {
	v := reflect.ValueOf(m)
	result := make([]interface{}, 0, v.Len())
	for _, k := range v.MapKeys() {
		result = append(result, v.MapIndex(k).Interface())
	}
	return result
}

// Map - id to model map
type Map map[string]Model

var (
	mtx    sync.Mutex
	nearby Map
)

func init(){
	nearby = make(map[string]Model)
}

// Init create collection
func Init() Map {
	mtx.Lock()
	nearby["1"] = Model{ID: 1, Me: 1, You: 2, ContactTime: makeTimestamp()}
	nearby["2"] = Model{ID: 2, Me: 2, You: 1, ContactTime: makeTimestamp()}
	mtx.Unlock()
	return nearby
}

