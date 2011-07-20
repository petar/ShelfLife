// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package db

import (
	//"log"
	"os"
	"github.com/petar/ShelfLife/thirdparty/bson"
	"github.com/petar/ShelfLife/thirdparty/mgo"
)

// Node document for a foreign (external) object, described by
// an opaque string foreign ID.
type ForeignDoc struct {
	ForeignID string "fid"
}

// Edge document for user that likes foreign object. The foreign
// ID is duplicated on the edge for faster queries. It is also present
// in the target node document.
type LikeForeignDoc struct {
	ForeignID string "fid"
}

// initLike adds and configures the like system database types
func (db *Db) initLike() os.Error {
	// external is a node type that represents opaque external string IDs 
	nt, err := db.kp.AddNodeType("foreign")
	if err != nil {
		return err
	}
	index := mgo.Index{
		Key: []string{"value.fid"},
		Unique: true,
		DropDups: false,
		Background: false,
		Sparse: false,
	}
	err = nt.C.EnsureIndex(index)
	if err != nil {
		return err
	}
	// Edge type for user likes of foreign objects
	et, err := db.kp.AddEdgeType("like_foreign", "user", "foreign")
	if err != nil {
		return err
	}
	index = mgo.Index{
		Key: []string{"_id"},
		Unique: true,
		DropDups: false,
		Background: false,
		Sparse: false,
	}
	return et.C.EnsureIndex(index)
}

type foreignFind struct {
	ID    bson.ObjectId  "_id"
	Value ForeignDoc     "value"
}

// Like updates the databse to indicate that user like foreign object fid.
// It is idempotent (multiple calls are OK).
func (db *Db) Like(user bson.ObjectId, fid string) os.Error {
	id, err := db.addOrGetForeign(fid)
	if err != nil {
		return err
	}
	_, err = db.kp.AddOrReplaceEdge("like_foreign", user, id, &LikeForeignDoc{ ForeignID: fid })
	return err
}

// addOrGetForeign create a node for the given foreign ID and returns its node object ID
func (db *Db) addOrGetForeign(fid string) (bson.ObjectId, os.Error) {
	q, err := db.kp.FindNode("foreign", bson.D{{"fid", fid}})
	if err != nil {
		return "", err
	}
	ff := &foreignFind{}
	if err = q.One(ff); err != nil {
		return db.kp.AddNode("foreign", &ForeignDoc{ ForeignID: fid })
	}
	return ff.ID, nil
}

// Unlike updates the databse to indicate that user does not like foreign object fid.
// It is idempotent (multiple calls are OK).
func (db *Db) Unlike(user bson.ObjectId, fid string) os.Error {
	id, err := db.addOrGetForeign(fid)
	if err != nil {
		return err
	}
	return db.kp.RemoveEdgeAnchors("like_foreign", user, id)
}

// LikeCount returns the number of users that like the foreign object fid
func (db *Db) LikeCount(fid string) (int, os.Error) {
	q, err := db.kp.FindEdge("like_foreign", bson.D{{"fid", fid}})
	if err != nil {
		return 0, err
	}
	return q.Count()
}

// Likes returns true, if the user with the given node ID likes the foreign object fid
func (db *Db) Likes(user bson.ObjectId, fid string) (bool, os.Error) {
	id, err := db.addOrGetForeign(fid)
	if err != nil {
		return false, err
	}
	return db.kp.IsEdge("like_foreign", user, id)
}
