package axtools

import (
	"testing"
)

func TestNameFromUnderlinedToCamel(t *testing.T) {
	if NameFromUnderlinedToCamel("tg_et") != "TgEt" {
		t.Fatalf("name is not equal %v", NameFromUnderlinedToCamel("tg_et"))
	}
}

func TestNameFromCamelToUnderline(t *testing.T) {
	if NameFromCamelToUnderline("TgEt") != "tg_et" {
		t.Fatalf("name is not equal %v", NameFromCamelToUnderline("TgEt"))
	}
}
