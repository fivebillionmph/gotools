package db

import (
	gss "github.com/fivebillionmph/gotools/server"
	"net/http"
	"errors"
)

type Server struct {
	db_types []DBType
}

func New_server() *Server {
	s := Server{
		db_types: make([]DBType, 0),
	}
	return &s
}

func (self *Server) Add_type(db_type DBType) error {
	for _, dt := range self.db_types {
		if dt.Name() == db_type.Name() {
			return errors.New("Duplicated type name: " + dt.Name())
		}
	}
	self.db_types = append(self.db_types, db_type)

	return nil
}

func (self *Server) Run_server(port int) error {
	var server *gss.Server

	server.AddRouterPath("/types", "GET", false, func(w http.ResponseWriter, r *http.Request) {
		type_names := make([]string, 0, len(self.db_types))
		for _, dt := range self.db_types {
			type_names = append(type_names, dt.Name())
		}
		gss.Send_json_response(w, 200, type_names)
	})

	server, err := gss.New("", port)
	if err != nil { return err }

	server.Start()

	return nil
}
