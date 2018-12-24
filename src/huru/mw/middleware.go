package mw

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"huru/models/person"
	"huru/utils"
	"net/http"
	"os"
)

func Deny(roles ...string) Adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user, ok := context.Get(r, "logged_in_user").(person.Model);
			if ok == false {
				panic("'logged_in_user' is not in the context map.")
			}

			if user.Role == "" {
				panic("logged in user does not have a role.");
			}

			for _, rl := range roles {
				if rl == user.Role {
					panic("user role was not allowed: " + user.Role)
				}
			}

			next(w, r)
			//next.ServeHTTP(w, r)
		})
	}
}

func Allow(roles ...string) Adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user, ok := context.Get(r, "logged_in_user").(person.Model);
			if ok == false {
				panic("'logged_in_user' is not in the context map.")
			}

			if user.Role == "" {
				panic("logged in user does not have a role.");
			}

			for _, rl := range roles {
				if rl == user.Role {
					next(w, r)
					//next.ServeHTTP(w, r)
					return
				}
			}

			panic("user role was not allowed: " + user.Role)
		})
	}
}

func Logging(v interface{}) Adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Do stuff here
			log.Println("Here is the request URI:", r.RequestURI)
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next(w, r)
		})
	}
}

func LoggingX(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("Here is the request URI:", r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

//func LoggingHandler() mux.MiddlewareFunc {
//	logger := Logging(struct{}{});
//	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
//		log.Println("logging middleware registered");
//		return logger(next.ServeHTTP);
//	})
//}

func LoggingHandler() mux.MiddlewareFunc {
	logger := Logging(struct{}{});
	return func(next http.Handler) http.Handler {
		log.Println("logging middleware registered");
		return logger(next.ServeHTTP);
	}
}

type Exception struct {
	Message string `json:"message"`
}

func Auth(huruApiToken string) Adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {

		huruEnv := os.Getenv("huru_env");
		isLocal := huruEnv == "local"

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			params := r.URL.Query()
			fmt.Println("the params are:", params);
			tok := params.Get(huruApiToken)

			if tok == "" {
				tok = r.Header.Get(huruApiToken)
			}

			if tok == "" {
				cookie, err := r.Cookie(huruApiToken);
				if err != nil {
					log.Error(err)
				}
				if cookie != nil {
					tok = cookie.Value
				}
			}

			if tok == "" {

				if isLocal != true {
					panic("No token in request and huru_env is not local.");
				}

				tok = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
					"eyJoYW5kbGUiOiIiLCJwYXNzd29yZCI6IiIsInVzZXJuYW1lIjoiIn0." +
					"4-pak7FLp_83Px7N0NZHPt4mpC5LzqDgkWsnk4JaHTg"

			}

			token, _ := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error")
				}
				return []byte("secret"), nil
			})

			claims, ok := token.Claims.(jwt.MapClaims)

			if ! (ok && token.Valid) {
				json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				return;
			}

			var user person.Model
			mapstructure.Decode(claims, &user)

			fmt.Printf("%+v\n", user)

			//log.Info("the logged-in user is:", user)
			context.Set(r, "logged_in_user", user)

			log.Warn("About to call next in AUTH.");
			//next.ServeHTTP(w, r)
			next(w,r);
		})
	}
}

type Z = func(next http.HandlerFunc) http.Handler;

func AuthHandler() mux.MiddlewareFunc {  // used to return Z, now returns mux.MiddlewareFunc
	huruApiToken := "x-huru-api-token"
	auth := Auth(huruApiToken)
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		log.Println("auth middleware registered.")
		return auth(next.ServeHTTP)
	})
}

func Error() Adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {

					log.Error("Caught error in defer/recover middleware: ", err)

					message := ""
					statusCode := 0

					if v, ok := err.(struct{ StatusCode int }); ok {
						statusCode = v.StatusCode
					}

					if statusCode != 0 {
						w.WriteHeader(statusCode)
					} else {
						w.WriteHeader(http.StatusInternalServerError)
					}

					if message, ok := err.(string); ok {
						json.NewEncoder(w).Encode(struct {
							Message string
						}{
							message,
						})
						return;
					}

					if v, ok := err.(struct{ OriginalError error }); ok {
						log.Error("Caught error in defer/recover middleware: ", err)
						originalError := v.OriginalError

						if originalError != nil {
							log.Error("Original error in defer/recover middleware: ", originalError)
						}
					}

					if v, ok := err.(struct{ Message string }); ok {
						message = v.Message
					}

					if message == "" {
						message = "Unknown error message."
					}

					json.NewEncoder(w).Encode(struct {
						Message string
					}{
						message,
					})

				}
			}()

			log.Info("We are handling errors with defer.")
			//next.ServeHTTP(w, r)
			next(w,r);
		})
	}
}

//
//func ErrorHandler() mux.MiddlewareFunc {
//	handler := Error();
//	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
//		return handler(next.ServeHTTP);
//	})
//}


func ErrorHandler() mux.MiddlewareFunc {
	handler := Error();
	return func(next http.Handler) http.Handler {
		return handler(next.ServeHTTP);
	}
}


type Adapter = func(http.HandlerFunc) http.HandlerFunc

func Adapt(h http.Handler, adapters ...mux.MiddlewareFunc) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func AdaptFuncs(h http.HandlerFunc, adapters ...Adapter) http.HandlerFunc {

	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func List(adapters ...Adapter) []Adapter {
	return adapters
}

var (
	counter = 1
)

func Middleware(adapters ...interface{}) http.HandlerFunc {

	counter = counter + 1;

	adapters = utils.FlattenDeep(adapters)

	if len(adapters) < 1 {
		panic("Adapters need to have length > 0.");
	}

	last := adapters[len(adapters)-1]

	h, ok := last.(http.HandlerFunc)

	if ok == false {
		log.Warn("Was super false");
		h = http.HandlerFunc(h);
	}

	adapters = adapters[:len(adapters)-1]

	for i := len(adapters)-1; i >= 0; i-- {

		adapt := adapters[i]

		adapt, _ = adapt.(Adapter);

		h = (adapt.(Adapter))(h)

	}

	return h

}

