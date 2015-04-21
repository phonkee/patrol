// common constants for whole patrol application
// all enums for system will be find here
package settings

import "compress/gzip"

const (
	// version
	VERSION = "0.1"

	// debug (only for developers)
	DEBUG = true

	// queue name where to publish messages
	EVENT_QUEUE_ID = "post-messages"

	// builtin plugin ids
	AUTH_PLUGIN_ID     = "auth"
	COMMON_PLUGIN_ID   = "common"
	EVENTS_PLUGIN_ID   = "event"
	PROJECTS_PLUGIN_ID = "project"
	STATIC_PLUGIN_ID   = "static"
	TEAMS_PLUGIN_ID    = "teams"

	// padding of command in list
	LIST_COMMANDS_COMMAND_PADDING = 30
	LIST_ROUTES_COMMAND_PADDING   = 70

	// event plugin constants
	EVENT_WORKER_DEFAULT_GOROUTINES_COUNT = 2

	HTTP_SERVER_DEFAULT_HOST = "127.0.0.1:4434"

	AUTH_TOKEN_HEADER_NAME = "X-Patrol-Token"

	PAGING_DEFAULT_LIMIT_PARAM_NAME = "limit"
	PAGING_DEFAULT_PAGE_PARAM_NAME  = "page"

	PAGING_MAX_LIMIT     = 50
	PAGING_MIN_LIMIT     = 10
	PAGING_DEFAULT_LIMIT = PAGING_MIN_LIMIT

	ORDERING_DEFAULT_PARAM_NAME = "order"

	// http methods
	HTTP_DELETE  = "DELETE"
	HTTP_HEAD    = "HEAD"
	HTTP_GET     = "GET"
	HTTP_POST    = "POST"
	HTTP_PUT     = "PUT"
	HTTP_PATCH   = "PATCH"
	HTTP_OPTIONS = "OPTIONS"
	HTTP_TRACE   = "TRACE"

	SENTRY_AUTH_HEADER_NAME = "X-Sentry-Auth"
	SENTRY_AUTH_KEY         = "sentry_key"
	SENTRY_AUTH_SECRET      = "sentry_secret"
	SENTRY_AUTH_VERSION     = "sentry_version"

	SENTRY_TIMESTAMP_LAYOUT  = "2006-01-02T15:04:05"
	EVENT_PARSER_PROTOCOL_V4 = "4"

	// compression for queue
	RAW_EVENT_COMPRESSION_LEVEL = gzip.DefaultCompression

	// compression level for all database GzippedMap fields
	GZIPPED_MAP_COMPRESSION_LEVEL = gzip.DefaultCompression
)
