package server

import (
	"errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	http_server *http.ServeMux
	logger      *log.Logger
	router      *mux.Router
	port        int
}

func New(logfile string, port int) (*Server, error) {
	var f io.Writer
	var err error
	if logfile != "" {
		f, err = os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			return nil, err
		}
	} else {
		f = os.Stdout
	}
	logger := log.New(f, "--- ", log.Ldate|log.Ltime|log.Lshortfile)
	http_server := http.NewServeMux()
	router := mux.NewRouter()

	server := &Server{http_server, logger, router, port}
	return server, nil
}

func (self *Server) Start() error {
	self.http_server.Handle("/", self.router)
	port_str := "127.0.0.1:" + strconv.Itoa(self.port)
	self.logger.Println("starting server on port " + port_str)
	return http.ListenAndServe(port_str, self.http_server)
}

func (self *Server) AddRouterPath(path string, method string, prefix bool, handler func(http.ResponseWriter, *http.Request)) error {
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		return errors.New("method must be get or post")
	}

	parent_handler := func(w http.ResponseWriter, r *http.Request) {
		self.logger.Printf("%s\t%s\t%s\n", r.Method, r.Header.Get("X-Forwarded-For"), r.URL.String())
		handler(w, r)
	}

	if prefix {
		self.router.PathPrefix(path).HandlerFunc(parent_handler).Methods(method)
	} else {
		self.router.Path(path).HandlerFunc(parent_handler).Methods(method)
	}

	return nil
}

func (self *Server) AddStaticRouterPathPrefix(path string, dir string, before func(http.ResponseWriter, *http.Request) bool) error {
	file_server := http.StripPrefix(path + "/", http.FileServer(http.Dir(dir)))
	if before == nil {
		self.router.PathPrefix(path).Handler(file_server)
	} else {
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			cont := before(w, r)
			if cont {
				file_server.ServeHTTP(w, r)
			}
		}
		self.router.PathPrefix(path).Handler(handler)
	}
	return nil
}

func (self *Server) AddSingleFilePath(path string, html_file string, prefix bool, before func(http.ResponseWriter, *http.Request) bool) error {
	file_handler := func(w http.ResponseWriter, r *http.Request) {
		self.logger.Printf("%s\t%s\t%s\n", r.Method, r.Header.Get("X-Forwarded-For"), r.URL.String())
		http.ServeFile(w, r, html_file)
	}

	var parent_handler func(w http.ResponseWriter, r *http.Request)
	if before == nil {
		parent_handler = file_handler
	} else {
		parent_handler = func(w http.ResponseWriter, r *http.Request) {
			cont := before(w, r)
			if cont {
				file_handler(w, r)
			}
		}
	}

	if prefix {
		self.router.PathPrefix(path).HandlerFunc(parent_handler)
	} else {
		self.router.Path(path).HandlerFunc(parent_handler)
	}

	return nil
}
