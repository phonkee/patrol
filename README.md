# patrol


#### Warning: this project is in initial state. It's nowhere near to complete!!

Patrol is implementation of sentry server in go language.
It uses sentry protocol so you can use raven clients (currently protocol version 4 is supported).
Frontend is written in angularjs.

For demo you can try:
    http://demo.patrol.name

    username: demo
    password: demo

It's still very limited, under heavy development.

As database patrol uses exclusively postgres (for its advanced field types).
For queue you can currently use either redis or rabbitmq.

Goal of this project is not to replace sentry with all its features, but
create simple, portable, easily deployable solution for logging.

In first stable version all the static data will be embedded to binary, so there
will be ne dependencies.
All configuration is provided with command line arguments.


Development of frontend is located at https://github.com/phonkee/patrol-frontend/


Authors:
phonkee

Contributions:
Any ideas are welcome