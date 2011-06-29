// mgo - MongoDB driver for Go
// 
// Copyright (c) 2010-2011 - Gustavo Niemeyer <gustavo@niemeyer.net>
// 
// All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// 
//     * Redistributions of source code must retain the above copyright notice,
//       this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright notice,
//       this list of conditions and the following disclaimer in the documentation
//       and/or other materials provided with the distribution.
//     * Neither the name of the copyright holder nor the names of its
//       contributors may be used to endorse or promote products derived from
//       this software without specific prior written permission.
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR
// CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
// EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
// PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package mgo_test

import (
	"github.com/petar/ShelfLife/thirdparty/bson"
	. "launchpad.net/gocheck"
	"launchpad.net/mgo"
	"os"
)

func (s *S) TestGridFSCreate(c *C) {
	session, err := mgo.Mongo("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	db := session.DB("mydb")

	before := bson.Now()

	gfs := db.GridFS("fs")
	file, err := gfs.Create("")
	c.Assert(err, IsNil)

	n, err := file.Write([]byte("some data"))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 9)

	err = file.Close()
	c.Assert(err, IsNil)

	after := bson.Now()

	// Check the file information.
	result := M{}
	err = db.C("fs.files").Find(nil).One(result)
	c.Assert(err, IsNil)

	fileId, ok := result["_id"].(bson.ObjectId)
	c.Assert(ok, Equals, true)
	c.Assert(fileId.Valid(), Equals, true)
	result["_id"] = "<id>"

	fileTs, ok := result["uploadDate"].(bson.Timestamp)
	c.Assert(ok, Equals, true)
	c.Assert(fileTs >= before && fileTs <= after, Equals, true)
	result["uploadDate"] = "<timestamp>"

	expected := M{
		"_id":        "<id>",
		"length":     9,
		"chunkSize":  262144,
		"uploadDate": "<timestamp>",
		"md5":        "1e50210a0202497fb79bc38b6ade6c34",
	}
	c.Assert(result, Equals, expected)

	// Check the chunk.
	result = M{}
	err = db.C("fs.chunks").Find(nil).One(result)
	c.Assert(err, IsNil)

	chunkId, ok := result["_id"].(bson.ObjectId)
	c.Assert(ok, Equals, true)
	c.Assert(chunkId.Valid(), Equals, true)
	result["_id"] = "<id>"

	expected = M{
		"_id":      "<id>",
		"files_id": fileId,
		"n":        0,
		"data":     []byte("some data"),
	}
	c.Assert(result, Equals, expected)

	// Check that an index was created.
	indexes, err := db.C("fs.chunks").Indexes()
	c.Assert(err, IsNil)
	c.Assert(len(indexes), Equals, 2)
	c.Assert(indexes[1].Key, Equals, []string{"files_id", "n"})
}

func (s *S) TestGridFSFileDetails(c *C) {
	session, err := mgo.Mongo("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	db := session.DB("mydb")

	gfs := db.GridFS("fs")

	file, err := gfs.Create("myfile1.txt")
	c.Assert(err, IsNil)

	n, err := file.Write([]byte("some"))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 4)

	c.Assert(file.Size(), Equals, int64(4))

	n, err = file.Write([]byte(" data"))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 5)

	c.Assert(file.Size(), Equals, int64(9))

	id, _ := file.Id().(bson.ObjectId)
	c.Assert(id.Valid(), Equals, true)
	c.Assert(file.Name(), Equals, "myfile1.txt")
	c.Assert(file.ContentType(), Equals, "")

	var info interface{}
	err = file.GetInfo(&info)
	c.Assert(err, IsNil)
	c.Assert(info, IsNil)

	file.SetId("myid")
	file.SetName("myfile2.txt")
	file.SetContentType("text/plain")
	file.SetInfo(M{"any": "thing"})

	c.Assert(file.Id(), Equals, "myid")
	c.Assert(file.Name(), Equals, "myfile2.txt")
	c.Assert(file.ContentType(), Equals, "text/plain")

	err = file.GetInfo(&info)
	c.Assert(err, IsNil)
	c.Assert(info, Equals, bson.M{"any": "thing"})

	err = file.Close()
	c.Assert(err, IsNil)

	c.Assert(file.MD5(), Equals, "1e50210a0202497fb79bc38b6ade6c34")

	result := M{}
	err = db.C("fs.files").Find(nil).One(result)
	c.Assert(err, IsNil)

	result["uploadDate"] = "<timestamp>"

	expected := M{
		"_id":         "myid",
		"length":      9,
		"chunkSize":   262144,
		"uploadDate":  "<timestamp>",
		"md5":         "1e50210a0202497fb79bc38b6ade6c34",
		"filename":    "myfile2.txt",
		"contentType": "text/plain",
		"metadata":    bson.M{"any": "thing"},
	}
	c.Assert(result, Equals, expected)
}

