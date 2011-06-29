// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package db

import (
	"log"
	"os"
	"github.com/petar/ShelfLife/thirdparty/bson"
	"github.com/petar/ShelfLife/thirdparty/mgo"
)

type UserDoc struct {
	Name         string "name"
	Login        string "login"
	Email        string "email"
	HashPassword string "password"
}

// initUser configures the the 'user' collection
func (db *Db) initUser() os.Error {
	nt, err := db.kp.AddNodeType("user")
	if err != nil {
		return err
	}
	index := mgo.Index{
		Key: []string{"value.email", "value.login"},
		Unique: true,
		DropDups: false,
		Background: false,
		Sparse: false,
	}
	return nt.C.EnsureIndex(index)
}

// AddUser adds a new user, while returning an error if 'login' or 'email' are duplicates.
// Ir returns the user record's ID in the KPartite schema and any relevant error.
func (db *Db) AddUser(u *UserDoc) (string, os.Error) {
	return db.kp.AddNode("user", u)
}

type userFind struct {
	ID    string  "_id"
	Value UserDoc "value"
}

// FindUserByEmail looks up a user record with the given email.
// A non-nil error indicates a connectivity problem. 
// A missing user returns u == nil and err == nil.
func (db *Db) FindUserByEmail(email string) (u *UserDoc, err os.Error) {
	q, err := db.kp.FindNode("user", bson.D{{"email", email}})
	if err != nil {
		return nil, err
	}
	uf := &userFind{}
	err = q.One(uf)
	if err == mgo.NotFound {
		return nil, nil
	}
	if err != nil {
		log.Printf("DB error: %s\n", err)
		return nil, err
	}
	return &uf.Value, nil
}

// FindUserByLogin looks up a user record with the given login (i.e. username).
// A non-nil error indicates a connectivity problem. 
// A missing user returns u == nil and err == nil.
func (db *Db) FindUserByLogin(login string) (u *UserDoc, err os.Error) {
	q, err := db.kp.FindNode("user", bson.D{{"login", login}})
	if err != nil {
		return nil, err
	}
	uf := &userFind{}
	err = q.One(uf)
	if err == mgo.NotFound {
		return nil, nil
	}
	if err != nil {
		log.Printf("DB error: %s\n", err)
		return nil, err
	}
	return &uf.Value, nil
}
