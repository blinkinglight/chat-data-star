package handlers

import (
	"net/http"

	"github.com/blinkinglight/chat-data-star/web/views/layoutview"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupHome(router chi.Router, session sessions.Store, ns *embeddednats.Server) error {

	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		layoutview.Index().Render(r.Context(), w)
	}

	router.Get("/", homeHandler)

	return nil
}
