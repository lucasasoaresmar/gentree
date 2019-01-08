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
func (m *PersonsDAO) FindAll() (persons []Person, err error) {
	err = db.C(COLLECTION).Find(bson.M{"name": bson.M{"$exists": true}}).All(&persons)
	return
}

// Find a person by its id
func (m *PersonsDAO) FindById(id string) (person Person, err error) {
	_id, err := stringToObjectId(id)
	if err != nil {
		return
	}
	err = db.C(COLLECTION).FindId(_id).One(&person)
	return
}

// Insert a person into database
func (m *PersonsDAO) Insert(person Person) (err error) {
	err = db.C(COLLECTION).Insert(&person)
	return
}

// Delete a person by its id
func (m *PersonsDAO) DeleteById(id string) (err error) {
	_id, err := stringToObjectId(id)
	if err != nil {
		return
	}
	err = db.C(COLLECTION).RemoveId(_id)
	return
}

//Delete all persons
func (m *PersonsDAO) DeleteAll() (err error) {
	_, err = db.C(COLLECTION).RemoveAll(bson.M{})
	return
}

// Update a person by its id
func (m *PersonsDAO) Update(id string, person Person) (err error) {
	_id, err := stringToObjectId(id)
	if err != nil {
		return
	}
	err = db.C(COLLECTION).UpdateId(_id, &person)
	return
}

// Find parents
func (m *PersonsDAO) FindParents(childId string) (parents []Person, err error) {
	_childId, err := stringToObjectId(childId)
	if err != nil {
		return
	}
	err = db.C(COLLECTION).Find(bson.M{"children": bson.M{"$in": []bson.ObjectId{_childId}}}).All(&parents)
	return
}

// Find children
func (m *PersonsDAO) FindChildren(parentId string) (children []Person, err error) {
	_parentId, err := stringToObjectId(parentId)
	if err != nil {
		return
	}
	err = db.C(COLLECTION).Find(bson.M{"parents": bson.M{"$in": []bson.ObjectId{_parentId}}}).All(&children)
	return
}

//Order persons
func (m *PersonsDAO) Order() (err error) {
	var lasts []Person
	order := 1
	err = db.C(COLLECTION).Find(bson.M{"children": bson.M{"$size": 0}}).All(&lasts)
	if err != nil {
		return
	}
	for len(lasts) > 0 || lasts != nil {
		var tempLasts, _tempLasts []Person
		for _, last := range lasts {
			last.Order = order
			err = m.Update(last.ID.Hex(), last)
			if err != nil {
				return
			}
			if len(last.Parents) == 0 {
				continue
			}
			_tempLasts, err = m.FindParents(last.ID.Hex())
			if err != nil {
				return
			}
			tempLasts = append(tempLasts, _tempLasts...)
		}
		lasts = tempLasts
		order++
	}
	return nil
}

// Add a child to a Person
func (m *PersonsDAO) addChild(parentId string, childId string) (err error) {
	parent, err := m.FindById(parentId)
	if err != nil {
		return
	}
	_childId, err := stringToObjectId(childId)
	if err != nil {
		return
	}
	if contains(parent.Children, _childId) {
		return errors.New("This relation already exists")
	}
	parent.Children = append(parent.Children, _childId)
	err = m.Update(parentId, parent)
	return
}

// Add a father to a Person
func (m *PersonsDAO) addParent(childId string, parentId string) (err error) {
	child, err := m.FindById(childId)
	if err != nil {
		return
	}
	_parentId, err := stringToObjectId(parentId)
	if err != nil {
		return
	}
	if contains(child.Parents, _parentId) {
		return errors.New("This relation already exists")
	}
	child.Parents = append(child.Parents, _parentId)
	err = m.Update(childId, child)
	return
}

// Relate a child to a parent
func (m *PersonsDAO) RelateChildToParent(parentId string, childId string) (err error) {
	if parentId == childId {
		return errors.New("Same ids")
	}
	parent, err := m.FindById(parentId)
	if err != nil {
		return
	}
	child, err := m.FindById(childId)
	if err != nil {
		return
	}
	//Need to improve - check if child is an ancestor
	if (parent.Order < child.Order) && (len(parent.Children) != 0 || len(parent.Parents) != 0) {
		return errors.New("You just can't do that")
	}
	err = m.addChild(parentId, childId)
	if err != nil {
		return
	}
	err = m.addParent(childId, parentId)
	if err != nil {
		return
	}
	err = m.Order()
	return
}

// Remove Relation between a child and a parent
func (m *PersonsDAO) RemoveRelation(parentId string, childId string) (err error) {
	_childId, err := stringToObjectId(childId)
	if err != nil {
		return
	}
	_parentId, err := stringToObjectId(parentId)
	if err != nil {
		return
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
	return
}

// Find relatives of a person if they have order equal or greater than entry order
func (m *PersonsDAO) FindRelativesWithOrderGreaterThan(id string, order int) (relatives []Person, err error) {
	var parents, children []Person
	_id, err := stringToObjectId(id)
	if err != nil {
		return
	}
	err = db.C(COLLECTION).Find(bson.M{"children": bson.M{"$in": []bson.ObjectId{_id}}, "order": bson.M{"$gte": order}}).All(&parents)
	err = db.C(COLLECTION).Find(bson.M{"parents": bson.M{"$in": []bson.ObjectId{_id}}, "order": bson.M{"$gte": order}}).All(&children)
	if err != nil {
		return
	}
	relatives = append(relatives, parents...)
	relatives = append(relatives, children...)
	return
}

// Find genealogical tree of Person
func (m *PersonsDAO) GenTree(id string) (personsINeed []Person, err error) {
	var personsIHave, relatives []Person
	target, err := m.FindById(id)
	if err != nil {
		return
	}
	order := target.Order
	personsINeed = append(personsINeed, target)
	for len(personsIHave) < len(personsINeed) {
		diff := diff(personsIHave, personsINeed)
		personsIHave = personsINeed
		for _, target := range diff {
			relatives, err = m.FindRelativesWithOrderGreaterThan(target.ID.Hex(), order)
			if err != nil {
				return
			}
			personsINeed = appendUnique(relatives, personsINeed)
		}
	}
	return
}
