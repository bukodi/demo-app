package main

import (
	"os"

	_ "github.com/bukodi/demo-app/pkg/init_by_tags"
)

func main() {
	var err error
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Running in Lambda
		err = ServeInAWSLambda()
	} else {
		err = rootCmd.Execute()
	}

	if err != nil {
		panic(err)
	}
}
