package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PostBody struct {
	Choice string `json:"choice"`
}

type ResponseBody struct {
	Plot    string   `json:"plot"`
	Choices []string `json:"choices"`
}

func (s *Server) addNarratorRoutes(r chi.Router) {
	r.Route("/narrate/{sessionKey}", func(r chi.Router) {
		// Begin story
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			sessionKey := chi.URLParam(r, "sessionKey")
			s.logger.Info("Creating story", "sessionKey", sessionKey)
			chatResponse, err := s.Narrator.CreateStory(r.Context(), sessionKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Trim response data
			response := ResponseBody{
				Plot:    chatResponse.Plot,
				Choices: chatResponse.Choices,
			}

			resp, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Return JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(resp)
		})

		// Update story
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			sessionKey := chi.URLParam(r, "sessionKey")
			var choice PostBody
			if err := json.NewDecoder(r.Body).Decode(&choice); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response, err := s.Narrator.UpdateStory(r.Context(), sessionKey, choice.Choice)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			resp, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Return JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(resp)
		})
	})
}
