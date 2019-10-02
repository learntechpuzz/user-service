aws dynamodb create-table --cli-input-json Users.json --endpoint-url http://localhost:8000

aws dynamodb create-table --cli-input-json NextIdTable.json --endpoint-url http://localhost:8000

aws dynamodb put-item --table-name NextIdTable --item "{\"NextKey\": {\"S\": \"Users\"},\"NextId\": {\"N\": \"0\"}}" --return-consumed-capacity TOTAL --endpoint-url http://localhost:8000
