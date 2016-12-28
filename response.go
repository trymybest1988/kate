package kate

import (
	"bufio"
	"encoding/json"
	"net"
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter

	StatusCode() int

	Data() interface{}
	SetData(interface{})

	RawBody() []byte

	EncodeJson(v interface{}) ([]byte, error)

	WriteJson(v interface{}) error
}

type responseWriter struct {
	http.ResponseWriter

	wroteHeader bool
	statusCode  int
	rawBody     []byte
	data        interface{}
}

func (w *responseWriter) StatusCode() int {
	return w.statusCode
}

func (w *responseWriter) Data() interface{} {
	return w.data
}

func (w *responseWriter) SetData(v interface{}) {
	w.data = v
}

func (w *responseWriter) RawBody() []byte {
	return w.rawBody
}

func (w *responseWriter) WriteHeader(code int) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.ResponseWriter.WriteHeader(code)
	w.wroteHeader = true
	w.statusCode = code
}

func (w *responseWriter) EncodeJson(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (w *responseWriter) WriteJson(v interface{}) error {
	b, err := w.EncodeJson(v)
	if err != nil {
		return err
	}
	w.SetData(v)
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.rawBody = b
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	flusher := w.ResponseWriter.(http.Flusher)
	flusher.Flush()
}

func (w *responseWriter) CloseNotify() <-chan bool {
	notifier := w.ResponseWriter.(http.CloseNotifier)
	return notifier.CloseNotify()
}

func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker := w.ResponseWriter.(http.Hijacker)
	return hijacker.Hijack()
}
