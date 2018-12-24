package type_creator

import (
	"fmt"
	"huru/utils"
	"net/http"
	"reflect"
)


type Docs struct {
	Methods []string
	Route string
}


type TypeList = []interface{}

func ExtractType(s TypeList, h http.HandlerFunc) http.HandlerFunc {

	m := make(map[string]string)

	for _, v := range s {
		t := reflect.TypeOf(v)
		f, _ := t.FieldByName("TypeCreatorMeta")
		//fmt.Println("All tags?:",f.Tag);
		v, ok := f.Tag.Lookup("type")
		if !ok {
			fmt.Println("no 'type' tag.");
			continue;
		}
		for _, key := range []string{"tc_req_body_type", "tc_resp_body_type"} {
			_, ok := f.Tag.Lookup(key)
			//fmt.Println(ok,"key:",key)
			if ok {
				m[key] = v
			}
		}

		//fmt.Println(v, ok)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("tc_req_body_type") != m["tc_req_body_type"] {
			fmt.Println(
				utils.JoinArgs(
					"Request body types are different,",
					"actual:", r.Header.Get("tc_req_body_type"),
					"expected:", m["tc_req_body_type"],
				),
			)
		}

		if r.Header.Get("tc_resp_body_type") != m["tc_resp_body_type"] {
			fmt.Println(
				utils.JoinArgs(
					"Response body types are different,",
					"actual:", r.Header.Get("tc_resp_body_type"),
					"expected:", m["tc_resp_body_type"],
				),
			)
		}

		w.Header().Set("zoom", "bar")


		for _, key := range []string{"tc_resp_body_type"} {
			v, ok := m[key]
			if ok {
				w.Header().Set(key, v)
			}
		}


		fmt.Printf("Req: %s\n", r.URL.Path)
		h.ServeHTTP(w, r)
	}
}
