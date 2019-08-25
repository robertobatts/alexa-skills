module scorekeeper

go 1.12

require (
	dynamodb v0.0.0
	github.com/aws/aws-lambda-go v1.12.0
	github.com/aws/aws-sdk-go v1.23.8
	github.com/davecgh/go-spew v1.1.1
	golexa v0.0.0
)

replace (
	dynamodb => ../../dynamodb
	golexa => ../..
)
