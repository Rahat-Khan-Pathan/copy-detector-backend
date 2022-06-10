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

func CoursesNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Courses
	var self Models.Courses
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	_, err = collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	c.Set("Content-Type", "application/json")
	c.Status(200).Send([]byte("Course Created Successfully"))
	return nil
}
func CoursesModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Courses
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("Course Not Found")
	}
	var self Models.Courses
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("An error occurred when modifying course object")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}
func CoursesDelete(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Courses
	collectionExams := DBManager.SystemCollections.Exams
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	courseObj, err := CoursesGetByID(objID)
	if err != nil {
		return errors.New("Course Not Found")
	}
	for _, val := range courseObj.Exams {
		_, err = collectionExams.DeleteOne(context.Background(), bson.M{
			"_id": val,
		})
		if err != nil {
			continue
		}
	}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.Status(500)
		return err
	}
	c.Set("Content-Type", "application/json")
	c.Status(200).Send([]byte("Course Deleted Successfully"))
	return nil
}
func CoursesGetByID(id primitive.ObjectID) (Models.Courses, error) {
	collection := DBManager.SystemCollections.Courses
	filter := bson.M{"_id": id}
	var self Models.Courses
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
func CoursesGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Courses
	var self Models.CoursesSearch
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
func CoursesGetByIDPopulated(objID primitive.ObjectID, ptr *Models.Courses) (Models.CoursesPopulated, error) {
	var currentDoc Models.Courses
	if ptr == nil {
		currentDoc, _ = CoursesGetByID(objID)
	} else {
		currentDoc = *ptr
	}
	populatedResult := Models.CoursesPopulated{}
	populatedResult.CloneFrom(currentDoc)
	for _, exam := range currentDoc.Exams {
		examObj, err := ExamsGetByID(exam)
		if err != nil {
			return populatedResult, err
		}
		populatedResult.Exams = append(populatedResult.Exams, examObj)
	}
	return populatedResult, nil
}
func CoursesGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Courses

	var results []bson.M
	var searchRequests Models.CoursesSearch
	c.BodyParser(&searchRequests)

	b, results := Utils.FindByFilter(collection, searchRequests.GetUsersSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("Object Not Found")
	}

	// Convert
	var allRequestsDocuments []Models.Courses
	byteArr, _ := json.Marshal(results)
	json.Unmarshal(byteArr, &allRequestsDocuments)

	populatetedResults := make([]Models.CoursesPopulated, len(allRequestsDocuments))

	for i, v := range allRequestsDocuments {
		populatetedResults[i], _ = CoursesGetByIDPopulated(v.ID, &v)
	}

	allpopulated, _ := json.Marshal(bson.M{"results": populatetedResults})

	c.Set("Content-Type", "application/json")
	c.Status(200).Send(allpopulated)
	return nil
}