func (s *S) TestGridFSCreateWithChunking(c *C) {
	session, err := mgo.Mongo("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	db := session.DB("mydb")

	gfs := db.GridFS("fs")

	file, err := gfs.Create("")
	c.Assert(err, IsNil)

	file.SetChunkSize(5)

	// Smaller than the chunk size.
	n, err := file.Write([]byte("abc"))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 3)

	// Boundary in the middle.
	n, err = file.Write([]byte("defg"))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 4)

	// Boundary at the end.
	n, err = file.Write([]byte("hij"))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 3)

	// Larger than the chunk size, with 3 chunks.
	n, err = file.Write([]byte("klmnopqrstuv"))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 12)

	err = file.Close()
	c.Assert(err, IsNil)

	// Check the file information.
	result := M{}
	err = db.C("fs.files").Find(nil).One(result)
	c.Assert(err, IsNil)

	fileId, _ := result["_id"].(bson.ObjectId)
	c.Assert(fileId.Valid(), Equals, true)
	result["_id"] = "<id>"
	result["uploadDate"] = "<timestamp>"

	expected := M{
		"_id":        "<id>",
		"length":     22,
		"chunkSize":  5,
		"uploadDate": "<timestamp>",
		"md5":        "44a66044834cbe55040089cabfc102d5",
	}
	c.Assert(result, Equals, expected)

	// Check the chunks.
	iter, err := db.C("fs.chunks").Find(nil).Sort(M{"n": 1}).Iter()
	c.Assert(err, IsNil)

	dataChunks := []string{"abcde", "fghij", "klmno", "pqrst", "uv"}

	for i := 0; ; i++ {
		result = M{}
		err := iter.Next(result)
		if err == mgo.NotFound {
			if i != 5 {
				c.Fatalf("Expected 5 chunks, got %d", i)
			}
			break
		}

		result["_id"] = "<id>"

		expected = M{
			"_id":      "<id>",
			"files_id": fileId,
			"n":        i,
			"data":     []byte(dataChunks[i]),
		}
		c.Assert(result, Equals, expected)
	}
}

