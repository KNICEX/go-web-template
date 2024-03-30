package snowflake

import "testing"

func TestGenID(t *testing.T) {
	Init(1, 1)
	for i := 0; i < 100; i++ {
		t.Log(GenID())
	}
}
