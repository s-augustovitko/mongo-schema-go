package validation

import (
	"reflect"
	"testing"
	"time"

	"github.com/s-augustovitko/mongo-schema-go/internal/tags"
)

type createConfigTestItem struct {
	Arg2 string    `validation:"pattern=@gmail.com$,patternProperties=gi,multipleOf=2.3"`
	Arg3 int       `validation:"required,multipleOf=2.3" type:"int,long"`
	Arg4 int64     `validation:"required,min=2.1,pattern=@gmail.com$,patternProperties=gi"`
	Arg5 float32   `validation:"multipleOf=2.3"`
	Arg6 float64   `validation:"required"`
	Arg7 bool      `validation:"required"`
	Arg8 []*string `enum:"a,b,c,d" validation:"uniqueItems,required" items:"min=3,max=7" itemsType:"string"`
}

type createConfigTest struct {
	Arg1                 string      `enum:"a,b , c,d" validation:"min=2,max=50"`
	Id                   interface{} `bson:"_id" validation:"required,max=5"`
	createConfigTestItem `field:",inline" validation:"required"`
	Date                 *time.Time             `type:"date" validation:"required"`
	Arr                  []createConfigTestItem `validation:"min=1,max=5,required"`
	Obj                  createConfigTestItem   `field:"obj1" bson:"obj2" validation:"required"`
	M                    BsonM                  `description:"some cool description" validation:"required, min=1"`
}

func TestCreateConfig(t *testing.T) {
	str := ""
	obj := createConfigTest{
		createConfigTestItem: createConfigTestItem{Arg8: []*string{&str}},
		Obj:                  createConfigTestItem{Arg8: []*string{&str}},
		Arr:                  []createConfigTestItem{{Arg8: []*string{&str}}},
	}
	wantArr := []config{
		{
			Validation: Validation{Min: CreateVal[float64](2), Max: CreateVal[float64](50)},
			Tag:        "arg1",
			BsonType:   []string{"string"},
			Enum:       CreateVal([]string{"a", "b", "c", "d"}),
		},
		{
			Validation: Validation{Required: true, Max: CreateVal[float64](5)},
			Tag:        "_id",
			BsonType:   []string{"objectId"},
			Enum:       WithVal[[]string]{Val: nil},
		},
		{
			Validation: Validation{},
			Tag:        "createConfigTestItem",
			BsonType:   []string{"object"},
			IsInline:   true,
			IsStruct:   true,
			Enum:       WithVal[[]string]{Val: nil},
		},
		{
			Validation: Validation{Required: true},
			Tag:        "date",
			BsonType:   []string{"date"},
			Enum:       WithVal[[]string]{Val: nil},
		},
		{
			Validation:      Validation{Required: true, Min: CreateVal[float64](1), Max: CreateVal[float64](5)},
			IsArray:         true,
			IsArrayOfStruct: true,
			Tag:             "arr",
			BsonType:        []string{"array"},
			Enum:            WithVal[[]string]{Val: nil},
		},
		{
			Validation: Validation{Required: true},
			Tag:        "obj1",
			BsonType:   []string{"object"},
			IsStruct:   true,
			Enum:       WithVal[[]string]{Val: nil},
		},
		{
			Validation:  Validation{Required: true, Min: CreateVal[float64](1)},
			BsonType:    []string{"object"},
			Tag:         "m",
			Description: CreateVal("some cool description"),
			Enum:        WithVal[[]string]{Val: nil},
		},
	}

	val := reflect.ValueOf(obj)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		val := val.Field(i)
		field := typ.Field(i)
		want := wantArr[i]

		have, err := createConfig(val, field)

		if !reflect.DeepEqual(have, want) || err != nil {
			t.Errorf("Field:%#v;\nGot: %#v;\nWant: %#v;\nErr: %#v", have.Tag, have, want, err)
		}
	}
}

