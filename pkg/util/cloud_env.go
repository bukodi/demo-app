package util

import "os"

func IsAWSLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}
