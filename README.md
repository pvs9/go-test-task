# go-test-task

This is a test-task with TodoItems implemented with such functionality

* Every 3 seconds a new TodoItem is created and sent to AWS SQS provided in the config
* `/todo` endpoint provides the latest TodoItem created as JSON response
* `APP_MODE=consumer` launches app AWS SQS consumer which consumes messages from configured queue and sends them to log

I used Clean Architecture approach building this application as
straightforward as I could, but there could be some architectural errors due to lack of time.

I did not have enough time to prepare some unit tests for my code, but, I guess, it is easy to understand how fast
they can be implemented with such architecture

## Two apps in the task

There were two apps needed in the test task, but I implemented quite a nice approach, inspired by the RoadRunner that we can
launch a single app in different modes (app, consumer) depending on the env APP_MODE variable with a default fallback to app.

It looks nice and clean and no unnecessary dependencies are used when consumer app mode is initialized.

## Consumer

Consumer app mode is built following Dependency Injection best practices, allowing us to swap the queue realisation.
Moreover, we can customize the handler for messages, but there will always be a default one. There are also some nice approaches
in receiver section - we use goroutines and can scale the receivers count if needed via config file provided.

## Swagger

Even though we have just one API endpoint, I decided to include Swagger documentation available via `/swagger/index.html` route
just to show that I love this approaches in app creation))

## Additional info

I used external cli tool to create database migrations in `./schema` folder as it looks quite convenient to me.

## How to use

- To start, please install: docker, docker-compose and go

- `'make up'` on root folder run it will start test-task docker
  containers with all required dependencies. 
- `'make bash'` to open container bash window 
- `localhost:3333` to access go container from localhost
- `localhost:3333/swagger/index.html` to access Swagger documentation
- `make aws-cli foo bar` to execute the aws cli with parameters foo bar
- `make aws-cli sqs list-queues` to access the sqs queue on aws localstack
- If you import dependencies to your go code, please use `make stop` and `make up` again to automatically download them
- To access the sql server, use mysql:3306 as address and port.
