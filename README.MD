# SMS Gateway

Sending and receiving SMS messages.

http://euromoby.com

### Running

```
go run cmd/main.go
```

### ENV variables

- PORT (default `8080`)
- WAIT_TIMEOUT (default `15`)
- DATABASE_URI (default `postgres://root:heslo@localhost:5432/smsgw_test?&pool_max_conns=10`)

### Testing

```
go test -v ./...
```

### Migrations

https://github.com/pressly/goose

```
cd migrations
goose postgres "user=root password=heslo dbname=smsgw_test" reset
goose postgres "user=root password=heslo dbname=smsgw_test" up
```

### Endpoints

#### Send Message

- POST /v1/sms/messages

- GET /v1/sms/messages/status/search
- GET /v1/sms/messages/status/:id

- GET /v1/sms/messages/outbound/search
- GET /v1/sms/messages/outbound/:id

- GET /v1/sms/messages/inbound/:shortcode/search
- GET /v1/sms/messages/inbound/:shortcode/:id
- PUT /v1/sms/messages/inbound/:shortcode/:id/ack

### Sandbox

Join Telegram channel @EuromobyBot
