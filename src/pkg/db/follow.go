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

type NotifyMsgDoc struct {
	Body string `bson:"body"`
}

// initFollow adds and configures the follow/notify system
func (db *Db) initFollow() os.Error {
	// Edge type for connecting users to objects they follow
	et, err := db.kp.AddEdgeType("notify_follow", "user", "foreign")
	if err != nil {
		return err
	}
	index := mgo.Index{
		Key:        []string{"_id"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	if err = et.C.EnsureIndex(index); err != nil {
		return err
	}
	// Node type for a notification
	_, err = db.kp.AddNodeType("notify_msg")
	if err != nil {
		return err
	}
	// Edge type object->notify_msg linking objects to notification messages
	et, err = db.kp.AddEdgeType("notify_push", "foreign", "notify_msg")
	if err != nil {
		return err
	}
	// Ensure simple graph
	index = mgo.Index{
		Key:        []string{"_id"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	if err = et.C.EnsureIndex(index); err != nil {
		return err
	}
	// Ensure efficient query 'list all messages for object X' and sort by creation time
	index = mgo.Index{
		Key:        []string{"from", "created"},
		Unique:     false,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	if err = et.C.EnsureIndex(index); err != nil {
		return err
	}
	// Edge type user->notification showing a user's notifications list
	et, err = db.kp.AddEdgeType("notify_pull", "user", "notify_msg")
	if err != nil {
		return err
	}
	// Ensure simple graph
	index = mgo.Index{
		Key:        []string{"_id"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	if err = et.C.EnsureIndex(index); err != nil {
		return err
	}
	// Ensure efficient query 'list all users that follow X' and sort by creation time
	index = mgo.Index{
		Key:        []string{"from", "created"},
		Unique:     false,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	if err = et.C.EnsureIndex(index); err != nil {
		return err
	}
	return nil
}

func (db *Db) SetFollow(user bson.ObjectId, foreign string) os.Error {
	id, err := db.GetOrMakeForeignID(foreign)
	if err != nil {
		return err
	}
	_, err = db.kp.AddOrReplaceEdge("notify_follow", user, id, nil)
	return err
}

func (db *Db) UnsetFollow(user bson.ObjectId, foreign string) os.Error {
	id, err := db.GetOrMakeForeignID(foreign)
	if err != nil {
		return err
	}
	return db.kp.RemoveEdgeAnchors("notify_follow", user, id)
}

func (db *Db) IsFollow(user bson.ObjectId, foreign string) (bool, os.Error) {
	id, err := db.GetOrMakeForeignID(foreign)
	if err != nil {
		return false, err
	}
	return db.kp.IsEdge("notify_follow", user, id)
}

func (db *Db) FollowerCount(foreign string) (int, os.Error) {
	id, err := db.GetOrMakeForeignID(foreign)
	if err != nil {
		return 0, err
	}
	return db.kp.ArrivingDegree("notify_follow", id)
}

func (db *Db) ListFollowed(user bson.ObjectId) ([]bson.ObjectId, os.Error) {
	q, err := db.kp.LeavingEdges("notify_follow", user)
	if err != nil {
		return nil, err
	}
	n, err := q.Count()
	if err != nil {
		return nil, err
	}
	r := make([]bson.ObjectId, 0, n)
	iter, err := q.Iter()
	if err != nil {
		return nil, err
	}
	edgeDoc := &EdgeDoc{}
	for err = iter.Next(edgeDoc); err != nil; err = iter.Next(edgeDoc) {
		r = append(r, edgeDoc.To)
	}
	return r, nil
}

/*
func (db *Db) AddNotification(obj bson.ObjectId, msgDoc *NotifyMsgDoc) os.Error {
}

func (db *Db) RemoveNotification(msg bson.ObjectId) os.Error {
}

type NotifyJoin struct {
}

func (db *Db) FetchNotifications(user bson.ObjectId) ([]*NotifyJoin, os.Error) {
}
*/
