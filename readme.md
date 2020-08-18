# Project binary

## Tech Stack

- Gorm
- Gin
- MySql
- Docker
- JWT

## Docker Image availble at DockerHub

### To pull image run following command

`$ docker pull imankurj/binary`

### To run container

`$ docker run --rm -p 8080:8080 --env-file .env imankurj/binary:auth:added`

## API Services

- GET
- POST
- PATCH
- DELETE

### API Endpoint localhost:8080

## Functions Example

### POST (LOGIN)

``` json
POST /login
request:

{
  "username":"<username>"
  "password":"<password>"
}

response:

"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImVhMzE1ZTBhLTA1NWEtNDc2ZS1hODRkLWE1YmZiNzJkYWYzZSIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTU5NzcxNjQ5MCwidXNlcl9pZCI6MX0.O3JuoY0Q78XgALIU3nLzA_G0YR8r-M2NHsMCku2vkmg"
```

## **Add Authentication Token in header on hitting service**

## Token Expiry Time is 1 Minute

## Header Format

`Authentication: Bearer <Token>`

### POST

``` json
POST /
request:

{
  "value":true,
   "key": "name(Optional)"
}

response:

{
  "id":"b7f32a21-b863-4dd1-bd86-e99e8961ffc6",
  "value": true,
  "key": "name"
}
```

### GET

``` json
GET /:id

response:

{
  "id":"b7f32a21-b863-4dd1-bd86-e99e8961ffc6",
  "value": true,
  "key": "name"
}
```

### PATCH

``` json
PATCH /:id
request:

{
  "value":false,
  "key": "new name(Optional)"
}

response:

{
  "id":"b7f32a21-b863-4dd1-bd86-e99e8961ffc6",
  "value": true,
  "key": "new name"
}
```

### DELETE

``` json
DELETE /:id

response:
HTTP 204 No Content
```

## To Run locally

`$ go install`

`$ ./Project_binary`

## For Testing

`$ go test`

## **Attention**: Before running the API edit the env file, add Database username and password
