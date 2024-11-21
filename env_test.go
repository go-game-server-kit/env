package env

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	os.Setenv("APP_NAME", "name")
	os.Setenv("APP_ENV", "env")
	os.Setenv("APP_GROUP", "group")
	os.Setenv("APP_GROUP_TAG", "tag")
	os.Setenv("APP_VERSION", "test")
	LoadEnv()
}
