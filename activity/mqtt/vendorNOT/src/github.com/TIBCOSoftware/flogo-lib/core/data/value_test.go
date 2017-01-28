package data

import (
	"fmt"
	"testing"
)

func TestGetAttrPath(t *testing.T) {

	a := "sensorData.temp"
	GetAttrPath(a)
	name, path, pt := GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "T.v"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "{T1.v}.myAttr"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "{T1.v}[0]"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "v[0]"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "v"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "{T1.v}"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
}
