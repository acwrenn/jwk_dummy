package server

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
)

// These formats are based on https://datatracker.ietf.org/doc/html/rfc7517
// See also - https://www.iana.org/assignments/jose/jose.xhtml

type jwkConf struct {
	pub  jwkPub
	priv jwks
}

type jwks struct {
	rsa []rsaKey
	ec  []ecdsaKey
}

type Hasher func(raw []byte) (disgest []byte)

type Key interface {
	crypto.Signer
	crypto.SignerOpts
	Hasher() Hasher
	Algo() string
	KID() string
}

type rsaKey struct {
	inner *rsa.PrivateKey
	algo  string
	kid   string

	hash   crypto.Hash
	hasher Hasher
}

func (r rsaKey) Algo() string {
	return r.algo
}

func (r rsaKey) HashFunc() crypto.Hash {
	return r.hash
}

func (r rsaKey) Hasher() Hasher {
	return r.hasher
}

func (r rsaKey) Public() crypto.PublicKey {
	return r.Public()
}

func (r rsaKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return r.inner.Sign(rand, digest, opts)
}

func (r rsaKey) KID() string {
	return r.kid
}

type ecdsaKey struct {
	inner *ecdsa.PrivateKey
	algo  string
	kid   string

	hash   crypto.Hash
	hasher Hasher
}

func (e ecdsaKey) Algo() string {
	return e.algo
}

func (e ecdsaKey) Hasher() Hasher {
	return e.hasher
}

func (e ecdsaKey) HashFunc() crypto.Hash {
	return e.hash
}

func (e ecdsaKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return e.inner.Sign(rand, digest, opts)
}

func (e ecdsaKey) Public() crypto.PublicKey {
	return e.Public()
}

func (e ecdsaKey) KID() string {
	return e.kid
}

type jwkPub struct {
	Keys []jwkKey `json:"keys"`
}

type jwkKey struct {
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`

	// EC Data
	Crv string `json:"crv,omitempty"`
	X   string `json:"x,omitempty"`
	Y   string `json:"y,omitempty"`

	// RSA Data
	Mod string `json:"n,omitempty"`
	Exp string `json:"e,omitempty"`
}

func generateJWKs(_ Config) (jwkConf, error) {
	c := jwkConf{}
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return c, fmt.Errorf("Could not generate RSA key: %w", err)
	}
	err = key.Validate()
	if err != nil {
		return c, fmt.Errorf("Stdlib generated invalid key: %w", err)
	}
	key.Precompute()

	// TODO acwrenn
	// Got myself a little turned around here.
	// Public exponents are generally known, and all the ones
	// I've ever heard of are + and < 2^64...
	// Casting int -> uint64 feels bad, but I can't think of a reason
	// why it could cause problems.
	exp := make([]byte, binary.Size(uint64(key.E)))
	binary.LittleEndian.PutUint64(exp, uint64(key.E))
	exp = trimZeros(exp)

	algo := "RS256"
	kid := "dummy-1"
	c.priv.rsa = []rsaKey{
		{
			inner: key,
			hash:  crypto.SHA256,
			hasher: func(in []byte) []byte {
				o := sha256.Sum256(in)
				return o[:]
			},
			kid:  kid,
			algo: algo,
		},
	}

	c.pub.Keys = []jwkKey{
		{
			Alg: algo,
			Use: "sig",
			Kty: "RSA",
			Kid: kid,
			Mod: base64.StdEncoding.EncodeToString(key.N.Bytes()),
			Exp: base64.RawURLEncoding.EncodeToString(exp),
		},
	}
	return c, nil
}

func serveJWKs(jwks jwkConf) http.HandlerFunc {
	p := jwks.pub
	return func(w http.ResponseWriter, req *http.Request) {
		respondWithJSON(w, http.StatusOK, p)
		return
	}
}

func trimZeros(b []byte) []byte {
	var j int
	l := len(b)
	for i := range b {
		if b[l-i-1] != 0 {
			j = l - i
			break
		}
	}
	return b[:j]
}
