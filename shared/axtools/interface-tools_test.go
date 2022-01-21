package axtools

import (
	"gopkg.in/yaml.v3"
	"reflect"
	"testing"
)

const isMap = `
name: tg_is
lines: [25]
type: normal
reels:
  - [0,1,2,3,4,5,6,7,8,9,10]
  - [0,1,2,3,4,5,6,7,8,9,10]
  - [0,1,2,3,4,5,6,7,8,9,10,11]
`

const isArray = `
- name: tg_is
  lines: [25]
  type: normal
  reels:
    - [0,1,2,3,4,5,6,7,8,9,10]
    - [0,1,2,3,4,5,6,7,8,9,10]
    - [0,1,2,3,4,5,6,7,8,9,10,11]
`

func TestOneItemInterface(t *testing.T) {
	var item interface{}
	if err := yaml.Unmarshal([]byte(isMap), &item); err != nil {
		t.Fatal(err)
	}
	arr, err := ReadInterfaceAsArray(item)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.ValueOf(arr).Kind() != reflect.Array && reflect.ValueOf(arr).Kind() != reflect.Slice {
		t.Fatalf("is not array. is %v", reflect.ValueOf(arr).Kind())
	}
}

func TestMultiItemInterface(t *testing.T) {
	var item interface{}
	if err := yaml.Unmarshal([]byte(isArray), &item); err != nil {
		t.Fatal(err)
	}
	arr, err := ReadInterfaceAsArray(item)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.ValueOf(arr).Kind() != reflect.Array && reflect.ValueOf(arr).Kind() != reflect.Slice {
		t.Fatalf("is not array. is %v", reflect.ValueOf(arr).Kind())
	}
}
