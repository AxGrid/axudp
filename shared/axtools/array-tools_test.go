package axtools

import (
	"testing"
)

func TestRemoveDuplicate(t *testing.T) {
	data := []int{1, 1, 3, 4, 5, 6, 6, 7, 1}
	data2 := RemoveIntDuplicate(data)
	if len(data2) != 6 {
		t.Fail()
	}

}
