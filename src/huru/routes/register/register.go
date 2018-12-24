package register

import (
	"encoding/json"
	"fmt"
	"huru/dbs"
	"huru/models/person"
	"huru/mw"
	tc "huru/type-creator"
	"huru/utils"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func getTableCreationCommands(v int) string {
	return fmt.Sprintf(`
		CREATE TABLE share_%d PARTITION OF share FOR VALUES IN (%d);
		CREATE TABLE nearby_%d PARTITION OF nearby FOR VALUES IN (%d);
	`, v, v, v, v)
}

var (
	mtx sync.Mutex
)

// Handler => RegisterHandler just what it says
type Handler struct{}

// CreateHandler whatever
func CreateHandler() struct{} {
	return Handler{}
}

type Injection struct {
	PeopleById, PeopleByHandle, PeopleByEmail person.Map
}

func (i Injection) GetVal() interface{} {
	return i
}

// Mount just what it says
func (h Handler) Mount(router *mux.Router, v Injection) Handler {
	mwList := mw.List(mw.Error())
	router.HandleFunc("/register", mw.Middleware(mwList,h.makeRegisterNewUser(v))).Methods("POST")
	return h
}


func (h Handler) makeRegisterNewUser(v Injection) http.HandlerFunc {

	type RespBody struct {
		TypeCreatorMeta string `type:"bar" tc_resp_body_type:"1"`
	}

	type ReqBody struct {
		TypeCreatorMeta string `type:"star" tc_req_body_type:"1"`
		Handle string
	}

	return tc.ExtractType(
		tc.TypeList{ReqBody{}, RespBody{}},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			//decoder := json.NewDecoder(r.Body)
			//var t ReqBody
			//err := decoder.Decode(&t)

			var user person.Model;
			json.NewDecoder(r.Body).Decode(&user)


			if user.Email == "" {
				panic(utils.AppError{
					StatusCode: 409,
					Message:    "Missing handle/email property in request body.",
				})
			}

			//if _, ok := v.PeopleByHandle[user.Handle]; ok {
			//	panic("already have a person with handle " + user.Handle)
			//}

			if _, ok := v.PeopleByEmail[user.Email]; ok {
				 panic("already have a person with email " + user.Email)
                 return;
			}


			//v.PeopleByHandle[user.Handle] = user
			v.PeopleByEmail[user.Email] = user
			//v.PeopleById[user.ID] = user;  // need to get ID from database

			json.NewEncoder(w).Encode(user);
			return




			db := dbs.GetDatabaseConnection()

			tx, err := db.Begin()

			if err != nil {
				panic(err)
			}

			p := person.Model{Handle: user.Handle}

			var id int
			err = tx.QueryRow("INSERT INTO person (handle,email,firstname,lastname) VALUES ($1, $2, $3, $4) RETURNING ID", p.Handle, p.Email, p.Firstname, p.Lastname).Scan(&id)

			if err != nil {
				panic(err)
			}

			db.Exec(getTableCreationCommands(id))
			err = tx.Commit()

			if err != nil {
				panic(err)
			}

			//json.NewEncoder(w).Encode(ResolutionValue{"foo"})

			//json.NewEncoder(w).Encode(struct {
			//	ID string
			//}{
			//	fmt.Sprintf("%d", id),
			//})

			json.NewEncoder(w).Encode(RespBody{});
		}))
}
