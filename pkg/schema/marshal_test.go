package schema

import (
	"reflect"
	"testing"
	"time"

	"github.com/s-augustovitko/mongo-schema-go/internal/validation"
)

type testAttachmentType string

type testAudit struct {
	CreatedAt *time.Time `bson:"createdAt" type:"date" validation:"required"`
	CreatedBy string     `bson:"createdBy" validation:"required"`
	UpdatedAt *time.Time `bson:"updatedAt" type:"date"`
	UpdatedBy string     `bson:"updatedBy"`
}

type testTag struct {
	ID      interface{} `bson:"_id"`
	BoardID interface{} `bson:"boardId" validation:"required"`

	Name  string `bson:"name" validation:"required,min=1,max=64"`
	Color string `bson:"color" validation:"required"`

	testAudit `field:",inline" bson:"audit"`
}

type testAttachment struct {
	ID interface{} `bson:"_id"`

	Name string             `bson:"title" binding:"required,min=1,max=128"`
	Type testAttachmentType `bson:"type" binding:"required" enum:"document,image,other"`
	Url  string             `bson:"url" validation:"required,pattern=^(ftp|http|https)://[^ \"]+$"`

	testAudit `field:",inline" bson:"audit"`
}

type marshalTest struct {
	ID       interface{} `bson:"_id"`
	BoardID  interface{} `bson:"boardId" validation:"required"`
	StatusID interface{} `bson:"statusId" validation:"required"`

	Assignee string `bson:"assignee"`
	Reporter string `bson:"reporter"`

	Title   string `bson:"title" validation:"required,min=1,max=256"`
	Content string `bson:"content"`
	Order   int    `bson:"order" validation:"required,min=1"`

	Tag         testTag          `validation:"required"`
	Attachments []testAttachment `bson:"attachments"`
	Test        []string         `bson:"test" validation:"required,min=1" items:"max=20"`
	TagIDs      []interface{}    `bson:"tagIDs"`

	testAudit `field:",inline" bson:"audit"`
}

func TestMarshal(t *testing.T) {
	want := validation.BsonM{"validator": validation.BsonM{"$jsonSchema": validation.BsonM{
		"additionalProperties": true,
		"bsonType":             "object",
		"properties": validation.BsonM{
			"_id":      validation.BsonM{"bsonType": []string{"objectId"}},
			"assignee": validation.BsonM{"bsonType": []string{"string"}},
			"attachments": validation.BsonM{"bsonType": []string{"array"},
				"items": validation.BsonM{
					"bsonType": []string{"object"},
					"properties": validation.BsonM{
						"_id":       validation.BsonM{"bsonType": []string{"objectId"}},
						"createdAt": validation.BsonM{"bsonType": []string{"date"}},
						"createdBy": validation.BsonM{"bsonType": []string{"string"}},
						"title":     validation.BsonM{"bsonType": []string{"string"}},
						"type":      validation.BsonM{"bsonType": []string{"string"}, "enum": []string{"document", "image", "other"}},
						"updatedAt": validation.BsonM{"bsonType": []string{"date"}},
						"updatedBy": validation.BsonM{"bsonType": []string{"string"}},
						"url":       validation.BsonM{"bsonType": []string{"string"}, "pattern": "^(ftp|http|https)://[^ \"]+$"}},
					"required": []string{"url", "createdAt", "createdBy"}},
				"uniqueItems": false},
			"boardId":   validation.BsonM{"bsonType": []string{"objectId"}},
			"content":   validation.BsonM{"bsonType": []string{"string"}},
			"createdAt": validation.BsonM{"bsonType": []string{"date"}},
			"createdBy": validation.BsonM{"bsonType": []string{"string"}},
			"order":     validation.BsonM{"bsonType": []string{"int"}, "minimum": float64(1)},
			"reporter":  validation.BsonM{"bsonType": []string{"string"}},
			"statusId":  validation.BsonM{"bsonType": []string{"objectId"}},
			"tag": validation.BsonM{"bsonType": []string{"object"},
				"properties": validation.BsonM{"_id": validation.BsonM{"bsonType": []string{"objectId"}},
					"boardId":   validation.BsonM{"bsonType": []string{"objectId"}},
					"color":     validation.BsonM{"bsonType": []string{"string"}},
					"createdAt": validation.BsonM{"bsonType": []string{"date"}},
					"createdBy": validation.BsonM{"bsonType": []string{"string"}},
					"name":      validation.BsonM{"bsonType": []string{"string"}, "maxLength": 64, "minLength": 1},
					"updatedAt": validation.BsonM{"bsonType": []string{"date"}},
					"updatedBy": validation.BsonM{"bsonType": []string{"string"}}},
				"required": []string{"boardId", "name", "color", "createdAt", "createdBy"}},
			"tagIDs": validation.BsonM{
				"bsonType":    []string{"array"},
				"items":       validation.BsonM{"bsonType": []string{"objectId"}},
				"uniqueItems": false},
			"test": validation.BsonM{
				"bsonType": []string{"array"},
				"items":    validation.BsonM{"bsonType": []string{"string"}, "maxLength": 20},
				"minItems": 1, "uniqueItems": false},
			"title":     validation.BsonM{"bsonType": []string{"string"}, "maxLength": 256, "minLength": 1},
			"updatedAt": validation.BsonM{"bsonType": []string{"date"}},
			"updatedBy": validation.BsonM{"bsonType": []string{"string"}}},
		"required": []string{"boardId", "statusId", "title", "order", "tag", "test", "createdAt", "createdBy"},
		"title":    "Schema Test"},
	}}

	data := marshalTest{Attachments: []testAttachment{{}}, Test: []string{""}, TagIDs: []interface{}{""}}
	have, warnings, err := Marshal(&data, "Schema Test", true)

	if err != nil || len(warnings) > 0 || !reflect.DeepEqual(want, have) {
		t.Errorf("\nGot: %#v;\nWant: %#v;\nWarns: %#v;\nErr: %#v;", have, want, warnings, err)
	}
}

func TestMarshalErr(t *testing.T) {
	want := validation.BsonM{
		"bsonType":             "object",
		"title":                "Schema Validation",
		"additionalProperties": false,
	}

	data := ""
	have, warnings, err := Marshal(data, "", false)

	if err == nil || len(warnings) > 0 || !reflect.DeepEqual(want, have) {
		t.Errorf("\nGot: %#v;\nWant: %#v;\nWarns: %#v;\nErr: %#v;", have, want, warnings, err)
	}
}
