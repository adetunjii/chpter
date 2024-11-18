package main

import (
	"github.com/chpter/user-svc/db"
)

func newTestServer(database db.UserQueries) *server {
	return &server{
		database: database,
	}
}
