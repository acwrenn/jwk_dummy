package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Run(address string, config Config) error {
	jwks, err := generateJWKs(config)
	if err != nil {
		return err
	}

	http.HandleFunc(config.ConfigRoute, onlyMethod(http.MethodGet, serveConfig(config)))
	http.HandleFunc("/.well-known/jwks.json", onlyMethod(http.MethodGet, serveJWKs(jwks)))
	http.HandleFunc("/sign", onlyMethod(http.MethodPost, signPayload(jwks)))
	return http.ListenAndServe(address, nil)
}

func onlyMethod(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			respondWithJSON(w, http.StatusMethodNotAllowed, J{"msg": "Method not allowed"})
			return
		}
		handler(w, req)
	}
}

type openIDConfig struct {
	JWKUri string `json:"jwks_uri"`
}

func serveConfig(conf Config) http.HandlerFunc {
	c := openIDConfig{
		JWKUri: fmt.Sprintf("%s://%s:%d/.well-known/jwks.json", conf.Protocol, conf.Address, conf.Port),
	}
	return func(w http.ResponseWriter, req *http.Request) {
		respondWithJSON(w, http.StatusInternalServerError, c)
	}

}

type J = map[string]string

func respondWithJSON(w http.ResponseWriter, code int, j interface{}) {
	w.Header().Add("Content-Type", "application/json")
	raw, err := json.Marshal(j)
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respondWithBytes(w, code, raw)
}

func respondWithString(w http.ResponseWriter, code int, s string) {
	respondWithBytes(w, code, []byte(s))
}

func respondWithBytes(w http.ResponseWriter, code int, b []byte) {
	respondWithReader(w, code, bytes.NewBuffer(b))
}

func respondWithReader(w http.ResponseWriter, code int, r io.Reader) {
	w.WriteHeader(code)
	_, err := io.Copy(w, r)
	if err != nil {
		fmt.Println(err)
	}
}
