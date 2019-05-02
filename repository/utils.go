package repository

import (
  "log"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
)

func Map(values []bson.M, f func(bson.M) interface{}) []interface{}{
  mapped := make([]interface{}, len(values))

  for k, v := range values {
    mapped[k] = f(v)
  }

  return mapped
}


func _toObjectIdFrom(from string) func (v bson.M) interface{} {
  return func(val bson.M) interface{} {
    value, _ := primitive.ObjectIDFromHex(val[from].(string))
    return value
  }
}

//This is a curried function that will first receive a string argument that represents
//what it needs to map to. It then returns the mapping closure with the argument
//toMap already applied
func _partialApplyMap(toMap string) func (v bson.M) interface{} {
    return func (v bson.M) interface{}{
      return v[toMap]
    }
}

func _mapTo(values []bson.M, toMap string) []interface{} {
  return Map(values, _partialApplyMap(toMap))
}

func mountResponses(query *mongo.Cursor) []bson.M {
  numbers := make([]bson.M, 0)

  for query.Next(ctx) {
     var result bson.M
     err := query.Decode(&result)
     if err != nil { log.Fatal(err) }
     numbers = append(numbers, result)
  }

  return numbers
}

func findQuery(collection *mongo.Collection, mongoQuery bson.D) []bson.M {
  query, err := collection.Find(ctx, mongoQuery)
  if err != nil { log.Fatal(err) }
  defer query.Close(ctx)

  results := mountResponses(query)

  if err := query.Err(); err != nil {
    log.Fatal(err)
  }

  return results
}

func findAll(collection *mongo.Collection) []bson.M{
  return findQuery(collection, bson.D{})
}
