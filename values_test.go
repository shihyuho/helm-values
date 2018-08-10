package main

import (
	"testing"
	"os"
	"io/ioutil"
	"path"
	"gopkg.in/yaml.v2"
	"fmt"
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

	merged, err := vals(v)
	if err != nil {
		t.Error(err)
	}

	expected := `affinity: {}
image:
  pullPolicy: IfNotPresent
  repository: nginx
  tag: alpine
service:
  port: 8080
  type: ClusterIP
`
	if actual := string(merged); actual != expected {
		t.Errorf("expected %s, but got %s", expected, actual)
	}
}

func TestValueFiles_Insert(t *testing.T) {
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

func TestMergeValues(t *testing.T) {
	b1 := []byte(`
service:
  port: 80
  type: ClusterIP
image:
  pullPolicy: IfNotPresent
  repository: nginx
  tag: stable
affinity: {}
`)
	src := yaml.MapSlice{}
	if err := yaml.Unmarshal(b1, &src); err != nil {
		t.Error(err)
	}

	dst := yaml.MapSlice{}
	dst = mergeValues(dst, src)
	expected := "[{service [{port 80} {type ClusterIP}]} {image [{pullPolicy IfNotPresent} {repository nginx} {tag stable}]} {affinity []}]"
	if actual := fmt.Sprintf("%v", dst); actual != expected {
		t.Errorf("expected %s, but got %s", expected, actual)
	}

	b2 := []byte(`
image:
 tag: alpine
`)
	src = yaml.MapSlice{}
	if err := yaml.Unmarshal(b2, &src); err != nil {
		t.Error(err)
	}

	dst = mergeValues(dst, src)
	expected = "[{service [{port 80} {type ClusterIP}]} {image [{pullPolicy IfNotPresent} {repository nginx} {tag alpine}]} {affinity []}]"
	if actual := fmt.Sprintf("%v", dst); actual != expected {
		t.Errorf("expected %s, but got %s", expected, actual)
	}
}

func TestSetValue(t *testing.T) {
	bytes := []byte(`
pullPolicy: IfNotPresent
repository: nginx
tag: stable
`)

	slice := yaml.MapSlice{}
	if err := yaml.Unmarshal(bytes, &slice); err != nil {
		t.Error(err)
	}

	slice = setValue(slice, "tag", "alpine")
	for _, item := range slice {
		if item.Key == "tag" {
			if actual := item.Value; actual != "alpine" {
				t.Errorf("expected %s, but got %s", "alpine", actual)
			}
			return
		}
	}
	t.Error("item of tag not in slice..")
}
