## Creating the binary file
`GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o idserver`

## Start the complete service with a postgres database and a hydra service
`docker-compose up`