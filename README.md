# IdService
Identity Service to work with [ory Hydra](https://github.com/ory/hydra) OAuth service.

## Run the Service standalone
`go run main.go`

## Create docker Image including go build
`docker build -t <tagname> .`

## Creating the binary file
`GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o idserver`

## Start the complete service with a postgres database and a hydra service
`docker-compose up`

## Creating a new User

POST 127.0.0.1:3000/user
````json
{
        "userName": "user-name",
        "name": "Homer",
        "lastName": "Simpson",
        "eMail": "mail@mailer.com",
        "password": "pwd",
        "applicationRoleDTO": [
            {
                "applicationName": "auth-code-client",
                "roles": [
                    "user",
                    "admin"
                ]
            }
        ]
    }

````

More should follow soon
