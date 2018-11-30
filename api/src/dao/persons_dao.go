package dao

import (
	"log"

	. "app/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PersonsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "persons"
)

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
