package dao

import (
	"log"
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	
	. "app/models"
)

type PersonsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "persons"
)

/*
	Helper functions
*/
// Check if id is in a slice
func contains(s []bson.ObjectId, e bson.ObjectId) bool {
  for _, a := range s {
      if a == e {
          return true
      }
  }
  return false
}

// Establish a connection to database
func (m *PersonsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of persons
func (m *PersonsDAO) FindAll() ([]Person, error) {
	var persons []Person
	err := db.C(COLLECTION).Find(bson.M{}).All(&persons)
	return persons, err
}

// Find a person by its id
func (m *PersonsDAO) FindById(id string) (Person, error) {
	var person Person
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&person)
	return person, err
}

// Insert a person into database
func (m *PersonsDAO) Insert(person Person) error {
	err := db.C(COLLECTION).Insert(&person)
	return err
}

// Delete a person by its id
func (m *PersonsDAO) DeleteById(id string) error {
	err := db.C(COLLECTION).RemoveId(bson.ObjectIdHex(id))
	return err
}

//Delete all persons
func (m *PersonsDAO) DeleteAll() error {
	_, err := db.C(COLLECTION).RemoveAll(bson.M{})
	return err
}

// Update a person by its id
func (m *PersonsDAO) Update(id string, person Person) error {
	err := db.C(COLLECTION).UpdateId(bson.ObjectIdHex(id), &person)
	return err
}

// Find parents
func (m *PersonsDAO) FindParents(childId string) ([]Person, error) {
	var parents []Person
	err := db.C(COLLECTION).Find(bson.M{"children": bson.M{"$in": []bson.ObjectId{bson.ObjectIdHex(childId)}}}).All(&parents)
	return parents, err
}

// Find children
func (m *PersonsDAO) FindChildren(parentId string) ([]Person, error) {
	var children []Person
	err := db.C(COLLECTION).Find(bson.M{"parents": bson.M{"$in": []bson.ObjectId{bson.ObjectIdHex(parentId)}}}).All(&children)
	return children, err
}

//Order persons
func (m *PersonsDAO) Order() error {
	var lasts []Person
	order := 1

	err := db.C(COLLECTION).Find(bson.M{"children": bson.M{"$size":0}}).All(&lasts)
	if err != nil {
		return err
	}

	for len(lasts) > 0 || lasts != nil {
		var tempLasts []Person
		for _, last := range lasts {
			last.Order = order
			err := m.Update(last.ID.Hex(), last)
			if err != nil {
				return err
			}
			if len(last.Parents) == 0 {
				continue
			}
			_tempLasts, err := m.FindParents(last.ID.Hex())
			if err != nil {
				return err
			}
			tempLasts = append(tempLasts, _tempLasts...)
		}
		lasts = tempLasts
		order++
	}
	return nil
}

// Add a child to a Person
func (m *PersonsDAO) addChild (parentId string, childId string) error {
	var parent Person
	parent, err := m.FindById(parentId)
	if err != nil { 
		return err
	}
	_childId := bson.ObjectIdHex(childId)
	if contains(parent.Children, _childId) { 
		return errors.New("This relation already exists")
	}
	parent.Children = append(parent.Children, _childId)
	err = m.Update(parentId, parent)
	return err
}

// Add a father to a Person
func (m *PersonsDAO) addParent (childId string, parentId string) error {
	child, err := m.FindById(childId)
	if err != nil {
		return err
	}
	_parentId := bson.ObjectIdHex(parentId)
	if contains(child.Parents, _parentId) { 
		return errors.New("This relation already exists")
	}
	child.Parents = append(child.Parents, _parentId)
	err = m.Update(childId, child)
	return err
}

// Relate a child to a parent
func (m *PersonsDAO) RelateChildToParent(parentId string, childId string) error {
	if parentId == childId {
		return errors.New("Same ids")
	}
	parent, err := m.FindById(parentId)
	if err != nil {
		return err
	}
	child, err := m.FindById(childId)
	if err != nil {
		return err
	}
	//Need to improve - check if child is an ancestor
	if parent.Order < child.Order && parent.Order != 0 {
		return errors.New("You just can't do that")
	}
	err = m.addChild(parentId, childId)
	if err != nil { 
		return err
	}
	err = m.addParent(childId, parentId)
	if err != nil { 
		return err
	}
	err = m.Order()
	return err
}
