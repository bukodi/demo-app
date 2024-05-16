package main

import _ "github.com/bukodi/demo-app/pkg/init_by_tags"

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
