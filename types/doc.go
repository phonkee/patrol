/*
Custom database types for patrol

We have following custom types that can be stored to database:

ForeignKey - wrapper around int64 providing some custom methods
PrimaryKey - primary key for all models
GzippedMap - basically map that is stored to database gzipped
IntSlice - intslice that is stored as postgres ARRAY type
StringSlice - slice of strings that is stored as postgres string type
IsField - basically boolean data type that is used by postgres schema
*/
package types