func (s *S) TestGridFSOpenNotFound(c *C) {
	session, err := mgo.Mongo("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	db := session.DB("mydb")

	gfs := db.GridFS("fs")
	file, err := gfs.OpenId("non-existent")
	c.Assert(err == mgo.NotFound, Equals, true)
	c.Assert(file, IsNil)

	file, err = gfs.Open("non-existent")
	c.Assert(err == mgo.NotFound, Equals, true)
	c.Assert(file, IsNil)
}

func (s *S) TestGridFSReadAll(c *C) {
	session, err := mgo.Mongo("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	db := session.DB("mydb")

	gfs := db.GridFS("fs")
	file, err := gfs.Create("")
	c.Assert(err, IsNil)
	id := file.Id()

	file.SetChunkSize(5)

	n, err := file.Write([]byte("abcdefghijklmnopqrstuv"))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 22)

	err = file.Close()
	c.Assert(err, IsNil)

	file, err = gfs.OpenId(id)
	c.Assert(err, IsNil)

	b := make([]byte, 30)
	n, err = file.Read(b)
	c.Assert(n, Equals, 22)
	c.Assert(err, IsNil)

	n, err = file.Read(b)
	c.Assert(n, Equals, 0)
	c.Assert(err == os.EOF, Equals, true)

	err = file.Close()
	c.Assert(err, IsNil)
}

func (s *S) TestGridFSReadChunking(c *C) {
	session, err := mgo.Mongo("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	db := session.DB("mydb")

	gfs := db.GridFS("fs")

	file, err := gfs.Create("")
	c.Assert(err, IsNil)

	id := file.Id()

	file.SetChunkSize(5)

	n, err := file.Write([]byte("abcdefghijklmnopqrstuv"))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 22)

	err = file.Close()
	c.Assert(err, IsNil)

	file, err = gfs.OpenId(id)
	c.Assert(err, IsNil)

	b := make([]byte, 30)

	// Smaller than the chunk size.
	n, err = file.Read(b[:3])
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 3)
	c.Assert(b[:3], Equals, []byte("abc"))

	// Boundary in the middle.
	n, err = file.Read(b[:4])
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 4)
	c.Assert(b[:4], Equals, []byte("defg"))

	// Boundary at the end.
	n, err = file.Read(b[:3])
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 3)
	c.Assert(b[:3], Equals, []byte("hij"))

	// Larger than the chunk size, with 3 chunks.
	n, err = file.Read(b)
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 12)
	c.Assert(b[:12], Equals, []byte("klmnopqrstuv"))

	n, err = file.Read(b)
	c.Assert(n, Equals, 0)
	c.Assert(err == os.EOF, Equals, true)

	err = file.Close()
	c.Assert(err, IsNil)
}

func (s *S) TestGridFSOpen(c *C) {
	session, err := mgo.Mongo("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	db := session.DB("mydb")

	gfs := db.GridFS("fs")

	file, err := gfs.Create("myfile.txt")
	c.Assert(err, IsNil)
	file.Write([]byte{'1'})
	file.Close()

	file, err = gfs.Create("myfile.txt")
	c.Assert(err, IsNil)
	file.Write([]byte{'2'})
	file.Close()

	file, err = gfs.Open("myfile.txt")
	c.Assert(err, IsNil)
	defer file.Close()

	var b [1]byte

	_, err = file.Read(b[:])
	c.Assert(err, IsNil)
	c.Assert(string(b[:]), Equals, "2")
}

func (s *S) TestGridFSRemoveId(c *C) {
	session, err := mgo.Mongo("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	db := session.DB("mydb")

	gfs := db.GridFS("fs")

	file, err := gfs.Create("myfile.txt")
	c.Assert(err, IsNil)
	file.Write([]byte{'1'})
	file.Close()

	file, err = gfs.Create("myfile.txt")
	c.Assert(err, IsNil)
	file.Write([]byte{'2'})
	id := file.Id()
	file.Close()

	err = gfs.RemoveId(id)
	c.Assert(err, IsNil)

	file, err = gfs.Open("myfile.txt")
	c.Assert(err, IsNil)
	defer file.Close()

	var b [1]byte

	_, err = file.Read(b[:])
	c.Assert(err, IsNil)
	c.Assert(string(b[:]), Equals, "1")

	n, err := db.C("fs.chunks").Find(M{"files_id": id}).Count()
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 0)
}

func (s *S) TestGridFSRemove(c *C) {
	session, err := mgo.Mongo("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	db := session.DB("mydb")

	gfs := db.GridFS("fs")

	file, err := gfs.Create("myfile.txt")
	c.Assert(err, IsNil)
	file.Write([]byte{'1'})
	file.Close()

	file, err = gfs.Create("myfile.txt")
	c.Assert(err, IsNil)
	file.Write([]byte{'2'})
	file.Close()

	err = gfs.Remove("myfile.txt")
	c.Assert(err, IsNil)

	_, err = gfs.Open("myfile.txt")
	c.Assert(err == mgo.NotFound, Equals, true)

	n, err := db.C("fs.chunks").Find(nil).Count()
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 0)
}
