package main

import (
	"github.com/labstack/echo"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
	"quandat10/htttdl/backend/api"
)

func main() {

	store := driver("neo4j://localhost:7687", neo4j.BasicAuth("neo4j", "quandat10", ""))
	router := echo.New()
	server, err := api.NewServer(store, router)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
	err = server.StartServer(":8001")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
	defer store.Close()
}

func driver(target string, token neo4j.AuthToken) neo4j.Driver {
	result, err := neo4j.NewDriver(target, token)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect database")
		panic(err)
	}
	return result
}
