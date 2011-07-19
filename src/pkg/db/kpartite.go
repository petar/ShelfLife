// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package db

import (
	"crypto/sha1"
	"hash"
	"os"
	"rand"
	"sync"
	"time"
	"github.com/petar/ShelfLife/thirdparty/bson"
	"github.com/petar/ShelfLife/thirdparty/mgo"
)

// KPartite is a thin layer in front of the database, which acts
// as a k-partite directed graph database. Values are stored at the
// nodes and edges.
type KPartite struct {
	name      string
	s         *mgo.Session
	sync.Mutex
	edgeHash  hash.Hash
	nodeTypes map[string]*NodeType
	edgeTypes map[string]*EdgeType
}

const IDLEN = 12

var (
	ErrArg  = os.NewError("bad argument")
	ErrDup  = os.NewError("duplicate")
	ErrType = os.NewError("unknown type")
)

func NewKPartite(name string, session *mgo.Session) *KPartite {
	return &KPartite{
		s:         session,
		name:      name,
		edgeHash:  sha1.New(),
		nodeTypes: make(map[string]*NodeType),
		edgeTypes: make(map[string]*EdgeType),
	}
}

// Type functions

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
		C:    kp.s.DB(kp.name).C(name),
	}
	kp.nodeTypes[name] = nt

	return nt, nil
}

func (kp *KPartite) GetNodeType(name string) *NodeType {
	kp.Lock()
	defer kp.Unlock()

	return kp.nodeTypes[name]
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
	et = &EdgeType{
		Name: name,
		C:    kp.s.DB(kp.name).C(name),
		From: ft,
		To:   tt,
	}
	kp.edgeTypes[name] = et

	return et, nil
}

func (kp *KPartite) GetEdgeType(name string) *EdgeType {
	kp.Lock()
	defer kp.Unlock()

	return kp.edgeTypes[name]
}

func (kp *KPartite) GetEdgeTypes() []*EdgeType {
	kp.Lock()
	defer kp.Unlock()

	ets := make([]*EdgeType, len(kp.edgeTypes))
	i := 0
	for _, et := range kp.edgeTypes {
		ets[i] = et
		i++
	}

	return ets
}

// Node functions

type NodeDoc struct {
	ID    bson.ObjectId "_id"
	Value interface{}   "value"
}

func (kp *KPartite) AddNode(nodeType string, value interface{}) (bson.ObjectId, os.Error) {
	nd := &NodeDoc{
		ID:    chooseID(),
		Value: value,
	}
	nt := kp.GetNodeType(nodeType)
	if nt == nil {
		return "", ErrArg
	}
	return nd.ID, nt.C.Insert(nd)
}

func (kp *KPartite) UpdateNode(nodeType string, nodeID bson.ObjectId, value interface{}) os.Error {
	nt := kp.GetNodeType(nodeType)
	if nt == nil {
		return ErrType
	}
	return nt.C.Update(bson.D{{"_id", nodeID}}, bson.D{{"value", value}})
}

func (kp *KPartite) FindNode(nodeType string, query interface{}) (*mgo.Query, os.Error) {
	nt := kp.GetNodeType(nodeType)
	if nt == nil {
		return nil, ErrType
	}
	return nt.C.Find(rewriteQuery(query)), nil
}

func rewriteQuery(q interface{}) interface{} {
	dq := q.(bson.D)
	for i, _ := range dq {
		dq[i].Name = "value." + dq[i].Name
	}
	return dq
}

func chooseID() bson.ObjectId {
	rand.Seed(time.Nanoseconds())
	b := make([]byte, IDLEN)
	for i := 0; i < IDLEN/4; i++ {
		u := rand.Uint32()
		b[4*i] = byte(u & 0xff)
		b[4*i+1] = byte((u >> 8) & 0xff)
		b[4*i+2] = byte((u >> 16) & 0xff)
		b[4*i+3] = byte((u >> 24) & 0xff)
	}
	return bson.ObjectId(b)
}

// Edge functions

type EdgeDoc struct {
	ID    bson.ObjectId "_id"
	From  bson.ObjectId "from"
	To    bson.ObjectId "to"
	Value interface{}   "value"
}

func (kp *KPartite) makeEdgeID(from, to bson.ObjectId) bson.ObjectId {
	kp.Lock()
	kp.edgeHash.Reset()
	kp.edgeHash.Write([]byte(from))
	kp.edgeHash.Write([]byte(to))
	h := kp.edgeHash.Sum()
	kp.Unlock()
	return bson.ObjectId(h[:IDLEN])
}

