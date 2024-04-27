package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

func NewHashMiddleware(key string) *HashMiddleware {
	return &HashMiddleware{
		key: key,
	}
}

type HashMiddleware struct {
	key string
}

func (hm HashMiddleware) WithHash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HashSHA256") != "" && hm.key != "" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "error reading request body", http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			requestHash := r.Header.Get("HashSHA256")
			if !checkHash(bodyBytes, hm.key, requestHash) {
				http.Error(w, "хэш не совпадает", http.StatusBadRequest)
				return
			}
		}

		hw := &hashWriter{
			ResponseWriter: w,
			key:            hm.key,
			rHash:          r.Header.Get("HashSHA256"),
		}

		h.ServeHTTP(hw, r)
	})
}

type hashWriter struct {
	http.ResponseWriter
	key   string
	rHash string
}

func (hw hashWriter) Write(b []byte) (int, error) {
	if hw.rHash != "" && hw.key != "" {
		hw.Header().Set("HashSHA256", string(getHash(b, hw.key)))
	}
	return hw.ResponseWriter.Write(b)
}

func checkHash(data []byte, key string, requestHash string) bool {
	return hmac.Equal([]byte(getHash(data, key)), []byte(requestHash))
}

func getHash(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)

	return hex.EncodeToString(h.Sum(nil))
}
