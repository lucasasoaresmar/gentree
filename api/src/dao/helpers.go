package dao

import (
	"gopkg.in/mgo.v2/bson"
)

func contains(s []bson.ObjectId, e bson.ObjectId) bool {
  for _, a := range s {
      if a == e {
          return true
      }
  }
  return false
}
