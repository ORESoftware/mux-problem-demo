package routes

import (
	"huru/routes/data-fields"
	"huru/routes/nearby"
	"huru/routes/share"
)

// Handlers

type NearbyHandler = nearby.Handler
type ShareHandler = share.Handler
type PersonDataHandler = data_fields.Handler


type HandlerCreator = func() struct{}



type NearbyInjection = nearby.Injection
type ShareInjection = share.Injection
type PersonDataInjection = data_fields.Injection


type HuruRouteHandlers struct {
	Nearby   nearby.Handler
	Share    share.Handler
}

