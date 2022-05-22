package shortener

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Asymmetriq/shortener/internal/cookie"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") &&
			!strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzWriter, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gzWriter.Close()

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") &&
			!strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzWriter}, r)
			return
		}

		gzReader, err := gzip.NewReader(r.Body)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		var buf bytes.Buffer
		_, err = buf.ReadFrom(gzReader)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		r.ContentLength = int64(len(buf.Bytes()))
		r.Body = ioutil.NopCloser(&buf)

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzWriter}, r)
	})
}

type cookieWriter struct {
	http.ResponseWriter
	authTicket string
}

func (rw cookieWriter) WriteHeader(statusCode int) {
	http.SetCookie(rw.ResponseWriter, &http.Cookie{
		Name:  string(cookie.Name),
		Value: rw.authTicket,
	})
	rw.ResponseWriter.WriteHeader(statusCode)
}

func cookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDCookie, err := r.Cookie(string(cookie.Name))
		var id, authTicket string
		switch err {
		case nil:
			id, authTicket = cookie.CheckUserID(userIDCookie.Value)
		case http.ErrNoCookie:
			id, authTicket = cookie.GetSignedUserID()
		default:
			http.Error(w, "cookie parsing", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(cookieWriter{ResponseWriter: w, authTicket: authTicket}, r.WithContext(context.WithValue(r.Context(), cookie.Name, id)))
	})
}
