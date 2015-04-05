/*
This package contains the db helper functions
*/
package db

import (
	"gopkg.in/mgo.v2"
	"os"
)

var mongoSession *mgo.Session

var dbName = os.Getenv("LEMN_MG_DB_NAME")

// this initializes the db connection
// this *MUST* be called at the start of the web app
// currently this is being called from the main routine before
// starting the http handler
func InitMongo() error {
	var err error
	mongoSession, err = mgo.Dial(os.Getenv("LEMN_MG_URI"))
	if err != nil {
		return err
	}
	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	mongoSession.SetMode(mgo.Monotonic, true)
	return nil
}

// this returns the mongo db session copy
// every of the routine that is requesting this db session is responsible for closing this session
// New creates a new session with the same parameters as the original session, // including consistency, batch size, prefetching, safety mode, etc. The
// returned session will use sockets from the pool, so there's a chance that
// writes just performed in another session may not yet be visible.
// Login information from the original session will not be copied over into the
// new session unless it was provided through the initial URL for the Dial
// function.
// Copy is similar to new but it just maitains the auth information
func GetMongo() *mgo.Database {
	return mongoSession.Copy().DB(dbName)
}

func GetMongoSession() *mgo.Session {
	return mongoSession.Copy()
}

// this is for gracefull shutdown of db connection while shutting down the server
func CloseMongo() {
	mongoSession.Close()
}
