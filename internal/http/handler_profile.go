package http

import (
	"net/http"
)

type ProfileHandler struct {
	template *Template
}

func NewProfileHandler(
	template *Template,
) *ProfileHandler {
	return &ProfileHandler{
		template: template,
	}
}

func (h *ProfileHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := h.template.Render(w, "profile", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
