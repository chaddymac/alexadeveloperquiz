GOOS=linux GOARCH=amd64 go build -o main main.go
zip main.zip main
aws lambda update-function-code --function-name alexadeveloperquiz --zip-file fileb://main.zip --publish
aws lambda invoke --function-name alexadeveloperquiz --log-type Tail response.json --query 'LogResult' --output text |  base64 -d
