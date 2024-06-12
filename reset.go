package main

import "net/http"

func (cfg *apiConfig) handlerResetServerHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server hits has been reset to 0."))
}
