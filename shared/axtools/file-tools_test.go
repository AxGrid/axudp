package axtools

import (
	"testing"
)

func TestGetFirstDir(t *testing.T) {
	testPath("./data/tg_et.yml", "data", t)
	testPath("/Users/TestUser/src/axp/data2/tg_et.yml", "data2", t)
	testPath("./tg_et.yml", "", t)
}

func testPath(path string, assert string, t *testing.T) {
	if GetFirstDir(path) != assert {
		t.Fatalf("Wrong path '%v' != '%v' for %v", GetFirstDir(path), assert, path)
	}
}
