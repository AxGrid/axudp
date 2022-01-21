package axtools

import (
	"testing"
)

func TestMapReadInterface(t *testing.T) {
	m := make(map[string]interface{})
	m["test1"] = 12
	m["test2"] = "15"
	if GetIntFromMap(m, "test1", 5) != 12 {
		t.Fatal("not equals")
	}
	if GetIntFromMap(m, "test2", 5) != 15 {
		t.Fatal("not equals")
	}
	if GetIntFromMap(m, "test3", 5) != 5 {
		t.Fatal("not equals")
	}
}

func TestMapReadString(t *testing.T) {
	m := make(map[string]string)
	m["test1"] = "17"
	m["test2"] = "15"
	if GetIntFromMap(m, "test1", 5) != 17 {
		t.Fatal("not equals")
	}
	if GetIntFromMap(m, "test2", 5) != 15 {
		t.Fatal("not equals")
	}
	if GetIntFromMap(m, "test3", 5) != 5 {
		t.Fatal("not equals")
	}
}
