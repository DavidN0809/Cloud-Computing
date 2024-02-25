// Package main implements a client for MovieInfo service.
package main

import (
	"context"
	"log"
	"time"

	"github.com/DavidN0809/Cloud-Computing/lab5/movieapi"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Context for the calls
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Example movie data to set
	movieData := &movieapi.MovieData{
		Title:    "Inception",
		Year:     2010,
		Director: "Christopher Nolan",
		Cast:     []string{"Leonardo DiCaprio", "Ellen Page", "Tom Hardy"},
	}
	// Set movie information
	_, err = c.SetMovieInfo(ctx, movieData)
	if err != nil {
		log.Fatalf("could not set movie info: %v", err)
	}
	log.Printf("Successfully set movie info for: %s", movieData.Title)

	// Now try to retrieve the movie info that we just set
	r, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: movieData.Title})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for %s: %d %s %v", movieData.Title, r.GetYear(), r.GetDirector(), r.GetCast())
}
