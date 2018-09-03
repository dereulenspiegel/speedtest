package main

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/NebulousLabs/fastrand"
)

const (
	HeaderCacheControl            = "Cache-Control"
	HeaderPragma                  = "Pragma"
	HeaderContentDescription      = "Content-Description"
	HeaderContentType             = "Content-Type"
	HeaderContentDisposition      = "Content-Disposition"
	HeaderContentTransferEncoding = "Content-Transfer-Encoding"
)

const (
	ParamRequestedSize = "ckSize"
)

var (
	DefaultSize      = 100
	DefaultChunkSize = 1048576
)

var (
	invalidMethod = errors.New("Invalid method")
)

var (
	serveAddr = flag.String("http.addr", ":8080", "Listen addr for the http server")
)

func main() {
	flag.Parse()
	mux := http.NewServeMux()

	mux.Handle("/", serveStatic())
	mux.Handle("/empty.php", emptyHandler())
	mux.Handle("/garbage.php", garbageHandler())
	mux.Handle("/getIP.php", getIpHandler())

	server := &http.Server{
		Addr:    *serveAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("HTTP server failed: %s", err)
	}
}

func getIpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, invalidMethod.Error(), http.StatusMethodNotAllowed)
			return
		}
		clientAddr := ""
		if headerRealIp := r.Header.Get("X-Real-IP"); headerRealIp != "" {
			clientAddr = headerRealIp
		} else if headerRealIp = r.Header.Get("HTTP_X_FORWARDED_FOR"); headerRealIp != "" {
			clientAddr = headerRealIp
		} else {
			clientAddr = r.RemoteAddr
		}
		host, _, err := net.SplitHostPort(clientAddr)
		if err != nil {
			host = "unknown"
		}
		w.Write([]byte(host))
	})
}

func garbageHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, invalidMethod.Error(), http.StatusMethodNotAllowed)
			return
		}

		requestedSize := DefaultSize

		requestedSizeString := r.URL.Query().Get(ParamRequestedSize)
		if requestedSizeString != "" {
			var err error
			requestedSize, err = strconv.Atoi(requestedSizeString)
			if err != nil {
				log.Printf("Invalid parameter %s: %s (%s)", ParamRequestedSize, requestedSizeString, err)
				requestedSize = DefaultSize
			}
		}

		w.Header().Add(HeaderCacheControl, "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Add(HeaderCacheControl, "post-check=0, pre-check=0")
		w.Header().Add(HeaderContentDescription, "File Transfer")
		w.Header().Add(HeaderContentType, "application/octet-stream")
		w.Header().Add(HeaderContentDisposition, "attachment; filename=random.dat")
		w.Header().Add(HeaderContentTransferEncoding, "binary")
		w.Header().Add(HeaderPragma, "no-cache")
		for i := 0; i < requestedSize; i++ {
			randBytes := fastrand.Bytes(DefaultChunkSize)
			n, err := w.Write(randBytes)
			if err != nil {
				log.Printf("Failed to write to client: %s", err)
				return
			}
			if n != DefaultChunkSize {
				log.Printf("Failed to write all random bytes. Expected %d to be written. Actually wrote %d bytes", DefaultChunkSize, n)
				return
			}
		}
	})
}

func emptyHandler() http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodPost:
			w.Header().Add(HeaderCacheControl, "no-store, no-cache, must-revalidate, max-age=0")
			w.Header().Add(HeaderCacheControl, "post-check=0, pre-check=0")
			w.Header().Add(HeaderPragma, "no-cache")
			w.WriteHeader(http.StatusOK)
			return
		default:
			http.Error(w, invalidMethod.Error(), http.StatusMethodNotAllowed)
		}
	})

	/*handlerIfElse := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			return
		} else if r.Method == http.MethodPost {
			w.Header().Add(HeaderCacheControl, "no-store, no-cache, must-revalidate, max-age=0")
			w.Header().Add(HeaderCacheControl, "post-check=0, pre-check=0")
			w.Header().Add(HeaderPragma, "no-cache")
			w.WriteHeader(http.StatusOK)
			return
		} else {
			http.Error(w, invalidMethod.Error(), http.StatusMethodNotAllowed)
		}
	})*/

	return handler
}
