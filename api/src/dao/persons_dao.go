package dao

import (
	"errors"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	txn "gopkg.in/mgo.v2/txn"

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
	err := db.C(COLLECTION).Find(bson.M{"name": bson.M{"$exists": true}}).All(&persons)
	return persons, err
}

// Find a person by its id
func (m *PersonsDAO) FindById(id string) (Person, error) {
	var person Person
	_id, err := stringToObjectId(id)
	if err != nil {
		return person, err
	}
	err = db.C(COLLECTION).FindId(_id).One(&person)
	return person, err
}

// Insert a person into database
func (m *PersonsDAO) Insert(person Person) error {
	err := db.C(COLLECTION).Insert(&person)
	return err
}

// Delete a person by its id
func (m *PersonsDAO) DeleteById(id string) error {
	_id, err := stringToObjectId(id)
	if err != nil {
		return err
	}
	err = db.C(COLLECTION).RemoveId(_id)
	return err
}

//Delete all persons
func (m *PersonsDAO) DeleteAll() error {
	_, err := db.C(COLLECTION).RemoveAll(bson.M{})
	return err
}

// Update a person by its id
func (m *PersonsDAO) Update(id string, person Person) error {
	_id, err := stringToObjectId(id)
	if err != nil {
		return err
	}
	err = db.C(COLLECTION).UpdateId(_id, &person)
	return err
}

// Find parents
func (m *PersonsDAO) FindParents(childId string) ([]Person, error) {
	var parents []Person
	_childId, err := stringToObjectId(childId)
	if err != nil {
		return parents, err
	}
	err = db.C(COLLECTION).Find(bson.M{"children": bson.M{"$in": []bson.ObjectId{_childId}}}).All(&parents)
	return parents, err
}

// Find children
func (m *PersonsDAO) FindChildren(parentId string) ([]Person, error) {
	var children []Person
	_parentId, err := stringToObjectId(parentId)
	if err != nil {
		return children, err
	}
	err = db.C(COLLECTION).Find(bson.M{"parents": bson.M{"$in": []bson.ObjectId{_parentId}}}).All(&children)
	return children, err
}

//Order persons
func (m *PersonsDAO) Order() error {
	var lasts []Person
	order := 1
	err := db.C(COLLECTION).Find(bson.M{"children": bson.M{"$size": 0}}).All(&lasts)
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
func (m *PersonsDAO) addChild(parentId string, childId string) error {
	var parent Person
	parent, err := m.FindById(parentId)
	if err != nil {
		return err
	}
	_childId, err := stringToObjectId(childId)
	if err != nil {
		return err
	}
	if contains(parent.Children, _childId) {
		return errors.New("This relation already exists")
	}
	parent.Children = append(parent.Children, _childId)
	err = m.Update(parentId, parent)
	return err
}

// Add a father to a Person
func (m *PersonsDAO) addParent(childId string, parentId string) error {
	child, err := m.FindById(childId)
	if err != nil {
		return err
	}
	_parentId, err := stringToObjectId(parentId)
	if err != nil {
		return err
	}
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
	if (parent.Order < child.Order) && (len(parent.Children) != 0 || len(parent.Parents) != 0) {
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

// Remove Relation between a child and a parent
func (m *PersonsDAO) RemoveRelation(parentId string, childId string) error {
	_childId, err := stringToObjectId(childId)
	if err != nil {
		return err
	}
	_parentId, err := stringToObjectId(parentId)
	if err != nil {
		return err
	}
	runner := txn.NewRunner(db.C(COLLECTION))
	ops := []txn.Op{{
		C:      COLLECTION,
		Id:     _parentId,
		Update: bson.M{"$pull": bson.M{"children": _childId}},
	}, {
		C:      COLLECTION,
		Id:     _childId,
		Update: bson.M{"$pull": bson.M{"parents": _parentId}},
	}}
	id := bson.NewObjectId() // Optional
	err = runner.Run(ops, id, nil)
	if err != nil {
		return err
	}
	return nil
}

// Find relatives of a person if they have order equal or greater than entry order
func (m *PersonsDAO) FindRelativesWithOrderGreaterThan(id string, order int) ([]Person, error) {
	var relatives, parents, children []Person
	_id, err := stringToObjectId(id)
	if err != nil {
		return relatives, err
	}
	err = db.C(COLLECTION).Find(bson.M{"children": bson.M{"$in": []bson.ObjectId{_id}}, "order": bson.M{"$gte": order}}).All(&parents)
	err = db.C(COLLECTION).Find(bson.M{"parents": bson.M{"$in": []bson.ObjectId{_id}}, "order": bson.M{"$gte": order}}).All(&children)
	if err != nil {
		return nil, err
	}
	relatives = append(relatives, parents...)
	relatives = append(relatives, children...)
	return relatives, nil
}

// Find genealogical tree of Person
func (m *PersonsDAO) GenTree(id string) ([]Person, error) {
	var personsIHave, personsINeed []Person
	target, err := m.FindById(id)
	if err != nil {
		return nil, err
	}
	order := target.Order
	personsINeed = append(personsINeed, target)
	for len(personsIHave) < len(personsINeed) {
		diff := diff(personsIHave, personsINeed)
		personsIHave = personsINeed
		for _, target := range diff {
			relatives, err := m.FindRelativesWithOrderGreaterThan(target.ID.Hex(), order)
			if err != nil {
				return nil, err
			}
			personsINeed = appendUnique(relatives, personsINeed)
		}
	}
	return personsINeed, nil
}
