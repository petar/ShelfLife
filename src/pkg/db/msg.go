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

// MsgDoc represents a user comment
type MsgDoc struct {
	Body string `bson:"body"`
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
func (db *Db) AddMsg(user bson.ObjectId, attachTo string, replyTo bson.ObjectId, body string) (bson.ObjectId, os.Error) {
	// Obtain object ID of foreign ID
	fID, err := db.GetOrMakeForeignID(attachTo)
	if err != nil {
		return "", err
	}
	// Add message node
	msgID, err := db.kp.AddNode("msg", &MsgDoc{ Body: body })
	if err != nil {
		return "", err
	}
	// Add written-by relation
	// XXX: check that user exists?
	if _, err = db.kp.AddEdge("msg_written_by", msgID, user, nil); err != nil {
		return "", err
	}
	// Add attach-to relation
	if _, err = db.kp.AddEdge("msg_attach_to", msgID, fID, nil); err != nil {
		return "", err
	}
	// Add reply-to relation
	if replyTo.Valid() {
		// XXX: check that replyTo msg exists and attaches to same object
		if _, err = db.kp.AddEdge("msg_replies_to", msgID, replyTo, nil); err != nil {
			return "", err
		}
	}
	return msgID, nil
}

// EditMsg updates the body of the given message
func (db *Db) EditMsg(editorID, msgID bson.ObjectId, body string) os.Error {
	join, err := db.joinMsg(msgID)
	if err != nil {
		return err
	}
	if join.Author != editorID {
		return ErrSec
	}
	join.Doc.Body = body
	return db.kp.UpdateNode("msg", msgID, &join.Doc)
}

// RemoveMsg removes a message and all incident edges
func (db *Db) RemoveMsg(editorID, msgID bson.ObjectId) os.Error {
	join, err := db.joinMsg(msgID)
	if err != nil {
		return err
	}
	if join.Author != editorID {
		return ErrSec
	}
	return db.kp.RemoveNode("msg", msgID)
}

// FindMsgAttachedTo returns all messages attached to a given foreign object
func (db *Db) FindMsgAttachedTo(attachTo string) ([]*MsgJoin, os.Error) {
	// Obtain object ID of foreign ID
	fID, err := db.GetOrMakeForeignID(attachTo)
	if err != nil {
		return nil, err
	}
	// Find messages attached to this object
	q, err := db.kp.ArrivingEdges("msg_attach_to", fID)
	if err != nil {
		return nil, err
	}
	iter, err := q.Iter()
	if err != nil {
		return nil, err
	}
	r := make([]*MsgJoin, 0)
	edgeDoc := &EdgeDoc{}
	for err = iter.Next(edgeDoc); err != nil; err = iter.Next(edgeDoc) {
		join, err := db.joinMsg(edgeDoc.From)
		if err != nil {
			continue
		}
		r = append(r, join)
	}
	return r, nil
}

// MsgJoin packs all message information in a single struct 
type MsgJoin struct {
	ID       bson.ObjectId  // msg_ID
	Doc      MsgDoc
	Author   bson.ObjectId  // user_ID of author
	AttachTo bson.ObjectId  // foreign_ID of object the message is attached to
	ReplyTo  bson.ObjectId  // msg_ID of message replying to
}

func (db *Db) joinMsg(msgID bson.ObjectId) (*MsgJoin, os.Error) {
	join := &MsgJoin{}
	join.ID = msgID

	// Get body
	nd, err := db.kp.FindNode("msg", bson.D{{"_id", msgID}})
	if err != nil {
		return nil, err
	}
	join.Doc.Body, _ = (nd.Value).(bson.M)["body"].(string)

	// Find Author
	ed, err := db.kp.LeavingEdge("msg_written_by", msgID)
	if err != nil {
		return nil, err
	}
	join.Author = ed.To

	// Find attach-to object
	ed, err = db.kp.LeavingEdge("msg_attach_to", msgID)
	if err != nil {
		return nil, err
	}
	join.AttachTo = ed.To

	// Find ReplyTo message
	ed, err = db.kp.LeavingEdge("msg_replies_to", msgID)
	if err != nil {
		join.ReplyTo = ""
	} else {
		join.ReplyTo = ed.To
	}
	log.Printf("ok")

	return join, nil
}
