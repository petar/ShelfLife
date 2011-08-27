// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package sociability

import (
	"log"
	"os"
	"time"
	"github.com/petar/GoHTTP/server/rpc"
	"github.com/petar/ShelfLife/thirdparty/bson"
)

const msgFormat = "02 Jan 3:04pm"

// AddMsg adds a new message to the database. The author is the currently
// logged in user. The message is attached to the object given by the string 
// argument "AttachTo". Optionally, the message is in response to another message
// with message ID "ReplyTo". AddMsg returns the message ID of the newly added
// message, in the return field "ID".
func (a *API) AddMsg(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	authorDoc, authorID, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	attachTo, err := args.QueryString("AttachTo")
	if err != nil || attachTo == "" {
		return ErrArg
	}
	replyTo, _ := args.QueryString("ReplyTo")
	body, err := args.QueryString("Body")
	if err != nil || body == "" {
		return ErrArg
	}
	msgID, err := a.db.AddMsg(authorID, attachTo, ObjectIDOfWebString(replyTo), body)
	if err != nil {
		return err
	}
	j := msgJoinJSON {
		ID:        WebStringOfObjectID(msgID),
		Body:      body,
		AuthorID:  WebStringOfObjectID(authorID),
		AuthorNym: authorDoc.Login,
		AttachTo:  attachTo,
		ReplyTo:   replyTo,
		Modified:  time.NanosecondsToLocalTime(int64(bson.Now())).Format(msgFormat),
	}
	r.SetInterface("Msg", j)
	return nil
}

// EditMsg changes the body of an existing message
func (a *API) EditMsg(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	_, editorID, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	msg, err := args.QueryString("Msg")
	if err != nil {
		return err
	}
	body, err := args.QueryString("Body")
	if err != nil || body == "" {
		return ErrArg
	}
	r.SetInt("XPad", 0)
	return a.db.EditMsg(editorID, ObjectIDOfWebString(msg), body)
}

// RemoveMsg deletes a message
func (a *API) RemoveMsg(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	_, editorID, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	msg, err := args.QueryString("Msg")
	if err != nil {
		return err
	}
	r.SetInt("XPad", 0)
	return a.db.RemoveMsg(editorID, ObjectIDOfWebString(msg))
}

type msgJoinJSON struct {
	ID        string `json:"id"`
	Body      string `json:"body"`
	AuthorID  string `json:"author_id"`
	AuthorNym string `json:"author_nym"`
	AttachTo  string `json:"attach"`
	ReplyTo   string `json:"reply"`
	Modified  string `json:"modified"`
}

func (a *API) FindMsgAttachedTo(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	attachTo, err := args.QueryString("AttachTo")
	if err != nil || attachTo == "" {
		return ErrArg
	}
	joins, err := a.db.FindMsgAttachedTo(attachTo)
	if err != nil {
		return err
	}
	q := make([]msgJoinJSON, len(joins))
	for i, join := range joins {
		author, err := a.whoIsID(join.Author)
		if err != nil {
			log.Printf("Unresolved author ID: %s", join.Author)
			q[i].AuthorNym = "anonymous"
		} else {
			q[i].AuthorNym = author.Login
		}
		q[i].ID = WebStringOfObjectID(join.ID)
		q[i].Body = join.Doc.Body
		q[i].AuthorID = WebStringOfObjectID(join.Author)
		q[i].AttachTo = WebStringOfObjectID(join.AttachTo)
		q[i].ReplyTo = WebStringOfObjectID(join.ReplyTo)
		modtm := time.NanosecondsToLocalTime(int64(join.Modified)).Format(msgFormat)
		q[i].Modified = modtm
	}
	r.SetInterface("Results", q)
	return nil
}