func (kp *KPartite) AddEdge(edgeType string, from, to bson.ObjectId, value interface{}) (bson.ObjectId, os.Error) {
	ed := &EdgeDoc{
		ID:    kp.makeEdgeID(from, to),
		From:  from,
		To:    to,
		Value: value,
	}
	et := kp.GetEdgeType(edgeType)
	if et == nil {
		return "", ErrArg
	}
	return ed.ID, et.C.Insert(ed)
}

func (kp *KPartite) AddOrReplaceEdge(edgeType string, from, to bson.ObjectId, value interface{}) (bson.ObjectId, os.Error) {
	ed := &EdgeDoc{
		ID:    kp.makeEdgeID(from, to),
		From:  from,
		To:    to,
		Value: value,
	}
	et := kp.GetEdgeType(edgeType)
	if et == nil {
		return "", ErrArg
	}
	_, err := et.C.Upsert(bson.D{{"_id", ed.ID}}, ed)
	return ed.ID, err
}

func (kp *KPartite) UpdateEdge(edgeType string, edgeID bson.ObjectId, value interface{}) os.Error {
	et := kp.GetEdgeType(edgeType)
	if et == nil {
		return ErrArg
	}
	return et.C.Update(bson.D{{"_id", edgeID}}, bson.D{{"value", value}})
}

func (kp *KPartite) UpdateEdgeAnchors(edgeType string, from, to bson.ObjectId, value interface{}) os.Error {
	return kp.UpdateEdge(edgeType, kp.makeEdgeID(from, to), value)
}

func (kp *KPartite) FindEdge(edgeType string, query interface{}) (*mgo.Query, os.Error) {
	et := kp.GetEdgeType(edgeType)
	if et == nil {
		return nil, ErrType
	}
	return et.C.Find(rewriteQuery(query)), nil
}

func (kp *KPartite) IsEdge(edgeType string, from, to bson.ObjectId) (bool, os.Error) {
	et := kp.GetEdgeType(edgeType)
	if et == nil {
		return false, ErrArg
	}
	n, err := et.C.Find(bson.D{{"_id", kp.makeEdgeID(from, to)}}).Count()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (kp *KPartite) RemoveEdge(edgeType string, edgeID bson.ObjectId) os.Error {
	et := kp.GetEdgeType(edgeType)
	if et == nil {
		return ErrArg
	}
	return et.C.Remove(bson.D{{"_id", edgeID}})
}

func (kp *KPartite) RemoveEdgeAnchors(edgeType string, from, to bson.ObjectId) os.Error {
	return kp.RemoveEdge(edgeType, kp.makeEdgeID(from, to))
}

// Node-edge functions

func (kp *KPartite) RemoveNode(nodeType string, nodeID bson.ObjectId) os.Error {
	nt := kp.GetNodeType(nodeType)
	if nt == nil {
		return ErrArg
	}
	if err := nt.C.Remove(bson.D{{"_id", nodeID}}); err != nil {
		return err
	}
	ets := kp.GetEdgeTypes()
	for _, et := range ets {
		et.C.RemoveAll(bson.D{{"from", nodeID}})
		et.C.RemoveAll(bson.D{{"to", nodeID}})
	}

	return nil
}

func (kp *KPartite) LeavingEdges(edgeType string, from bson.ObjectId) (*mgo.Query, os.Error) {
	et := kp.GetEdgeType(edgeType)
	if et == nil {
		return nil, ErrType
	}
	return et.C.Find(bson.D{{"from", from}}), nil
}

func (kp *KPartite) ArrivingEdges(edgeType string, to bson.ObjectId) (*mgo.Query, os.Error) {
	et := kp.GetEdgeType(edgeType)
	if et == nil {
		return nil, ErrType
	}
	return et.C.Find(bson.D{{"to", to}}), nil
}

func (kp *KPartite) LeavingDegree(edgeType string, from bson.ObjectId) (int, os.Error) {
	q, err := kp.LeavingEdges(edgeType, from)
	if err != nil {
		return 0, err
	}
	return q.Count()
}

func (kp *KPartite) ArrivingDegree(edgeType string, to bson.ObjectId) (int, os.Error) {
	q, err := kp.ArrivingEdges(edgeType, to)
	if err != nil {
		return 0, err
	}
	return q.Count()
}
