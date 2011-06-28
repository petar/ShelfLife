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

// KPartite is a thin layer in front of the database, which acts
// as a k-partite directed graph database. Values are stored at the
// nodes and edges.
type KPartite struct {
	sync.Mutex
	name      string
	s         *mgo.Session
	nodeTypes map[string]*NodeType
	edgeTypes map[string]*EdgeType
}

var (
	ErrArg = os.NewError("bad argument")
	ErrDup = os.NewError("duplicate")
)

func NewKPartite(name string, session *mgo.Session) *KPartite {
	return &KPartite{
		s:         session,
		name:      name,
		nodeTypes: make(map[string]*NodeType),
		edgeTypes: make(map[string]*EdgeType),
	}
}

// Type manipulation

type NodeType struct {
	Name string
	C    mgo.Collection
}

func (kp *KPartite) AddNodeType(name string) (*NodeType, os.Error) {
	kp.Lock()
	defer kp.Unlock()

	nt, ok := kp.nodeTypes[name]
	if ok {
		return nil, ErrDup
	}
	nt = &NodeType{
		Name: name,
		C:    s.DB(kp.name).C(name),
	}

	return nt, nil
}

func (kp *KPartite) GetNodeType(name string) *NodeType {
	kp.Lock()
	defer kp.Unlock()

	return kp.nodeTypes(name)
}

type EdgeType struct {
	Name string
	C    mgo.Collection
	From *NodeType
	To   *NodeType
}

func (kp *KPartite) AddEdgeType(name string, from, to string) (*EdgeType, os.Error) {
	kp.Lock()
	defer kp.Unlock()

	et, ok := kp.edgeTypes[name]
	if ok {
		return nil, ErrDup
	}
	ft, ok := kp.nodeTypes[from]
	if !ok {
		return nil, ErrArg
	}
	tt, ok := kp.nodeTypes[to]
	if !ok {
		return nil, ErrArg
	}
	et = &NodeType{
		Name: name,
		C:    s.DB(kp.name).C(name),
		From: ft,
		To:   tt,
	}

	return et, nil
}

func (kp *KPartite) GetEdgeType(name string) *EdgeType {
	kp.Lock()
	defer kp.Unlock()

	return kp.edgeTypes(name)
}

// Node manipulation

type NodeDoc struct {
	ID    string      "_id"
	Value interface{} "value"
}

func (kp *KPartite) AddNode() {
}

func (kp *KPartite) RemoveNode() {
}

// Edge manipulation

type EdgeDoc struct {
	ID    string      "_id"
	From  string      "from"
	To    string      "to"
	Value interface{} "value"
}

func (kp *KPartite) AddEdge() {
}

func (kp *KPartite) RemoveEdge() {
}
