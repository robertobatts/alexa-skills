module scorekeeper

go 1.12

require (
	dynago v0.0.0
	github.com/aws/aws-lambda-go v1.12.0
	github.com/davecgh/go-spew v1.1.1
)

replace dynago => ../dynamodb
