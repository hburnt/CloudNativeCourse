// Package main imlements a client for movieinfo service
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/hburnt/CloudNativeCourse/lab5/movieapi"
	"google.golang.org/grpc"
)

const (
	address      = "localhost:50051"
	defaultTitle = "Pulp fiction"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Contact the server and print out its response.
	title := defaultTitle
	if len(os.Args) > 1 {
		title = os.Args[1]
	}
	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for %s %d %s %v", title, r.GetYear(), r.GetDirector(), r.GetCast())

  req, err := c.SetMovieInfo(ctx, &movieapi.MovieData{Title: "Iron Man", Year: 2008, Director: "Jon Favreau", Cast: []string{"Robert Downey Jr.", "Jon Favreau", "Jeff Bridges"}})

  if err != nil {
    log.Fatalf("Could not set movie info: %w", err)
  }
 
log.Printf("Set movie info: %s", req.GetCode())
}
