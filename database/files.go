package main

import (
 "context"
 "fmt"
 "log"
 "regexp"
 "time"

 "go.mongodb.org/mongo-driver/mongo"
 "go.mongodb.org/mongo-driver/mongo/options"
 "go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
 DATABASE_URL    = ""
 DATABASE_NAME   = "Goku"
 MAX_BTN         = 8
)

type Media struct {
 FileUniqueID string
 FileName     string
 FileID       string
 FileSize     int64
}

type Caption struct {
 HTML string
}

type File struct {
 ID       string `bson:"_id"`
 FileName string `bson:"file_name"`
 FileID   string `bson:"file_id"`
 FileSize int64  `bson:"file_size"`
 Type     string `bson:"type"`
 Caption  string `bson:"caption,omitempty"`
}

var client *mongo.Client
var mycol *mongo.Collection

func init() {
 var err error
 client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(DATABASE_URL))
 if err != nil {
  log.Fatal(err)
 }
 
 if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
  log.Fatal(err)
 }
 mycol = client.Database(DATABASE_NAME).Collection("Files")
}

func saveFile(media Media, typ string, l Caption) (string, error) {
 fileName := regexp.MustCompile(`[_\-.+]`).ReplaceAllString(media.FileName, " ")
 file := File{
  ID:       media.FileUniqueID,
  FileName: fileName,
  FileID:   media.FileID,
  FileSize: media.FileSize,
  Type:     typ,
  Caption:  l.HTML,
 }

 _, err := mycol.InsertOne(context.TODO(), file)
 if err != nil {
  if mongo.IsDuplicateKeyError(err) {
   fmt.Printf("%s is already saved\n", fileName)
   return "Duplicates", nil
  }
  return "", err
 }

 fmt.Printf("%s is saved\n", fileName)
 return "Saved", nil
}

func getDatabaseSize() (int64, error) {
 var stats struct {
  DataSize int64 `bson:"dataSize"`
 }
 err := mycol.Database().RunCommand(context.TODO(), map[string]int{"dbstats": 1}).Decode(&stats)
 if err != nil {
  return 0, err
 }
 return stats.DataSize, nil
}

func totalFilesCount() (int64, error) {
 count, err := mycol.CountDocuments(context.TODO(), bson.M{})
 if err != nil {
  return 0, err
 }
 return count, nil
}

func getSearchResults(query string, maxResults int64, offset int64, filter bool, language string) ([]File, interface{}, int64, error) {
 query = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(query, "")

 var rawPattern string
 if filter {
  query = regexp.MustCompile(`\s`).ReplaceAllString(query, `(\s|\.|\+|-|_)`)
  rawPattern = `(\s|_|\-|\.|\+)` + query + `(\s|_|\-|\.|\+)`
 } else if query == "" {
  rawPattern = `.` // Match anything
 } else if !strings.Contains(query, " ") {
  rawPattern = `(\b|[\.\+\-_])` + query + `(\b|[\.\+\-_])`
 } else {
  rawPattern = regexp.QuoteMeta(query)
  rawPattern = strings.ReplaceAll(rawPattern, " ", ".*[\\s\\.\\+\\-\\_]") 
 }

 regex, err := regexp.Compile(rawPattern)
 if err != nil {
  return nil, nil, 0, err
 }

 filterCriteria := bson.M{"file_name": bson.M{"$regex": regex}}
 totalResults, err := mycol.CountDocuments(context.TODO(), filterCriteria)
 if err != nil {
  return nil, nil, 0, err
 }

 nextOffset := offset + maxResults
 if nextOffset >= totalResults {
  nextOffset = 0 // Resetting if exceeds total
 }

 cursor, err := mycol.Find(context.TODO(), filterCriteria, options.Find().SetSkip(offset).SetLimit(maxResults))
 if err != nil {
  return nil, nil, 0, err
 }
 defer cursor.Close(context.TODO())

 var files []File
 if err := cursor.All(context.TODO(), &files); err != nil {
  return nil, nil, 0, err
 }

 return files, nextOffset, totalResults, nil
}
