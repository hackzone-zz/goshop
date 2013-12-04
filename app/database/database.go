package database

import (
	//"log"
	"os"
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
)

type M map[string]interface{}

type cfg struct {
	session *mgo.Session
	db *mgo.Database
	host, database string
}

type FindMethods interface {
	One(result interface{}) error
	All(result interface{})
}

type finder struct {
	collection string
	query M
}

var c cfg


// Connect to a database
func Connect(host, database string) *mgo.Session {
	// connect to host
	session, err := mgo.Dial(host)

	if err != nil {
		panic(err)

		// kill everything
		os.Exit(1)
	}

	// save
	c = cfg{
		session  : session,
		host     : host,
		database : database,
		db       : session.DB(database),
	}

	return session
}

// Make a query
func Find(collection string, query M) FindMethods {
	return FindMethods(&finder{
		collection: collection,
		query     : query,
	})
}

func (f *finder) One(result interface{}) error {
	err := c.db.C(f.collection).Find(f.query).One(result)
	return err
}

func (f *finder) All(result interface{}) {
	c.db.C(f.collection).Find(f.query).All(result)
}

