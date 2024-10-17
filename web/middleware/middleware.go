package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/delaneyj/datastar"
	"github.com/gorilla/sessions"
	"github.com/invopop/ctxi18n"
)

func Auth(sessionStore sessions.Store) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			session, err := sessionStore.Get(r, "auth")
			if err != nil {
				log.Printf("%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return

			}

			if session.Values["auth"] != true {
				if r.Header.Get("Datastar-Request") == "true" {
					sse := datastar.NewSSE(w, r)
					datastar.Redirect(sse, "/login")
				} else {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
				}
				return
			}

			ctx := context.WithValue(r.Context(), "user", session.Values["name"].(string))
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Lang() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// TODO: Implement language from session
			lang := "en" // Default language

			ctx, err := ctxi18n.WithLocale(r.Context(), lang)
			if err != nil {
				log.Printf("error setting locale: %v", err)
				http.Error(w, "error setting locale", http.StatusBadRequest)
				return
			}
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NoCacheHeaders() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			w.Header().Set("Vary", "*")
			w.Header().Set("Surrogate-Control", "no-store") // This header can sometimes help with reverse proxies or CDNs.
			h.ServeHTTP(w, r)
		})
	}
}
