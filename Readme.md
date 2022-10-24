curl -X GET http://localhost:8000/user/{id}
curl -X GET -I http://localhost:8000/user/{id}

curl -X GET http://localhost:8000/user/{id}/tx
curl -X GET -I http://localhost:8000/user/{id}/tx

curl -X POST http://localhost:8000/user/{id}/{amount}
curl -X PUT http://localhost:8000/user/{id}/{amount}
curl -X PUT -I http://localhost:8000/user/{id}/{amount}


curl -X GET http://localhost:8000/tx/{id}

curl -X POST http://localhost:8000/transfer/{fromId}/{toId}/{amount}


curl -X GET http://localhost:8000/report/{month}/{year}


curl -X POST http://localhost:8000/reserve -H 'Content-Type: application/json' -d '{"userId":, "orderId":, "serviceId":, "amount":}'

curl -X POST http://localhost:8000/cancel -H 'Content-Type: application/json' -d '{"userId":, "orderId":, "serviceId":, "amount":}'

curl -X POST http://localhost:8000/approve -H 'Content-Type: application/json' -d '{"userId":, "orderId":, "serviceId":, "amount":}'
