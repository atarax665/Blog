Commands:

Run

```
docker compose up -d --build
```

See Logs

Producer

```
docker compose logs producer
```

Consumer

```
docker compose logs consumer
```

Sample Post Request

```
curl --location 'http://localhost:8080/bet' \
--header 'Content-Type: application/json' \
--data '{
    "username": "Abhinav",
    "team": "CSK",
    "amount": 1200
}'
```
