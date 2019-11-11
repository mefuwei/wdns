package storage

import "testing"

func TestGetStorage(t *testing.T) {

	cases := []string{"redis"}
	for _, s := range cases {
		switch s {
		case "redis":
			GetStorage("redis", "localhost:6379", "", 1)
		}
	}
}
