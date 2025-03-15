package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/blinkinglight/chat-data-star/web/views/chatview"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go"
	datastar "github.com/starfederation/datastar/sdk/go"
)

func SetupChat(router chi.Router, session sessions.Store, ns *embeddednats.Server) error {

	type Message struct {
		Message string `json:"message"`
		Sender  string `json:"sender"`
	}

	nc, err := ns.Client()
	if err != nil {
		return err
	}

	nc.Subscribe("chat.>", func(msg *nats.Msg) {
		var message = Message{}
		err := json.Unmarshal(msg.Data, &message)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		if message.Sender != "bot" {
			// do something with user message and send back response
			nc.Publish("chat."+message.Sender, []byte(`{"message":"hello, i am bot. You send me: `+message.Message+`", "sender":"bot"}`))
		}
	})

	_ = nc
	chatIndex := func(w http.ResponseWriter, r *http.Request) {
		_, err := upsertSessionID(session, r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		chatview.Index().Render(r.Context(), w)
	}

	chatMessage := func(w http.ResponseWriter, r *http.Request) {
		id, err := upsertSessionID(session, r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var state = Message{}

		err = datastar.ReadSignals(r, &state)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		state.Sender = id
		b, err := json.Marshal(state)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		nc.Publish("chat."+id, b)
		sse := datastar.NewSSE(w, r)
		_ = sse
	}

	chatMessages := func(w http.ResponseWriter, r *http.Request) {

		id, err := upsertSessionID(session, r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var ch = make(chan *nats.Msg)
		sub, err := nc.Subscribe("chat."+id, func(msg *nats.Msg) {
			ch <- msg
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer close(ch)
		defer sub.Unsubscribe()

		sse := datastar.NewSSE(w, r)

		for {
			select {
			case <-r.Context().Done():
				return

			case msg := <-ch:
				var message = Message{}
				err := json.Unmarshal(msg.Data, &message)
				if err != nil {
					// datastar.Error(sse, err)
					return
				}
				if message.Sender == "bot" {
					sse.MergeFragmentTempl(chatview.Bot(message.Message), datastar.WithMergeAppend(), datastar.WithSelector("#chatbox"))
				} else {
					sse.MergeFragmentTempl(chatview.Me(message.Message), datastar.WithMergeAppend(), datastar.WithSelector("#chatbox"))
				}
			}
		}
	}

	router.Get("/chat", chatIndex)
	router.Post("/chat", chatMessage)
	router.Get("/messages", chatMessages)

	return nil
}

func upsertSessionID(store sessions.Store, r *http.Request, w http.ResponseWriter) (string, error) {

	sess, err := store.Get(r, "chatters")
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}
	id, ok := sess.Values["id"].(string)
	if !ok {
		id = toolbelt.NextEncodedID()
		sess.Values["id"] = id
		if err := sess.Save(r, w); err != nil {
			return "", fmt.Errorf("failed to save session: %w", err)
		}
	}
	return id, nil
}
