package handler

import (
	"net/http"
)

func (h *Handler) createCharacter(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) getAllCharacters(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) getCharacter(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) updateCharacter(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) deleteCharacter(w http.ResponseWriter, r *http.Request) {}
