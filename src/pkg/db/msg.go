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

// MsgDoc represents a user comment
type MsgDoc struct {
	Body string "body"
}

// initMsg adds and configures the messaging system database types
func (db *Db) initMsg() os.Error {
	// Add message node type
	_, err := db.kp.AddNodeType("msg")
	if err != nil {
		return err
	}
	// Edge type for connecting users to messges they write
	et, err := db.kp.AddEdgeType("msg_written_by", "msg", "user")
	if err != nil {
		return err
	}
	index := mgo.Index{
		Key:        []string{"_id", "from"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	if err = et.C.EnsureIndex(index); err != nil {
		return err
	}
	// Edge type for messages' in-response-to relationship
	et, err = db.kp.AddEdgeType("msg_replies_to", "msg", "msg")
	if err != nil {
		return err
	}
	index = mgo.Index{
		Key:        []string{"_id", "from"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	if err = et.C.EnsureIndex(index); err != nil {
		return err
	}
	// Edge type attaching a message to a foreign object
	et, err = db.kp.AddEdgeType("msg_attach_to", "msg", "foreign")
	if err != nil {
		return err
	}
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
	return nil
}

// AddMsg adds a new message, attached to attachTo foreign object, and in reply to an existing
// message with object ID relyTo (if replyTo is a valid object ID).
func (db *Db) AddMsg(user bson.ObjectId, attachTo string, replyTo bson.ObjectId, body string) os.Error {
	// Obtain object ID of foreign ID
	fID, err := db.addOrGetForeign(attachTo)
	if err != nil {
		return err
	}
	// Add message node
	msgID, err := db.kp.AddNode("msg", &MsgDoc{ Body: body })
	if err != nil {
		return err
	}
	// Add written-by relation
	// XXX: check that user exists?
	if _, err = db.kp.AddEdge("msg_written_by", msgID, user, nil); err != nil {
		return err
	}
	// Add attach-to relation
	if _, err = db.kp.AddEdge("msg_attach_to", msgID, fID, nil); err != nil {
		return err
	}
	// Add reply-to relation
	if replyTo.Valid() {
		// XXX: check that replyTo msg exists?
		if _, err = db.kp.AddEdge("msg_replies_to", msgID, replyTo, nil); err != nil {
			return err
		}
	}
	return nil
}

// EditMsg updates the body of the given message
func (db *Db) EditMsg(msg bson.ObjectId, body string) os.Error {
	return db.kp.UpdateNode("msg", msg, &MsgDoc{ Body: body })
}

// RemoveMsg removes a message and all incident edges
func (db *Db) RemoveMsg(msg bson.ObjectId) os.Error {
	return db.kp.RemoveNode("msg", msg)
}

// MsgJoin packs all message information in a single struct 
type MsgJoin struct {
	ID         string  // Message ID
	AttachToID string  // Foreign ID
	ReplyToID  string  // ID of message replying to
	AuthorID   string  // User ID of author
	Body       string
}

// FindMsgAttachedTo returns all messages attached to a given foreign object
func (db *Db) FindMsgAttachedTo(attachTo string) ([]*MsgJoin, os.Error) {
	// Obtain object ID of foreign ID
	fID, err := db.addOrGetForeign(attachTo)
	if err != nil {
		return err
	}
	// Find messages attached to this object
	q, err := db.kp.ArrivingEdges("msg_attach_to", fID)
	if err != nil {
		return nil, err
	}
	?
}
