package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	. "app/dao"
	. "app/models"
)

var dao = PersonsDAO{}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, errorMessage string) {
	errorResponse := struct {
		ErrorMessage string
	}{
		errorMessage,
	}
	respondWithJson(w, code, errorResponse)
}

func respondWithSuccess(w http.ResponseWriter) {
	successResponse := struct {
		Result string
	}{
		Result: "success",
	}
	respondWithJson(w, http.StatusOK, successResponse)
}

// GET list of persons
func AllPersonsEndPoint(w http.ResponseWriter, r *http.Request) {
	persons, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, persons)
}

// GET a person by its ID
func FindPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	person, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Person ID")
		return
	}
	respondWithJson(w, http.StatusOK, person)
}

// POST new person
func CreatePersonEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var person Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	person.ID = bson.NewObjectId()
	if err := dao.Insert(person); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, person)
}

// PUT update a person by id
func UpdatePersonEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	var person Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(params["id"], person); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithSuccess(w)
}

// PATCH update a person by id
func UpdateValueOfPersonEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	person, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Person ID")
		return
	}
	var updatedPerson Person
	if err := json.NewDecoder(r.Body).Decode(&updatedPerson); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if updatedPerson.Order != person.Order {
		respondWithError(w, http.StatusBadRequest, "It's not possible to change the Order")
		return
	}
	if err := dao.Update(params["id"], updatedPerson); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithSuccess(w)
}

//DELETE all persons
func DeleteAllEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if err := dao.DeleteAll(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithSuccess(w)
}

// DELETE a person
func DeletePersonEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	if err := dao.DeleteById(params["id"]); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithSuccess(w)
}

// GET parents by person ID
func FindParentsEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	parents, err := dao.FindParents(params["id"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, parents)
}

// GET children by person ID
func FindChildrenEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	parents, err := dao.FindChildren(params["id"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, parents)
}

// PATCH child and parent to relate one each other
func RelateParentToChildEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	if err := dao.RelateChildToParent(params["parent_id"], params["child_id"]); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithSuccess(w)
}

// PATCH removerelation between child and parent
func RemoveRelation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	if err := dao.RemoveRelation(params["parent_id"], params["child_id"]); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithSuccess(w)
}

// GET genealogical tree of a Person
func TreeEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tree, err := dao.GenTree(params["id"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, tree)
}
