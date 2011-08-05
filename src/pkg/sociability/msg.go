// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package sociability

import (
	//"log"
	"os"
	"github.com/petar/GoHTTP/server/rpc"
)

// AddMsg adds a new message to the database. The author is the currently
// logged in user. The message is attached to the object given by the string 
// argument "AttachTo". Optionally, the message is in response to another message
// with message ID "ReplyTo". AddMsg returns the message ID of the newly added
// message, in the return field "ID".
func (a *API) AddMsg(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	_, authorID, err := a.whoAmI(args)
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
	r.SetString("ID", WebStringOfObjectID(msgID))
	return nil
}

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
	return a.db.EditMsg(editorID, ObjectIDOfWebString(msg), body)
}

func (a *API) RemoveMsg(args *rpc.Args, r *rpc.Ret) (err os.Error) {
	_, editorID, err := a.whoAmI(args)
	if err != nil {
		return err
	}
	msg, err := args.QueryString("Msg")
	if err != nil {
		return err
	}
	return a.db.RemoveMsg(editorID, ObjectIDOfWebString(msg))
}

type msgJoinJSON struct {
	ID       string `json:"id"`
	Body     string `json:"body"`
	Author   string `json:"author_id"`
	//XXX: AuthorUser   string `json:"author_id"`
	AttachTo string `json:"attach"`
	ReplyTo  string `json:"reply"`
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
		q[i].ID = WebStringOfObjectID(join.ID)
		q[i].Body = join.Body
		q[i].Author = WebStringOfObjectID(join.Author)
		q[i].AttachTo = WebStringOfObjectID(join.AttachTo)
		q[i].ReplyTo = WebStringOfObjectID(join.ReplyTo)
	}
	r.SetInterface("Results", q)
	return nil
}
