package routes

import (
	"huru/routes/data-fields"
	"huru/routes/login"
	"huru/routes/nearby"
	"huru/routes/person"
	"huru/routes/register"
	"huru/routes/share"
)

// Handlers
type RegisterHandler = register.Handler
type LoginHandler = login.Handler
type NearbyHandler = nearby.Handler
type ShareHandler = share.Handler
type PersonHandler = person.Handler
type PersonDataHandler = data_fields.Handler


type HandlerCreator = func() struct{}

var Handlers = map[string]HandlerCreator{
	"Register": register.CreateHandler,
}

type NearbyInjection = nearby.Injection
type ShareInjection = share.Injection
type PersonInjection = person.Injection
type PersonDataInjection = data_fields.Injection
type RegisterInjection = register.Injection
type LoginInjection = login.Injection


type HuruRouteHandlers struct {
	Register register.Handler
	Login    login.Handler
	Nearby   nearby.Handler
	Share    share.Handler
	Person   person.Handler
}

func (h HuruRouteHandlers) GetHandlers() HuruRouteHandlers {
	return HuruRouteHandlers{
		register.Handler{},
		login.Handler{},
		nearby.Handler{},
		share.Handler{},
		person.Handler{},
	}
}
