/*
Package static is responsible for handling static data (frontend)

static uses build flag named EMBED_STATIC
On this depends whether all the static data will be compiled to binary, or not.
Most of the time you'll don't need to handle this, but in case of fronted
development it's good not to have to recompile binary again.

@TODO: this is still not finished, need some work, mostly on go-bindata
serving.

*/
package static
