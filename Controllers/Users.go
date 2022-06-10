package Controllers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"example.com/example/DBManager"
	"example.com/example/Models"
	"example.com/example/Utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UsersValidateUsersUsername(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Users
	username := c.Params("username")
	filter := bson.M{
		"username": username,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) > 0 {
		// Decode
		response, _ := json.Marshal(
			bson.M{
				"results": bson.M{
					"has": true,
				},
			},
		)
		c.Set("Content-Type", "application/json")
		c.Status(200).Send(response)
		return nil
	} else {
		// Decode
		response, _ := json.Marshal(
			bson.M{
				"results": bson.M{
					"has": false,
				},
			},
		)
		c.Set("Content-Type", "application/json")
		c.Status(200).Send(response)
		return nil
	}
}

func UsersValidateUsersUsernameFunction(username string) (bool, interface{}) {
	collection := DBManager.SystemCollections.Users
	filter := bson.M{
		"username": username,
	}
	b, results := Utils.FindByFilter(collection, filter)
	id := ""
	if len(results) > 0 {
		id = results[0]["_id"].(primitive.ObjectID).Hex()
	}
	return b, id
}

func UsersValidateUsersEmail(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Users
	email := c.Params("email")
	filter := bson.M{
		"email": email,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) > 0 {
		// Decode
		response, _ := json.Marshal(
			bson.M{
				"results": bson.M{
					"has": true,
				},
			},
		)
		c.Set("Content-Type", "application/json")
		c.Status(200).Send(response)
		return nil
	} else {
		// Decode
		response, _ := json.Marshal(
			bson.M{
				"results": bson.M{
					"has": false,
				},
			},
		)
		c.Set("Content-Type", "application/json")
		c.Status(200).Send(response)
		return nil
	}
}

func UsersValidateUsersEmailFunction(email string) (bool, interface{}) {
	collection := DBManager.SystemCollections.Users
	filter := bson.M{
		"email": email,
	}
	b, results := Utils.FindByFilter(collection, filter)
	id := ""
	if len(results) > 0 {
		id = results[0]["_id"].(primitive.ObjectID).Hex()
	}
	return b, id
}

