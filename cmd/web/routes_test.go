package main

import (
	"fmt"
	"github.com/SnehilSundriyal/bookings-go/internal/config"
	"github.com/go-chi/chi"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, but is %T", v))
	}
}
