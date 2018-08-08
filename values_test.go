package main

import (
	"testing"
	"os"
	"io/ioutil"
	"path"
)

func TestVals(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmp)

	y1 := path.Join(tmp, "y1")
	ioutil.WriteFile(y1, []byte(`
affinity: {}
image:
  pullPolicy: IfNotPresent
  repository: nginx
  tag: stable
service:
  port: 80
  type: ClusterIP
`), defaultDirectoryPermission)

	y2 := path.Join(tmp, "y2")
	ioutil.WriteFile(y2, []byte(`
image:
  tag: alpine
`), defaultDirectoryPermission)

	y3 := path.Join(tmp, "y3")
	ioutil.WriteFile(y3, []byte(`
service:
  port: 8080
`), defaultDirectoryPermission)

	v := valueFiles{
		y1, y2, y3,
	}

	mreged, err := vals(v)
	if err != nil {
		t.Error(err)
	}

	//fmt.Println(string(actual))

	expected := `affinity: {}
image:
  pullPolicy: IfNotPresent
  repository: nginx
  tag: alpine
service:
  port: 8080
  type: ClusterIP
`

	if actual := string(mreged); actual != expected {
		t.Errorf("expected %s, but got %s", expected, actual)
	}
}

func TestValueFiles(t *testing.T) {
	v := valueFiles{"a", "b"}

	v.Insert("x", 0)
	expected := "[x a b]"
	if actual := v.String(); actual != expected {
		t.Errorf("expected %s, but got %s", expected, actual)
	}

	v.Insert("y", 2)
	expected = "[x a y b]"
	if actual := v.String(); actual != expected {
		t.Errorf("expected %s, but got %s", expected, actual)
	}

	v.Insert("z", -1)
	expected = "[x a y b z]"
	if actual := v.String(); actual != expected {
		t.Errorf("expected %s, but got %s", expected, actual)
	}

}
