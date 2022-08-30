package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/EleisonC/User-API.git/config"
	"github.com/EleisonC/User-API.git/models"
	"github.com/EleisonC/User-API.git/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// var DogOwner models.DogOwner
var dogOwnerCol *mongo.Collection = config.GetCollection(config.DB, "DogOwner")
var validate = validator.New()

func CreateDogOwner(w http.ResponseWriter, r *http.Request) {
	var dogOwner models.DogOwner
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("I got here")

	err := utils.ParseBody(r, &dogOwner)
	if err != nil {utils.ErrorHandler(w, err, "Failed To Parse The Body")}

	if validationErr := validate.Struct(&dogOwner); validationErr != nil {
		utils.ErrorHandler(w, validationErr, "There was an error validating your data")
		return
	}

	resultDogOwner, err := dogOwnerCol.InsertOne(ctx, dogOwner)
	if err != nil {
		utils.ErrorHandler(w, err, "There was an error entering the data into the DB")
		return
	}

	res, err:=json.Marshal(resultDogOwner)
	if err != nil {
		utils.ErrorHandler(w, err, "There was an error marshalling the data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetAllDogOwners(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dogOwners [] models.DogOwner
	results, err := dogOwnerCol.Find(ctx, bson.M{})
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var dogOwner models.DogOwner
		if err := results.Decode(&dogOwner); err != nil {
			utils.ErrorHandler(w, err, "Error Retriving Data")
			return
		}
		dogOwners = append(dogOwners, dogOwner)
	}

	res, err := json.Marshal(dogOwners)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetDogOwnerById(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dogOwner models.DogOwner
	params := mux.Vars(r)
	dogOwnerId := params["ownerID"]

	objId, err := primitive.ObjectIDFromHex(dogOwnerId)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")
		return
	}

	findErr := dogOwnerCol.FindOne(ctx, bson.M{"_id": objId}).Decode(&dogOwner)
	if findErr != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")

	}

	res, err := json.Marshal(dogOwner)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateDogOwner(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dogOwner models.DogOwner
	params := mux.Vars(r)
	dogOwnerId := params["ownerID"]

	objId, err := primitive.ObjectIDFromHex(dogOwnerId)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Updating Data")
		return
	}

	err = utils.ParseBody(r, &dogOwner)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Updating Data")
		return
	}

	if validationErr := validate.Struct(&dogOwner); validationErr != nil {
		utils.ErrorHandler(w, err, "Error Updating Data")
		return
	}

	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.D{{Key:"$set", Value: bson.D{{Key: "firstname", Value: dogOwner.FirstName},
	{Key: "lastname", Value: dogOwner.LastName}, {Key: "telNumber", Value: dogOwner.TelNumber},
	{Key: "email", Value: dogOwner.Email}, {Key: "password", Value: dogOwner.Password}}}}

	updateResult, err := dogOwnerCol.UpdateOne(ctx, filter, update)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Updating Data")
		return
	}

	if updateResult.MatchedCount == 0 {
		message := "The Owner with ID " + dogOwnerId + " has not been updated or not exist"
		count := updateResult.ModifiedCount
		res := utils.ResMessage{
			Message: message,
			Count: count,
		}
		finalRes, err := json.Marshal(res)
		if err != nil {
			utils.ErrorHandler(w, err, "Error marshalling the data")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(finalRes)
		return
	}
	message := "The dog with ID " + dogOwnerId + " has been updated in the DB"
	count := updateResult.ModifiedCount

	res := utils.ResMessage{
		Message: message,
		Count: count,
	}

	finalRes, err := json.Marshal(res)
	if err != nil {
		utils.ErrorHandler(w, err, "Error marshalling the data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(finalRes)
}

func DeleteDogOwner(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := mux.Vars(r)
	dogOwnerId := params["ownerID"]
	objId, err := primitive.ObjectIDFromHex(dogOwnerId)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Completing Delete")
	}

	delResult, err := dogOwnerCol.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		utils.ErrorHandler(w, err, "Error Deleting User")
		return
	}

	if delResult.DeletedCount == 0 {
		message := "The dog with ID " + dogOwnerId + " has not been deleted from the DB or does not exist"
		count := delResult.DeletedCount
		res := utils.ResMessage{
			Message: message,
			Count: count,
		}
		finalRes, err := json.Marshal(res)
		if err != nil {
			utils.ErrorHandler(w, err, "Error Creating Response")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(finalRes)
		return
	}

	message := "The dog with ID " + dogOwnerId + " has been deleted from the DB"
	delCont := delResult.DeletedCount

	res := utils.ResMessage{
		Message: message,
		Count: delCont,
	}
	finalRes, err := json.Marshal(res)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Creating Response")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(finalRes)
}

