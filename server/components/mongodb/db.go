package mongodb

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"gopkg.in/mgo.v2"
)

var (
	ErrMongoDBDisabled = fmt.Errorf("MongoDB is not enabled.")
)

func DBConnect(address string) *mgo.Session {
	log.Printf("MongoDB connecting to %s\n", address)
	session, err := mgo.Dial(address)
	if err != nil {
		log.Fatalf("Can't connect to mongodb: %v", err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("%v captured - Closing database connection\n", sig)
			session.Close()
			os.Exit(0)
		}
	}()

	return session
}