func UsersGetByIDFunction(id primitive.ObjectID) (Models.Users, error) {
	collection := DBManager.SystemCollections.Users
	filter := bson.M{"_id": id}
	var self Models.Users
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
func UsersSignUp(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Users
	var self Models.Users
	c.BodyParser(&self)
	err := self.ValidateSignUp()
	if err != nil {
		return err
	}
	_, existing := UsersValidateUsersUsernameFunction(self.Username)
	if existing != "" {
		return errors.New("Username Already Exists")
	}
	_, existing = UsersValidateUsersEmailFunction(self.Email)
	if existing != "" {
		return errors.New("This Email Is Already Registered")
	}
	// converting the password to sha256 hash
	h := sha256.New()
	h.Write([]byte(self.Password))
	sha256_hash := hex.EncodeToString(h.Sum(nil))
	self.PasswordHash = sha256_hash
	self.Password = ""
	self.Registered = false
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	id := res.InsertedID.(primitive.ObjectID)
	usersObj, _ := UsersGetByIDFunction(id)
	usersObj.PasswordHash = ""
	// Decode
	response, _ := json.Marshal(
		bson.M{
			"results": usersObj,
		},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func UsersRegister(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Users
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	var results []Models.Users
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return errors.New("User Not Found")
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("User Not Found")
	}
	results[0].Registered = true
	updateData := bson.M{
		"$set": results[0].GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": results[0].ID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("An Error Occurred While Registering. Please Try Again Later")
	} else {
		usersObj, _ := UsersGetByIDFunction(results[0].ID)
		usersObj.PasswordHash = ""
		// Decode
		response, _ := json.Marshal(
			bson.M{
				"results": usersObj,
			},
		)
		c.Set("Content-Type", "application/json")
		c.Status(200).Send(response)
		return nil
	}
}
func UsersReject(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Users
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	whoID, _ := primitive.ObjectIDFromHex(c.Params("who"))
	// who
	filter := bson.M{
		"_id": whoID,
	}
	var results1 []Models.Users
	cur, _ := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results1)
	if len(results1) == 0 {
		c.Status(404)
		return errors.New("Admin Not Found")
	}
	if results1[0].Admin == false {
		return errors.New("You are not admin. You can't reject other users")
	}

	// reject user
	filter = bson.M{
		"_id": objID,
	}
	var results []Models.Users
	cur, _ = collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("User Not Found")
	}
	results[0].Registered = false
	updateData := bson.M{
		"$set": results[0].GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": results[0].ID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("An Error Occurred While Rejecting. Please Try Again Later")
	} else {
		usersObj, _ := UsersGetByIDFunction(results[0].ID)
		usersObj.PasswordHash = ""
		// Decode
		response, _ := json.Marshal(
			bson.M{
				"results": usersObj,
			},
		)
		c.Set("Content-Type", "application/json")
		c.Status(200).Send(response)
		return nil
	}
}
func UsersPasswordChange(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Users
	var self Models.UsersPasswordChange
	c.BodyParser(&self)
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	userObj, err := UsersGetByIDFunction(id)
	if err != nil {
		return errors.New("User Not Found")
	}
	// converting the password to sha256 hash
	h := sha256.New()
	h.Write([]byte(self.CurrentPassword))
	sha256_hash := hex.EncodeToString(h.Sum(nil))
	if sha256_hash != userObj.PasswordHash {
		return errors.New("Current Password is Incorrect")
	}
	h.Write([]byte(self.NewPassword))
	sha256_hash = hex.EncodeToString(h.Sum(nil))
	userObj.PasswordHash = sha256_hash
	updateData := bson.M{
		"$set": userObj.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": id}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("An Error Occurred While Registering. Please Try Again Later")
	}
	c.Set("Content-Type", "application/json")
	c.Status(200).Send([]byte("Password Changed Successfully"))
	return nil
}
func UsersLogin(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Users
	var self Models.UsersLogin
	c.BodyParser(&self)
	err := self.ValidateLogin()
	if err != nil {
		return err
	}
	// converting the password to sha256 hash
	h := sha256.New()
	h.Write([]byte(self.Password))
	sha256_hash := hex.EncodeToString(h.Sum(nil))

	filter := bson.M{
		"email":        self.EmailOrUsername,
		"passwordhash": sha256_hash,
	}
	filter2 := bson.M{
		"username":     self.EmailOrUsername,
		"passwordhash": sha256_hash,
	}
	var results []Models.Users
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.Status(500)
		return err
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	if len(results) == 1 {
		// Decode
		response, _ := json.Marshal(
			bson.M{
				"results": results[0],
			},
		)
		c.Set("Content-Type", "application/json")
		c.Status(200).Send(response)
		return nil
	}
	cur, err = collection.Find(context.Background(), filter2)
	if err != nil {
		c.Status(500)
		return err
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	if len(results) == 1 {
		// Decode
		response, _ := json.Marshal(
			bson.M{
				"results": results[0],
			},
		)
		c.Set("Content-Type", "application/json")
		c.Status(200).Send(response)
		return nil
	}
	c.Status(404)
	return errors.New("Incorrect Username/Email or Password")
}

func UsersGetByID(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Users
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{"_id": id}
	var results []bson.M
	var self Models.Users
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	if len(results) == 0 {
		c.Status(500)
		return errors.New("User Not Found")
	}
	// Decode
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	self.PasswordHash = ""
	response, _ := json.Marshal(
		bson.M{
			"results": self,
		},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}
func UsersGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Users
	var self Models.UsersSearch
	c.BodyParser(&self)
	var results []Models.Users
	cur, err := collection.Find(context.Background(), self.GetUsersSearchBSONObj())
	if err != nil {
		err := errors.New("Something Went Wrong. Please Try Again Later")
		c.Status(500).Send([]byte(err.Error()))
		return err
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)

	// Decode
	bsonBytes, _ := bson.Marshal(results) // Decode
	bson.Unmarshal(bsonBytes, &results)   // Encode
	for i, _ := range results {
		results[i].PasswordHash = ""
	}
	response, _ := json.Marshal(
		bson.M{"results": results},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}
