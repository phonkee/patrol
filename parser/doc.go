/*
Parser package is responsible for parsing sentry messages.
Currently patrol supports V4 sentry messages.
It's very easy to write parser for new protocol version though.
All parsers are registered to registry so patrol can user them.
*/
package parser
