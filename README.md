#  MEGTASK

MegTask is a task management application (HTTP API) that saves your tasks and
lets you update them when you are done.

# Features :zap
1. Create a user account.
2. Login to existing user account.
3. Create task.
4. Update task.
5. Delete task.

# Starting the Server: Perquisites

1. Go installed.
2. A database connection URL from mongodb.com


# How to start the application server :rocket

1. Ensure the latest version of `go` is installed on your device. Visit
   https://go.dev/doc/install to install `go`.

2. Clone this repo to your local device and run `cd megtask` on your terminal.

3. Lastly, run `go build` to build the executable and then run `./megtask --dbURL={enter your mongodb connection url here}` to start the HTTP server.

**NOTE**: Upload the [MEGTASK_POSTMAN_COLLECTION file](./MEGTASK_POSTMAN_COLLECTION.json) to postman to see the documented API endpoints.