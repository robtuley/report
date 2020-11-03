package report

import "testing"

func TestULIDGenerator(t *testing.T) {
	str1 := createULID()
	str2 := createULID()

	if str1 == str2 {
		t.Fatal("not different", str1, str2)
	}
	if len(str1) != 26 {
		t.Fatal("wrong length", str1)
	}
}

func TestRandomStringGenerator(t *testing.T) {
	str1 := randString(64)
	str2 := randString(64)

	if str1 == str2 {
		t.Fatal("not different", str1, str2)
	}
	if len(str1) != 64 {
		t.Fatal("wrong length", str1)
	}
}
