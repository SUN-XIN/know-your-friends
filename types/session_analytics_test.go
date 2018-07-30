package types

import (
	"testing"
)

func TestScyllaDBKey(t *testing.T) {
	si := SessionIntegrate{
		UserIDOwner:  "toto",
		UserIDFriend: "tutu",

		TotalDuration: 100,
		Day:           1532725644,
		Type:          SessionTypeMostSeen,

		IsInSignPlace: true,
	}

	k := si.ScyllaDBKey()
	kOK := "toto-tutu-1532725644-1"
	if k != kOK {
		t.Errorf("Expect key %s, but get %s", kOK, k)
	}
}
