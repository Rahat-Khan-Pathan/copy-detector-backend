package Utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindByFilter(collection *mongo.Collection, filter bson.M) (bool, []bson.M) {
	results := []bson.M{}

	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return false, results
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)

	return true, results
}

func ArrayStringContains(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func UploadImageBase64(stringBase64 string, imageDocType string) (string, error) {
	i := strings.Index(stringBase64, ",")
	if i != -1 {
		file, _ := base64.StdEncoding.DecodeString(stringBase64[i+1:])
		var filePath = fmt.Sprintf("Resources/Images/client_att_%d_%d.%s", rand.Intn(1024), MakeTimestamp(), imageDocType)

		f, err := os.Create("./" + filePath)
		if err != nil {
			return "", err
		}
		defer f.Close()

		if _, err := f.Write(file); err != nil {
			return "", err
		}
		f.Sync()
		return filePath, nil
	}
	return "", nil
}
