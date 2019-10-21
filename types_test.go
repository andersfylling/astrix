package astrix

import (
	"testing"
)

func TestGetTypes(t *testing.T) {
	if _, err := GetTypes(".", nil, nil); err != nil {
		t.Fatal(err)
	}
}
