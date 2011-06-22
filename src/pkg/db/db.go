// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package db

import (
	"log"
	"os"
	gobson "launchpad.net/gobson"
	"launchpad.net/mgo"
)

// Db encapsulates all interaction with the database
type Db struct {
	// Session with the Db server
	s *mgo.Session

	// Users collection
	u mgo.Collection
}

// NewDb creates a new Db interface to the database.
// addr is the IP address and port, in string form, of the database server.
func NewDb(addr string) (db *Db, err os.Error) {
	s, err := mgo.Mongo(addr)
	if err != nil {
		log.Printf("Problem connecting to MongoDB: %s", err)
		return nil, err
	}
	return &Db{ 
		s: s,
		u: s.DB("shelflife").C("users"),
	}, nil
}

type User struct {
	Login        string
	Email        string
	HashPassword string
}

func (db *Db) AddUser(u *User) os.Error {
	return db.u.Insert(u)
}

// UserByEmail looks up a user record with the given email.
// A non-nil error indicates a connectivity problem. 
// A missing user returns u == nil and err == nil.
func (db *Db) UserByEmail(email string) (u *User, err os.Error) {
	u = &User{}
	err = db.u.Find(gobson.M{ "Email": email }).One(u)
	if err == mgo.NotFound {
		return nil, nil
	}
	if err != nil {
		log.Printf("MongoDB error: %s\n", err)
		u = nil
	}
	return u, err
}

// UserByLogin looks up a user record with the given login (i.e. username).
// A non-nil error indicates a connectivity problem. 
// A missing user returns u == nil and err == nil.
func (db *Db) UserByLogin(login string) (u *User, err os.Error) {
	u = &User{}
	err = db.u.Find(gobson.M{ "Login": login }).One(u)
	if err == mgo.NotFound {
		return nil, nil
	}
	if err != nil {
		log.Printf("MongoDB error: %s\n", err)
		u = nil
	}
	return u, err
}

func (db *Db) Close() os.Error {
	db.s.Close()
	return nil
}
