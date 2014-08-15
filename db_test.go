package libpack

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func tmpdir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestInit(t *testing.T) {
	var err error
	// Init existing dir
	tmp1 := tmpdir(t)
	defer os.RemoveAll(tmp1)
	_, err = Init(tmp1, "refs/heads/test", "")
	if err != nil {
		t.Fatal(err)
	}
	// Test that tmp1 is a bare git repo
	if _, err := os.Stat(path.Join(tmp1, "refs")); err != nil {
		t.Fatal(err)
	}

	// Init a non-existing dir
	tmp2 := path.Join(tmp1, "new")
	_, err = Init(tmp2, "refs/heads/test", "")
	if err != nil {
		t.Fatal(err)
	}
	// Test that tmp2 is a bare git repo
	if _, err := os.Stat(path.Join(tmp2, "refs")); err != nil {
		t.Fatal(err)
	}

	// Init an already-initialized dir
	_, err = Init(tmp2, "refs/heads/test", "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetEmpty(t *testing.T) {
	tmp := tmpdir(t)
	defer os.RemoveAll(tmp)
	db, err := Init(tmp, "refs/heads/test", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Set("foo", ""); err != nil {
		t.Fatal(err)
	}
}

func TestList(t *testing.T) {
	tmp := tmpdir(t)
	defer os.RemoveAll(tmp)
	db, err := Init(tmp, "refs/heads/test", "")
	if err != nil {
		t.Fatal(err)
	}
	db.Set("foo", "bar")
	if db.tree == nil {
		t.Fatalf("%#v\n")
	}
	for _, rootpath := range []string{"", ".", "/", "////", "///."} {
		names, err := db.List(rootpath)
		if err != nil {
			t.Fatalf("%s: %v", rootpath, err)
		}
		if fmt.Sprintf("%v", names) != "[foo]" {
			t.Fatalf("List(%v) =  %#v", rootpath, names)
		}
	}
	for _, wrongpath := range []string{
		"does-not-exist",
		"sldhfsjkdfhkjsdfh",
		"a/b/c/d",
		"foo/sdfsdf",
	} {
		_, err := db.List(wrongpath)
		if err == nil {
			t.Fatalf("should fail: %s", wrongpath)
		}
		if !strings.Contains(err.Error(), "does not exist in the given tree") {
			t.Fatalf("wrong error: %v", err)
		}
	}
}

func TestSetGetSimple(t *testing.T) {
	tmp := tmpdir(t)
	defer os.RemoveAll(tmp)
	db, err := Init(tmp, "refs/heads/test", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Set("foo", "bar"); err != nil {
		t.Fatal(err)
	}
	if key, err := db.Get("foo"); err != nil {
		t.Fatal(err)
	} else if key != "bar" {
		t.Fatalf("%#v", key)
	}
}

func TestSetGetMultiple(t *testing.T) {
	tmp := tmpdir(t)
	defer os.RemoveAll(tmp)
	db, err := Init(tmp, "refs/heads/test", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Set("foo", "bar"); err != nil {
		t.Fatal(err)
	}
	if err := db.Set("ga", "bu"); err != nil {
		t.Fatal(err)
	}
	if key, err := db.Get("foo"); err != nil {
		t.Fatal(err)
	} else if key != "bar" {
		t.Fatalf("%#v", key)
	}
	if key, err := db.Get("ga"); err != nil {
		t.Fatal(err)
	} else if key != "bu" {
		t.Fatalf("%#v", key)
	}
}

func TestSetCommitGet(t *testing.T) {
	tmp := tmpdir(t)
	defer os.RemoveAll(tmp)
	db, err := Init(tmp, "refs/heads/test", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Set("foo", "bar"); err != nil {
		t.Fatal(err)
	}
	if err := db.Set("ga", "bu"); err != nil {
		t.Fatal(err)
	}
	if err := db.Commit("test"); err != nil {
		t.Fatal(err)
	}
	if err := db.Set("ga", "added after commit"); err != nil {
		t.Fatal(err)
	}
	db.Free()
	db, err = Init(tmp, "refs/heads/test", "")
	if err != nil {
		t.Fatal(err)
	}
	if val, err := db.Get("foo"); err != nil {
		t.Fatal(err)
	} else if val != "bar" {
		t.Fatalf("%#v", val)
	}
	if val, err := db.Get("ga"); err != nil {
		t.Fatal(err)
	} else if val != "bu" {
		t.Fatalf("%#v", val)
	}
}


func TestSetGetNested(t *testing.T) {
	tmp := tmpdir(t)
	defer os.RemoveAll(tmp)
	db, err := Init(tmp, "refs/heads/test", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Set("a/b/c/d/hello", "world"); err != nil {
		t.Fatal(err)
	}
	if key, err := db.Get("a/b/c/d/hello"); err != nil {
		t.Fatal(err)
	} else if key != "world" {
		t.Fatalf("%#v", key)
	}
}
