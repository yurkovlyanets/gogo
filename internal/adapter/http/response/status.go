package response

import "net/http"

func WriteOKJSON(w http.ResponseWriter, v any) {
	writeJSON(w, http.StatusOK, v)
}

func WriteCreated(w http.ResponseWriter, v any) {
	writeJSON(w, http.StatusCreated, v)
}

func WriteBadRequest(w http.ResponseWriter, msg string) {
	writeError(w, http.StatusBadRequest, msg)
}

func WriteNotFound(w http.ResponseWriter, msg string) {
	writeError(w, http.StatusNotFound, msg)
}

func WriteInternal(w http.ResponseWriter, msg string) {
	writeError(w, http.StatusInternalServerError, msg)
}
