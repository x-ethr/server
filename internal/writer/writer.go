package writer

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type Writer struct {
	w http.ResponseWriter

	status int
	buffer *bytes.Buffer
}

func Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		instance := &Writer{w: w, status: 200} // default 200 response code

		next.ServeHTTP(instance, r)

		log.Printf("response size: %d\n", instance.buffer.Len())
		log.Printf("response status: %v\n", instance.status)

		size, e := instance.Done()
		if e != nil {
			log.Println("error while writing response", e)
		}

		log.Printf("response size: %d\n", size)
	})
}

func (w *Writer) Header() http.Header {
	return w.w.Header()
}

func (w *Writer) Write(bytes []byte) (int, error) {
	return w.buffer.Write(bytes)
}

func (w *Writer) WriteHeader(status int) {
	w.status = status
}

func (w *Writer) Done() (int64, error) {
	if w.status > 0 {
		w.w.WriteHeader(w.status)
	}

	return io.Copy(w.w, w.buffer)
}
