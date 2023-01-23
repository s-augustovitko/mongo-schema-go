# Mongo Schema Go

This package is used to create a `MongoDB` schemas based on a golang struct, the most basic usage can be seen in the following example.

```go
package main

import (
    "log"

    "github.com/s-augustovitko/mongo-schema-go/pkg/schema"
)

type Obj struct {
    ID          interface{}
    Name        string
    Email       string
    Address     string
    Age         int
    Birthday    date `type:"date"`
}

func main() {
    obj := Obj{}
    out, warnings, err := schema.Marshal(obj, "Demo Object Schema", true)
    if err != nil {
		log.Panicf("Could not parse schema: %v", err)
    }

    // (...)
}
```

## Installation

```bash
go get -u github.com/s-augustovitko/mongo-schema-go/pkg/schema
```

## Important Notes

An important thing to consider is that ID fields instead of using the `primitive.ObjectId` type, make sure to use `interface{}`. For other complex structures like dates make sure to always use the `type` or `itemsType` tags. This is because this module uses no external dependencies, and focuses mostly on the reflect package.

Another thing to consider is that when ever you have arrays in your struct, the object you pass to the marshal function has to have an array with at least 1 value, or it will not be processed correctly.

The inline option is available in the field tag in order to merge the parent object and the child object so that they are at the same level.

```go
package main

import (
    "log"
    "time"

    "github.com/s-augustovitko/mongo-schema-go/pkg/schema"
)

type Data struct {
    Name    string
    // Needs to be filled to be considered in the schema
    Ages    []int
}

type DataParent struct {
    // Will have type of objectId in the schema
    ID          interface{}
    // Name and Ages from the Data struct will be at the same level as Id, Address and Dates (DataParent.Name|DataParent.Ages)
    Obj         Data `field:",inline"`
    // Will be a child struct of the struct (DataParent.Data.Name|DataParent.Data.Number|DataParent.Data.Apt)
    Address     struct{
        Name    string
        Number  int
        Apt     int
    }
    // needs to be filled to be considered in the schema
    Dates       []time.Time `itemsType:"date"`
}

func main() {
    obj := &DataParent{
        Obj: Data{Ages: []int{1}},
        Dates: []*time.Time{time.Now()}
    }
    out, warnings, err := schema.Marshal(obj, "Demo Object Schema", true)
    if err != nil {
		log.Panicf("Could not parse schema: %v", err)
    }

    // (...)
}
```

## Marshal Response Structure

The response of the marshal function has 3 parameters:

- The first one being the jsonSchema object has a type of `map[string]interface{}` which can be used together with the `CreateCollection` mongo function in order to create a schema or using the command `collMod` to update the schema.
- The second value is a list of errors of type `ErrorWithTag`, this is used so that you can get the Tag or Name of the value where the error occurs, the fields in this list of errors will not be in the final bson model, since it could not be processed correctly, but the rest of field will be processed normally.
- The third value is an error, if this error ocurrs, it means you are not sending a struct to the `Marshal` function, and the schema was not created.

```go
type ErrorWithTag interface {
	Name() string
	Tag() string
	Error() string
}
```

## More information

For more details on mongo schema please follow this link
[MongoDB Schema](https://www.mongodb.com/docs/manual/core/schema-validation/)

## Tags

This are the possible tags that the struct can contain with this package, all fields are optional, if everything is empty the schema will have no validations, but the main schema will still be created, for context:

- (string?,inline?) = first value is any string and the second value is optional and can only be inline (eg. "example,inline" || "example" || ",inline")
- (string,...) = comma separated strings (eg. "double,int,long")
- (string|string=string,...) = comma separated validations (eg. "required,min=1,max=20")
- (string) = string value

```go
// field name for Marshal function
// Can use bson instead (inline does not work with bson tag)
var field      = string?,inline?
// Field bson type
var type       = string,...
// Only for arrays
// field items bson type
var itemsType  = string,...
// Comma separated values for mongo schema validations
var validation = string|string=string,...
// Mongo schema description (Error message for validations)
descriptionvar = string
// Enum values (usually used for array and string values)
var enum       = string,...
// Same as validations but for array items
var items      = string|string=string,...
```

### Type && ItemsType

Possible type values are the following, these values are based on the mongodb documentation. This are all the possible values for the itemsType and type tags, multiple values can be used.
[Bson Types](https://www.mongodb.com/docs/manual/reference/bson-types/)

```go
var validBsonTypes = []string{
	"string",
	"double",
	"object",
	"array",
	"binData",
	"undefined",
	"objectId",
	"bool",
	"date",
	"null",
	"regex",
	"dbPointer",
	"javascript",
	"symbol",
	"javascriptWithScope",
	"int",
	"timestamp",
	"long",
	"decimal",
	"minKey",
	"maxKey",
}
```

### Validation && Items

This works for the validation and items tags.
For more information on the validations and how it works on the mongo schema, please refer to the following link.
[JSON Schema](https://www.mongodb.com/docs/manual/reference/operator/query/jsonSchema/#mongodb-query-op.-jsonSchema)

```go
var required            = "required" // All values
var uniqueItems         = "uniqueItems" // Array or Slice values
var min                 = "min=int|float" // All Values (adapts depending on the type)
var max                 = "max=int|float" // All Values (adapts depending on the type)
var multipleOf          = "multipleOf=int|float" // Int, Float or all number type values
var pattern             = "pattern=string" // String values
var patternProperties   = "patternProperties=string" // Strings when pattern exists
```
