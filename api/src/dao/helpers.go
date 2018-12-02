package dao

import (
	"gopkg.in/mgo.v2/bson"

	. "app/models"
)

func contains(s []bson.ObjectId, e bson.ObjectId) bool {
  for _, a := range s {
      if a == e {
          return true
      }
  }
  return false
}

func diff(a []Person, b []Person) []Person {
	var difference []Person
	for _, _b := range b {
		equal := false
		for _, _a := range a {
			if _a.ID == _b.ID {
				equal = true
			}
		}
		if !equal { difference = append(difference, _b) }
	}
	return difference
}

func appendUnique(a []Person, b []Person) []Person {
	difference := diff(a, b)
	return append(a, difference...)
}

func removeId(slice []bson.ObjectId, id bson.ObjectId) []bson.ObjectId {
	var newSlice []bson.ObjectId
	for _, _id := range slice {
		if _id != id {
			newSlice = append(newSlice, _id)
		}
	}
	return newSlice
}
