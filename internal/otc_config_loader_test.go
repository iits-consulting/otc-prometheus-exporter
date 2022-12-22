package internal

import (
	"testing"
	"time"
)

func TestTokenIsValidAt(t *testing.T) {

	t1 := Token{
		Secret:    "",
		IssuedAt:  "2022-12-21T17:25:52.202000Z",
		ExpiresAt: "2022-12-23T17:25:52.202000Z",
	}

	testTimeAfter, _ := time.Parse(time.RFC3339, "2022-12-23T17:25:52.202000Z")
	testTimeBefore, _ := time.Parse(time.RFC3339, "2022-12-21T17:25:52.202000Z")

	valid, _ := t1.IsValidAt(testTimeAfter)
	if valid {
		t.Fatalf("Token is expired")
	}
	valid, _ = t1.IsValidAt(testTimeBefore)
	if !valid {
		t.Fatalf("Token")
	}

}
