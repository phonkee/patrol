# patrol


#### Warning: this project is in initial state. It's nowhere near to complete!!

Patrol is implementation of sentry server in go language. It uses sentry protocol so you can use raven clients (currently protocol version 4 is supported). Frontend is written in angularjs so backend server serves only rest api.
All static data will be compiled into binary so there will be no need to copy resources etc..
Just copy binary, and you're ready to go. 

As database patrol uses exclusively postgres (for its advanced field types). For queue you can currently use either redis or rabbitmq.

Authors:
phonkee

Contributions:
Pull requests or ideas are welcome
