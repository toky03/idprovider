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

## Unit test the Application
Example for adapter package:
`go test --coverprofile cover.out ./adapter && go tool cover -html=cover.out`

# Play the whole Flow
Complete Frlow from [Ory Hydra 5 Minute Tutorial](https://www.ory.sh/hydra/docs/5min-tutorial/)

`docker-compose up --build`

Create a new User:
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

Create a client in Hydra:
````yaml
docker-compose exec hydra \
    hydra clients create \
    --endpoint http://127.0.0.1:4445 \
    --id auth-code-client \
    --secret secret \
    --grant-types authorization_code,refresh_token \
    --response-types code,id_token \
    --scope openid,offline \
    --callbacks http://127.0.0.1:5555/callback
````

Create a sample App from Hydra:
````yaml
docker-compose exec hydra \
    hydra token user \
    --client-id auth-code-client \
    --client-secret secret \
    --endpoint http://127.0.0.1:4444/ \
    --port 5555 \
    --scope openid,offline
````
Navigate to http://127.0.0.1:5555/ and login with user-name and pwd.
