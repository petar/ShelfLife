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

type User struct {
	Name         string
	Login        string
	Email        string
	HashPassword string
}

func (db *Db) AddUser(u *User) os.Error {
	return db.u.Insert(u)
}

// FindUserByEmail looks up a user record with the given email.
// A non-nil error indicates a connectivity problem. 
// A missing user returns u == nil and err == nil.
func (db *Db) FindUserByEmail(email string) (u *User, err os.Error) {
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

// FindUserByLogin looks up a user record with the given login (i.e. username).
// A non-nil error indicates a connectivity problem. 
// A missing user returns u == nil and err == nil.
func (db *Db) FindUserByLogin(login string) (u *User, err os.Error) {
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
