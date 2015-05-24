/*
rest package

contains multiple helper packages for rest api:

metadata:
	metadata is basically description of rest requests for given endpoint.
	It is served via OPTIONS request, and frontend can use this information

response:
	This is helper to create json responses easily.

validator:
	Validator package provides functionality to create validators for
	serializers (structs mapped to requests)

*/
package rest
