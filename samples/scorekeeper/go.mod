module scorekeeper

go 1.12

require (
	dynamodb v0.0.0
	github.com/aws/aws-lambda-go v1.13.0
	github.com/aws/aws-sdk-go v1.23.8
	golexa v0.0.0
)

replace (
	dynamodb => ../../dynamodb
	golexa => ../..
)
