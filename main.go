package main

import (
	"fmt"
	"github.com/djumanoff/amqp"
	movie_store "github.com/kirigaikabuto/recommendation-system-movie-store"
	"log"
)

var (
	postgresUser         = "setdatauser"
	postgresPassword     = "123456789"
	postgresDatabaseName = "recommendation_system"
	postgresHost         = "localhost"
	postgresPort         = 5432
	postgresParams       = "sslmode=disable"
	amqpUrl              = "amqp://localhost:5672"
)

var cfg = amqp.Config{
	AMQPUrl:  amqpUrl,
	LogLevel: 5,
}

var srvCfg = amqp.ServerConfig{
	ResponseX: "response",
	RequestX:  "request",
}

func main() {
	sess := amqp.NewSession(cfg)

	if err := sess.Connect(); err != nil {
		fmt.Println(err)
		return
	}
	defer sess.Close()

	srv, err := sess.Server(srvCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	config := movie_store.Config{
		Host:     postgresHost,
		Port:     postgresPort,
		User:     postgresUser,
		Password: postgresPassword,
		Database: postgresDatabaseName,
		Params:   postgresParams,
	}
	movieStore, err := movie_store.NewPostgreStore(config)
	if err != nil {
		log.Fatal(err)
	}
	movieService := movie_store.NewMovieService(movieStore)
	moviesAmqpEndpoints := movie_store.NewAMQPEndpointFactory(movieService)
	srv.Endpoint("movie.get", moviesAmqpEndpoints.GetMovieByIdAMQPEndpoint())
	srv.Endpoint("movie.create", moviesAmqpEndpoints.CreateMovieAMQPEndpoint())
	srv.Endpoint("movie.list", moviesAmqpEndpoints.ListMoviesAMQPEndpoint())
	srv.Endpoint("movie.update", moviesAmqpEndpoints.UpdateProductAMQPEndpoint())
	srv.Endpoint("movie.delete", moviesAmqpEndpoints.DeleteMovieAMQPEndpoint())
	srv.Endpoint("movie.getByName", moviesAmqpEndpoints.GetMovieByNameAMQPEndpoint())
	fmt.Println("Start server")
	if err := srv.Start(); err != nil {
		fmt.Println(err)
		return
	}
}
