# Project binary

## Tech Stack

- Gorm
- Gin
- MySql
- Docker

## Docker Image availble at DockerHub

### To pull image run following command

`$ docker pull imankurj/binary`

### To run container

`$ docker run --rm -p 8080:8080 imankurj/binary`

## API Services

- GET
- POST
- PATCH
- DELETE

### API Endpoint localhost:8080

## Functions Example

### POST

``` json
POST /
request:

{
  "value":true,
   "key": "name" // this is optional
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
  "key": "new name" // this is optional
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

## For Testing

`$ go test`
