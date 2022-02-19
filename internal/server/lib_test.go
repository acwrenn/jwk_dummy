package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
)

func TestSignature(t *testing.T) {
	port := 10000 + rand.Intn(20000)
	conf := Config{
		Protocol: "http",
		Address:  "localhost",
		Port:     port,

		ConfigRoute: "/.well-known/openid-configuration",
	}
	go func() {
		Run(fmt.Sprintf("localhost:%d", port), conf)
	}()
	ctx, cls := context.WithTimeout(context.Background(), time.Minute)
	defer cls()
	for {
		select {
		case <-ctx.Done():
			break
		default:
		}
		_, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	iss := fmt.Sprintf("http://localhost:%d", port)
	aud := "TestAPI"

	tcv := map[string]string{
		"aud": aud,
	}
	v := jwtverifier.JwtVerifier{
		Issuer:           iss,
		ClaimsToValidate: tcv,
	}
	jv := v.New()

	m := "Hello!"
	claims := map[string]interface{}{
		"msg": m,
		"aud": aud,
		"iss": iss,
		"exp": time.Now().Add(100 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	raw, err := json.Marshal(claims)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/sign", iss), "application/json", bytes.NewBuffer(raw))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	to := tokenOutput{}
	err = json.Unmarshal(out, &to)
	if err != nil {
		t.Fatal(err)
	}

	token, err := jv.VerifyAccessToken(to.Token)
	if err != nil {
		t.Fatal(err)
	}
	if token.Claims["msg"] != m {
		t.Fatalf("Incorrect token returned: %s vs %+v", m, token.Claims["msg"])
	}
	t.Logf("Token validation success: %+v\n", token)
}
