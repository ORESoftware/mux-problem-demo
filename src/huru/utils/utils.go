package utils

import (
	"bytes"
	"os"
	"path"
	"reflect"
)

// APIDoc is docs
type APIDoc struct {
	Route           string
	ResolutionValue interface{}
}

var docsFileName = path.Join(os.Getenv("HOME"), ".huru", "docs.json")

func getKind(v interface{}) string {

	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Slice:
		return "slice"
	case reflect.Array:
		return "array"
	default:
		return "unknown"
	}

}

func getKind2(v interface{}) interface{} {
	return reflect.TypeOf(v).Kind()
}

func FlattenDeepInternal(args []interface{}, v reflect.Value) []interface{} {

	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			args = FlattenDeepInternal(args, v.Index(i))
		}
	} else {
		args = append(args, v.Interface())
	}

	return args
}

func FlattenDeep(args ...interface{}) []interface{} {
	return FlattenDeepInternal(nil,reflect.ValueOf(args));
}

func FlattenDeep2(args ...interface{}) []interface{} {

	list := []interface{}{}
	//var list []interface{} = make([]interface{})

	for _, v := range args {

		kind := getKind(v)

		if kind != "unknown" {

			a, _ := v.([]interface{})

			for _, z := range FlattenDeep2(a...) {
				list = append(list, z)
			}

		} else {
			list = append(list, v);
		}
	}

	return list;
}

//func WriteToDocs(v interface{}) {
//
//	z := typings.Entities{}
//
//	log.Fatal(z)
//
//	f, err := os.OpenFile(docsFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
//
//	defer f.Close()
//
//	if err != nil {
//		panic(err)
//	}
//
//	buf := bytes.NewBuffer(nil)
//	json.NewEncoder(buf).Encode(v)
//
//	if _, err = f.WriteString(buf.String()); err != nil {
//		panic(err)
//	}
//}

// JoinArgs joins strings
func JoinArgs(strangs ...string) string {
	buffer := bytes.NewBufferString("")
	for _, s := range strangs {
		buffer.WriteString(s + " ")
	}
	return buffer.String()
}

// AppError send this response
type AppError struct {
	StatusCode    int
	Message       string
	OriginalError error
}

// SetFields on obj
func SetFields(dst, src interface{}, names ...string) {
	d := reflect.ValueOf(dst).Elem()
	s := reflect.ValueOf(src).Elem()
	for _, name := range names {
		df := d.FieldByName(name)
		sf := s.FieldByName(name)
		switch sf.Kind() {
		case reflect.String:
			if v := sf.String(); v != "" {
				df.SetString(v)
			}

		case reflect.Bool:
			if v := sf.Bool(); v != false {
				df.SetBool(v)
			}
		}

	}
}

// SetExistingFields only set fields on dst that exist in the names array
func SetExistingFields1(src interface{}, dst interface{}, names ...string) {

	fields := reflect.TypeOf(src)
	// values := reflect.ValueOf(src)

	num := fields.NumField()
	s := reflect.ValueOf(src).Elem()
	d := reflect.ValueOf(dst).Elem()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		// value := values.Field(i)
		fsrc := s.FieldByName(field.Name)
		fdest := d.FieldByName(field.Name)

		switch fsrc.Kind() {
		case reflect.String:
			if v := fsrc.String(); v != "" {
				fdest.SetString(v)
			}

		case reflect.Bool:
			if v := fsrc.Bool(); v != false {
				fdest.SetBool(v)
			}
		}

		// if fdest.IsValid() && fsrc.IsValid() {
		// 	// A Value can be changed only if it is
		// 	// addressable and was not obtained by
		// 	// the use of unexported struct fields.
		// 	if fdest.CanSet() && fsrc.CanSet() {
		// 		// change value of N

		// 		fdest.Set(value)

		// 	}
		// } else {
		// 	log.Fatal("not valid", fdest)
		// }

	}
}

func contains(strangs []string, e string) bool {
	for _, a := range strangs {
		if a == e {
			return true
		}
	}
	return false
}

func SetExistingFields(src interface{}, dst interface{}, limit bool, names ...string) {

	srcFields := reflect.TypeOf(src).Elem()
	srcValues := reflect.ValueOf(src).Elem()
	dstValues := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcFields.NumField(); i++ {
		sf := srcFields.Field(i)

		if limit == true && contains(names, sf.Name) == false {
			continue
		}

		sv := srcValues.Field(i)
		dv := dstValues.FieldByName(sf.Name)

		if dv.IsValid() && dv.CanSet() {
			dv.Set(sv)
		}

	}
}

// SetExistingFields only set fields on dst that exist in the names array
func SetExistingFields2(src, dst interface{}, limit bool, names ...string) {

	fields := reflect.TypeOf(src)
	values := reflect.ValueOf(src)

	s := reflect.ValueOf(src).Elem()
	d := reflect.ValueOf(dst).Elem()

	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)

		if limit == true && contains(names, field.Name) == false {
			continue
		}

		value := values.Field(i)

		fsrc := s.FieldByName(field.Name)
		fdest := d.FieldByName(field.Name)

		if fdest.IsValid() && fsrc.IsValid() {
			// A Value can be changed only if it is
			// addressable and was not obtained by
			// the use of unexported struct fields.
			if fdest.CanSet() && fsrc.CanSet() {
				// change value of N

				fdest.Set(value)

			}
		}

	}
}
