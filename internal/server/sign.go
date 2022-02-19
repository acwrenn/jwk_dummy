package server

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type opts struct {
	hasher crypto.Hash
}

type tokenOutput struct {
	Token string `json:"token"`
}

func signPayload(jwks jwkConf) http.HandlerFunc {
	// TODO acwrenn
	// Come up with a more sophisticated method for selecting keys.
	// We could pretty simply have the user request it via a QP.
	priv := jwks.priv
	var signer Key
	if len(priv.rsa) > 0 {
		signer = priv.rsa[0]
	} else if len(priv.ec) > 0 {
		signer = priv.ec[0]
	}
	return func(w http.ResponseWriter, req *http.Request) {
		raw, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			respondWithString(w, http.StatusNotAcceptable, `{"msg": "Not acceptable - could not read stream"}`)
			return
		}
		header := map[string]interface{}{
			"alg": signer.Algo(),
			"typ": "JWT",
			"kid": signer.KID(),
		}
		msg, err := assembleMessage(header, raw)
		if err != nil {
			fmt.Println(err)
			respondWithString(w, http.StatusInternalServerError, `{"msg": "Internal Server Error - could not assemble JWT"}`)
			return
		}

		sig, err := signer.Sign(rand.Reader, signer.Hasher()([]byte(msg)), signer)
		if err != nil {
			fmt.Println(err)
			respondWithString(w, http.StatusInternalServerError, `{"msg": "Internal Server Error- could not sign payload."}`)
			return
		}
		jwt := assembleJWT(msg, sig)
		respondWithJSON(w, http.StatusOK, tokenOutput{Token: jwt})
		return
	}
}

func assembleMessage(header map[string]interface{}, body []byte) (string, error) {
	h, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", base64.RawURLEncoding.EncodeToString(h), base64.RawURLEncoding.EncodeToString(body)), nil
}

func assembleJWT(body string, sig []byte) string {
	return fmt.Sprintf("%s.%s", body, base64.RawURLEncoding.EncodeToString(sig))
}
