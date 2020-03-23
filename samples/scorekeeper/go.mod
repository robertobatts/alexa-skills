module github.com/robertobatts/golexa/samples/scorekeeper

go 1.13

require (
	github.com/aws/aws-lambda-go v1.13.0
	github.com/aws/aws-sdk-go v1.23.8
	github.com/robertobatts/golexa v1.0.0
	github.com/robertobatts/golexa/samples/scorekeeper/dynamodb v1.0.0
)

replace (
	github.com/robertobatts/golexa => ../../
	github.com/robertobatts/golexa/samples/scorekeeper/dynamodb => ./dynamodb
)