func TestCreateJSONSchema(t *testing.T) {
	str := ""
	obj := createConfigTest{
		createConfigTestItem: createConfigTestItem{Arg8: []*string{&str}},
		Obj:                  createConfigTestItem{Arg8: []*string{&str}},
		Arr:                  []createConfigTestItem{{Arg8: []*string{&str}}},
	}
	want := BsonM{
		"_id":  BsonM{"bsonType": []string{"objectId"}},
		"arg1": BsonM{"bsonType": []string{"string"}, "enum": []string{"a", "b", "c", "d"}, "maxLength": 50, "minLength": 2},
		"arg2": BsonM{"bsonType": []string{"string"}, "pattern": "@gmail.com$", "patternProperties": "gi"},
		"arg3": BsonM{"bsonType": []string{"int", "long"}, "multipleOf": 2.3},
		"arg4": BsonM{"bsonType": []string{"long"}, "minimum": 2.1},
		"arg5": BsonM{"bsonType": []string{"double"}, "multipleOf": 2.3},
		"arg6": BsonM{"bsonType": []string{"decimal"}},
		"arg7": BsonM{"bsonType": []string{"bool"}},
		"arg8": BsonM{
			"bsonType": []string{"array"},
			"items": BsonM{
				"bsonType":  []string{"string"},
				"enum":      []string{"a", "b", "c", "d"},
				"maxLength": 7, "minLength": 3},
			"uniqueItems": true},
		"arr": BsonM{
			"bsonType": []string{"array"},
			"items": BsonM{"bsonType": []string{"object"},
				"properties": BsonM{
					"arg2": BsonM{"bsonType": []string{"string"}, "pattern": "@gmail.com$", "patternProperties": "gi"},
					"arg3": BsonM{"bsonType": []string{"int", "long"}, "multipleOf": 2.3},
					"arg4": BsonM{"bsonType": []string{"long"}, "minimum": 2.1},
					"arg5": BsonM{"bsonType": []string{"double"}, "multipleOf": 2.3},
					"arg6": BsonM{"bsonType": []string{"decimal"}},
					"arg7": BsonM{"bsonType": []string{"bool"}},
					"arg8": BsonM{"bsonType": []string{"array"},
						"items": BsonM{
							"bsonType":  []string{"string"},
							"enum":      []string{"a", "b", "c", "d"},
							"maxLength": 7, "minLength": 3},
						"uniqueItems": true}},
				"required": []string{"arg3", "arg4", "arg6", "arg7", "arg8"}},
			"maxItems": 5, "minItems": 1,
			"uniqueItems": false},
		"date": BsonM{"bsonType": []string{"date"}},
		"m": BsonM{
			"bsonType":      []string{"object"},
			"description":   "some cool description",
			"minProperties": 1},
		"obj1": BsonM{
			"bsonType": []string{"object"},
			"properties": BsonM{
				"arg2": BsonM{"bsonType": []string{"string"}, "pattern": "@gmail.com$", "patternProperties": "gi"},
				"arg3": BsonM{"bsonType": []string{"int", "long"}, "multipleOf": 2.3},
				"arg4": BsonM{"bsonType": []string{"long"}, "minimum": 2.1},
				"arg5": BsonM{"bsonType": []string{"double"}, "multipleOf": 2.3},
				"arg6": BsonM{"bsonType": []string{"decimal"}},
				"arg7": BsonM{"bsonType": []string{"bool"}},
				"arg8": BsonM{
					"bsonType": []string{"array"},
					"items": BsonM{
						"bsonType":  []string{"string"},
						"enum":      []string{"a", "b", "c", "d"},
						"maxLength": 7, "minLength": 3},
					"uniqueItems": true}},
			"required": []string{"arg3", "arg4", "arg6", "arg7", "arg8"}}}
	wantReq := []string{"_id", "arg3", "arg4", "arg6", "arg7", "arg8", "date", "arr", "obj1", "m"}

	have := BsonM{}
	required, errs := CreateJSONSchema(reflect.ValueOf(obj), &have)

	if !reflect.DeepEqual(want, have) {
		t.Errorf("Field: Object;\nGot: %#v;\nWant: %#v;", have, want)
	}
	if !tags.CompareArr(required, wantReq) || len(errs) > 0 {
		t.Errorf("Field: Required;\nGot: %#v;\nWant: %#v;\nErrs: %#v", required, wantReq, errs)
	}
}

type createJSONSchemaTestErrs struct {
	Invalid          string   `validation:"invalid"`
	InvalidType      string   `type:"invalid"`
	InvalidItems     []string `items:"invalid"`
	InvalidItemsType []string `itemsType:"invalid"`
	InvalidEmpty     []string
	InvalidLen0      []string
}

func TestCreateJSONSchemaErrs(t *testing.T) {
	obj := createJSONSchemaTestErrs{
		InvalidItems:     []string{""},
		InvalidItemsType: []string{""},
		InvalidLen0:      []string{},
	}
	want := BsonM{}

	have := BsonM{}
	reqs, errs := CreateJSONSchema(reflect.ValueOf(obj), &have)
	mustHaves := []string{
		"invalid",
		"invalidType",
		"invalidItems",
		"invalidItemsType",
		"invalidEmpty",
		"invalidLen0",
	}

	errsM := map[string]bool{}
	for _, err := range errs {
		errsM[err.(ErrorWithTag).Tag()] = true
	}

	for _, item := range mustHaves {
		if !errsM[item] {
			t.Errorf("Error: %v not found;\nErrs Map: %v\nErrs: %#v;", item, errsM, errs)
		}
	}

	if len(reqs) > 0 || !reflect.DeepEqual(have, want) {
		t.Errorf("Reqs: %#v;\nGot: %#v;\nWant: %#v;", reqs, have, want)
	}
}
