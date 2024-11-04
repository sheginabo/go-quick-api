package middlewares

import (
	"bytes"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"time"
)

type bodyLogWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		var reqBody []byte
		var err error
		if r.Body != nil {
			reqBody, err = io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: w}
		w = blw

		next.ServeHTTP(w, r)

		elapsed := time.Since(startTime)

		respBody := blw.body.String()

		requestHeaderMap := map[string]string{}

		requestMap := map[string]interface{}{
			"url":     r.URL.Path,
			"method":  r.Method,
			"body":    string(reqBody),
			"headers": requestHeaderMap,
		}
		if q := r.URL.RawQuery; q != "" {
			requestMap["query_string"] = q
		}

		responseMap := map[string]interface{}{
			"status_code": blw.ResponseWriter.WriteHeader,
			"body":        respBody,
		}

		log.Info().
			Int64("duration_ms", elapsed.Milliseconds()).
			Interface("request", requestMap).
			Interface("response", responseMap).
			Msgf("[%s] API Request and Response", viper.GetString("APP_NAME"))
	})
}
