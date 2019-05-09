package repository

import (
  "log"
  "sort"
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

func mountQueryFromId(_id string) bson.M {
  parsedId, _   := primitive.ObjectIDFromHex(_id)
  mongoQuery := bson.M{"_id": parsedId}
  return mongoQuery
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
func PartialApplyMap(toMap string) func (v bson.M) interface{} {
    return func (v bson.M) interface{}{
      return v[toMap]
    }
}

func _mapTo(values []bson.M, toMap string) []interface{} {
  return Map(values, PartialApplyMap(toMap))
}

//this is to overcome the lack of sort function in the current mongodb driver
//or, at least, I couldn't find one, lol
//one other thing: it copies the results to another variable in order to avoid
//in-place modification (because it'd bring side effects)
func sortByUrgency(results []primitive.M) []primitive.M{
  results_clone := make([]primitive.M, len(results))

  copy(results_clone, results)

  sort.Slice(results_clone, func(i, j int) bool {
    first_operation, second_operation := results_clone[i]["operation"].(bson.M), results_clone[j]["operation"].(bson.M)
    first_charge_order, second_charge_order := first_operation["charge_order"].(float64), second_operation["charge_order"].(float64)

    if (first_charge_order == second_charge_order) {
      return results_clone[i]["event_date"].(primitive.DateTime) < results_clone[j]["event_date"].(primitive.DateTime)
    }

    return first_charge_order < second_charge_order
  })

  return results_clone
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

func insertOne(collection *mongo.Collection, value interface{}) (string, error){
  returnedValue, err := collection.InsertOne(ctx, value)

  return returnedValue.InsertedID.(primitive.ObjectID).Hex(), err
}

func findOneById(collection *mongo.Collection, _id string) *mongo.SingleResult {
  return collection.FindOne(ctx, mountQueryFromId(_id))
}

func findOneQuery(collection *mongo.Collection, mongoQuery bson.M) *mongo.SingleResult {
  return collection.FindOne(ctx, mongoQuery)
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
