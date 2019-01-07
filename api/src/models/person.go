package models

import "gopkg.in/mgo.v2/bson"

type Person struct {
	ID       bson.ObjectId   `bson:"_id" json:"id"`
	Name     string          `bson:"name" json:"name,omitempty"`
	Order    int             `bson:"order" json:"order"`
	Parents  []bson.ObjectId `bson:"parents" json:"parents,omitempty"`
	Children []bson.ObjectId `bson:"children" json:"-"`
}
