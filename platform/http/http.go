package http

import (
	"fmt"
	glog "log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type FwdToZeroWriter struct {
}

func (fw *FwdToZeroWriter) Write(p []byte) (n int, err error) {
	log.Error().Msg(string(p))
	return len(p), nil
}

func StartWebService(port int64) error {
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			normalizedPath := strings.TrimSpace(r.URL.Path)
			if !strings.HasPrefix(normalizedPath, "/") {
				normalizedPath = "/" + normalizedPath
			}
			r.URL.Path = normalizedPath
			next.ServeHTTP(w, r)
		})
	})
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(FilterApi)
	r.Get("/cdn/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	r.Get("/resource/{service}/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	})
	r.Mount("/api", adminRouter())
	r.Get("/{endpoint}/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Index"))
	})
	h2s := &http2.Server{}
	httpServer := &http.Server{
		Addr:     fmt.Sprintf(":%d", port),
		Handler:  h2c.NewHandler(r, h2s),
		ErrorLog: glog.New(&FwdToZeroWriter{}, "", 0),
	}
	go func() {
		log.Info().Int("port", int(port)).Msg("start http server")
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()
	return nil
}

func FilterApi(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqTyp := strings.TrimSpace(r.Header.Get("X-Request-Type"))
		switch strings.ToLower(reqTyp) {
		case "api":
			apiProxy(w, r)
		default:
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

func apiProxy(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(300 * time.Second))
	r.Get("/deployments", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/navigations", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/proxies", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	return r
}
