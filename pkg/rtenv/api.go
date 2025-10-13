package rtenv

import "flag"

func EnvValue(key string) string {
	flag.StringVar(&key, "key", "", "key")
	return ""
}
