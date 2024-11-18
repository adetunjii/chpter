# Chpter Assessment

This solution contains the implementation of a **User Service** and an **Order Service** using gRPC for efficient communication between the two services. Both services handle specific responsibilities while seamlessly interacting with each other through defined APIs.

**Keynotes**
- No actual database logic was implemented in this assessment. Instead, mock implementation of the database logic was used for testing, as there were no requirements concerning this.

- **User service** runs on port `5050`

- **Order service** runs on port `5051`


## Project Structure

This project is a monorepo that contains a go workspace with 3 sub-modules; `order-svc`, `user-svc`, and a `shared` module defined in the `go.work` file. 
The shared module acts as a utility library that is shared across all the other modules. It also contains the proto files and the generated protobufs for each service.

## Running the application 

### Running in Docker
The easiest way to run this project is by using `docker compose`.
At the root of the project, run:

```bash
$ docker compose up --build
```

### Running locally

To run locally, each module contains a Taskfile.yml file where specific commands have been defined. These commands help automate common tasks such as building, testing, running, or managing dependencies for the module. Using the [Task](https://github.com/go-task/task) runner, you can execute these tasks in a simple and consistent way across all modules in the project.

Ex: Running tests within a module:
```bash
$ task test
```
