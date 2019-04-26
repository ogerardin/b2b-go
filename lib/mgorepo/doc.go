// Provides a generic polymorphic MongoDB repository.
//
// SaveNew/Update methods accept any struct or pointer to a struut, and save it as a Mongo document in a given database/collection.
// Retrieval methods (GetById/GetAll) instantiate and populate the same struct type that was saved, by using reflection.
//
// The saved document has the following additional fields:
//
// - "_id" is a Mongo ObjectId that identifies the record
//
// - "_t" is a string representation of the actual struct type that was passed. Since Go doesn't have a native mechanism
// for converting a type to/from a string, we use a custom naming scheme (provided by package typeregistry). For
// deserialization to work, the type must have been registered with typeregistry.Register().
//
// Internally, this package uses https://github.com/mitchellh/mapstructure to map to and from structures. As a consequence:
//
// * If you want to retrieve the "_id" field in your struct, you may use mapstructure tag to map it to any field of your
// struct, e.g.:
//          Id      string `mapstructure:"_id"`
// * If you want a sub structure to be "flattened" (not nested) in the saved document, you may use the "squash" option
// of mapstructure's tag, e.g.:
//          EmbeddedStruct `mapstructure:",squash"`
//
package mgorepo
