package Controllers

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	"example.com/example/DBManager"
	"example.com/example/Models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const questionTypeContest = "contest"
const questionTypeWritten = "written"

func CodesSubmit(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Codes
	var self Models.Codes
	c.BodyParser(&self)
	examObj, _ := ExamsGetByID(self.ExamRef)
	if self.CourseRef == primitive.NilObjectID {
		return errors.New("Invalid Course")
	}
	if self.ExamRef == primitive.NilObjectID {
		return errors.New("Invalid Exam")
	}
	if self.QuestionType == questionTypeContest {
		if self.QuestionNo < 0 || self.QuestionNo > examObj.NumberOfQuestionsContest {
			return errors.New("Invalid Contest Question No")
		}
	}
	if self.QuestionType == questionTypeWritten {
		if self.QuestionNo < 0 || self.QuestionNo > examObj.NumberOfQuestionsWritten {
			return errors.New("Invalid Written Question No")
		}
	}
	if self.Email == "" {
		return errors.New("Invalid email")
	}

	var results1 []Models.Codes
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"courseref":    self.CourseRef,
				"examref":      self.ExamRef,
				"questionno":   self.QuestionNo,
				"questiontype": self.QuestionType,
				"email":        self.Email,
			},
		},
	}
	cur, _ := collection.Aggregate(context.Background(), pipeline)
	cur.All(context.Background(), &results1)
	if len(results1) > 0 {
		_, err := collection.DeleteOne(context.Background(), bson.M{"_id": results1[0].ID})
		if err != nil {
			c.Status(500)
			return err
		}
	}
	var results []Models.Codes
	pipeline = []bson.M{
		{
			"$match": bson.M{
				"courseref":  self.CourseRef,
				"examref":    self.ExamRef,
				"questionno": self.QuestionNo,
			},
		},
	}
	cur, _ = collection.Aggregate(context.Background(), pipeline)
	cur.All(context.Background(), &results)
	finalResults := []Models.CodesCopy{}
	for _, val := range results {
		first_code := val.Code
		second_code := self.Code
		total_character := 0
		total_copied := 0.0
		singleResult := Models.CodesCopy{}
		for i := 0; i < 255; i++ {
			character := string(i)
			if i == 10 || i == 32 {
				continue
			}
			first := strings.Count(first_code, character)
			second := strings.Count(second_code, character)
			min := math.Min(float64(first), float64(second))
			max := math.Max(float64(first), float64(second))
			if max == 0 {
				continue
			}
			total_character++
			total_copied = total_copied + (min / max)
		}
		singleResult.Email = val.Email
		singleResult.Percentage = (total_copied * 100) / float64(total_character)
		singleResult.SubmitDate = val.SubmitDate
		singleResult.Code = val.Code
		finalResults = append(finalResults, singleResult)
	}
	self.SubmitDate = primitive.NewDateTimeFromTime(time.Now())
	_, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	sort.SliceStable(finalResults, func(i, j int) bool {
		return finalResults[i].Percentage > finalResults[j].Percentage
	})
	response, _ := json.Marshal(bson.M{
		"results": finalResults,
	})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func CodesGetResults(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Codes
	var self Models.Codes
	c.BodyParser(&self)
	examObj, _ := ExamsGetByID(self.ExamRef)
	if self.CourseRef == primitive.NilObjectID {
		return errors.New("Invalid course")
	}
	if self.ExamRef == primitive.NilObjectID {
		return errors.New("Invalid exam")
	}
	if self.QuestionType == questionTypeContest {
		if self.QuestionNo < 0 || self.QuestionNo > examObj.NumberOfQuestionsContest {
			return errors.New("Invalid Contest Question No")
		}
	}
	if self.QuestionType == questionTypeWritten {
		if self.QuestionNo < 0 || self.QuestionNo > examObj.NumberOfQuestionsWritten {
			return errors.New("Invalid Written Question No")
		}
	}
	if self.Email == "" {
		return errors.New("Invalid email")
	}

	var results1 []Models.Codes
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"courseref":    self.CourseRef,
				"examref":      self.ExamRef,
				"questionno":   self.QuestionNo,
				"questiontype": self.QuestionType,
				"email":        self.Email,
			},
		},
	}
	cur, _ := collection.Aggregate(context.Background(), pipeline)
	cur.All(context.Background(), &results1)
	if len(results1) == 0 {

		return errors.New("Code not found for this email")
	}
	self.Code = results1[0].Code
	var results []Models.Codes
	pipeline = []bson.M{
		{
			"$match": bson.M{
				"courseref":  self.CourseRef,
				"examref":    self.ExamRef,
				"questionno": self.QuestionNo,
			},
		},
	}
	cur, _ = collection.Aggregate(context.Background(), pipeline)
	cur.All(context.Background(), &results)
	finalResults := []Models.CodesCopy{}
	for _, val := range results {
		if val.ID == results1[0].ID {
			continue
		}
		first_code := val.Code
		second_code := self.Code
		total_character := 0
		total_copied := 0.0
		singleResult := Models.CodesCopy{}
		for i := 0; i < 255; i++ {
			character := string(i)
			first := strings.Count(first_code, character)
			second := strings.Count(second_code, character)
			min := math.Min(float64(first), float64(second))
			max := math.Max(float64(first), float64(second))
			if max == 0 {
				continue
			}
			total_character++
			total_copied = total_copied + (min / max)
		}
		singleResult.Email = val.Email
		singleResult.Percentage = (total_copied * 100) / float64(total_character)
		singleResult.SubmitDate = val.SubmitDate
		singleResult.Code = val.Code
		finalResults = append(finalResults, singleResult)
	}
	sort.SliceStable(finalResults, func(i, j int) bool {
		return finalResults[i].Percentage > finalResults[j].Percentage
	})
	response, _ := json.Marshal(bson.M{
		"results": finalResults,
		"code":    self.Code,
	})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}
