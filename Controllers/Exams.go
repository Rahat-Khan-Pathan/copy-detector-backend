package Controllers

import (
	"context"
	"encoding/json"
	"errors"

	"example.com/example/DBManager"
	"example.com/example/Models"
	"example.com/example/Utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ExamsNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Exams
	collectionCourses := DBManager.SystemCollections.Courses
	var self Models.Exams
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	if self.CourseRef == primitive.NilObjectID {
		return errors.New("Invalid Course")
	}
	if self.NumberOfQuestionsContest == 0 && self.NumberOfQuestionsWritten == 0 {
		return errors.New("Must Have One Question")
	}
	objID, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	// get course object and add exam if to Exams array
	courseObj, err := CoursesGetByID(self.CourseRef)
	courseObj.Exams = append(courseObj.Exams, objID.InsertedID.(primitive.ObjectID))
	updateData := bson.M{
		"$set": courseObj.GetModifcationBSONObj(),
	}
	_, updateErr := collectionCourses.UpdateOne(context.Background(), bson.M{"_id": courseObj.ID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("An error occurred when modifying course object")
	}
	c.Set("Content-Type", "application/json")
	c.Status(200).Send([]byte("Exam Created Successfully"))
	return nil
}
func ExamsModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Exams
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("Exam Not Found")
	}
	var self Models.Exams
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	if self.CourseRef == primitive.NilObjectID {
		return errors.New("Invalid Course")
	}
	if self.NumberOfQuestionsContest == 0 && self.NumberOfQuestionsWritten == 0 {
		return errors.New("Must Have One Question")
	}
	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("An error occurred when odifying course")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}
func ExamsDelete(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Exams
	collectionCourses := DBManager.SystemCollections.Courses
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	examObj, err := ExamsGetByID(objID)
	if err != nil {
		return errors.New("Exam Not Found")
	}
	courseObj, _ := CoursesGetByID(examObj.CourseRef)
	var newExams []primitive.ObjectID
	for _, val := range courseObj.Exams {
		if val == objID {
			continue
		}
		newExams = append(newExams, val)
	}
	courseObj.Exams = newExams
	updateData := bson.M{
		"$set": courseObj.GetModifcationBSONObj(),
	}
	_, updateErr := collectionCourses.UpdateOne(context.Background(), bson.M{"_id": examObj.CourseRef}, updateData)
	if updateErr != nil {
		return errors.New("There is a problem while updating course. Please try again later")
	}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.Status(500)
		return err
	}
	c.Set("Content-Type", "application/json")
	c.Status(200).Send([]byte("Exam Deleted Successfully"))
	return nil
}
func ExamsGetByID(id primitive.ObjectID) (Models.Exams, error) {
	collection := DBManager.SystemCollections.Exams
	filter := bson.M{"_id": id}
	var self Models.Exams
	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return self, errors.New("object not found")
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	if len(results) == 0 {
		return self, errors.New("object not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}
func ExamsGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Exams
	var self Models.ExamsSearch
	c.BodyParser(&self)
	var results []bson.M
	cur, err := collection.Find(context.Background(), self.GetUsersSearchBSONObj())
	if err != nil {
		c.Status(500)
		return err
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	response, _ := json.Marshal(bson.M{
		"results": results,
	})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}
func ExamsGetByIDPopulated(objID primitive.ObjectID, ptr *Models.Exams) (Models.ExamsPopulated, error) {
	var currentDoc Models.Exams
	if ptr == nil {
		currentDoc, _ = ExamsGetByID(objID)
	} else {
		currentDoc = *ptr
	}
	populatedResult := Models.ExamsPopulated{}
	populatedResult.CloneFrom(currentDoc)
	populatedResult.CourseRef, _ = CoursesGetByID(currentDoc.CourseRef)
	return populatedResult, nil
}
func ExamsGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Exams

	var results []bson.M
	var searchRequests Models.ExamsSearch
	c.BodyParser(&searchRequests)

	b, results := Utils.FindByFilter(collection, searchRequests.GetUsersSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("Object Not Found")
	}

	// Convert
	var allRequestsDocuments []Models.Exams
	byteArr, _ := json.Marshal(results)
	json.Unmarshal(byteArr, &allRequestsDocuments)

	populatetedResults := make([]Models.ExamsPopulated, len(allRequestsDocuments))

	for i, v := range allRequestsDocuments {
		populatetedResults[i], _ = ExamsGetByIDPopulated(v.ID, &v)
	}

	allpopulated, _ := json.Marshal(bson.M{"results": populatetedResults})

	c.Set("Content-Type", "application/json")
	c.Status(200).Send(allpopulated)
	return nil
}
