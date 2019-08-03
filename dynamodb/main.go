package main 

import (
	"fmt"
	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateNewSession() *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sess
}


func main() {

	sess := CreateNewSession()

	fmt.Print((sess.Config.Credentials.Get()))


}