package tagfs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func isDir(dirname string, t *testing.T) bool {
	f, err := os.Open(dirname)
	if err != nil {
		t.Logf("Failed to open path %s: %v", dirname, err)
		return false
	}

	fi, err := f.Stat()
	if err != nil {
		t.Logf("Failed to retrieve file info on %s: %v", dirname, err)
		return false
	}

	return fi.IsDir()
}

var fsdir string
var fs *FileSystem

func TestMain(m *testing.M) {
	var err error
	var code int

	tempdir, _ := ioutil.TempDir("", "")
	fsdir = tempdir + "/tagfs"

	fmt.Printf("Using temp directory %s\n", tempdir)
	fs, err = Create(fsdir, 0700, 1)
	if err != nil {
		fmt.Printf("Failed to create filesystem in %s: %v\n", tempdir, err)
		code = -1
	} else {
		code = m.Run()
		fs.Close()
	}

	os.RemoveAll(tempdir)
	os.Exit(code)
}

func TestCreateFileSystem(t *testing.T) {
	t.Logf("Creating filesystem")
	if !isDir(fsdir+"/"+DB_DIR, t) {
		t.Fatalf("Expected directory %s to be created and it wasn't", fsdir+"/"+DB_DIR)
	}

	if !isDir(fsdir+"/"+FILE_DIR, t) {
		t.Fatalf("Expected directory %s to be created and it wasn't", fsdir+"/"+FILE_DIR)
	}
}

func createFile(name string, content []byte) error {
	file, err := fs.Create(name)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	file.Close()
	return err
}

func TestOpenFileSystem(t *testing.T) {
	expected := "The answer to the question is 42"
	createFile("test.txt", []byte(expected))
	t.Logf("Opening filesystem")
	fs, err := Open(fsdir)
	if fs == nil || err != nil {
		t.Fatalf("Failed to open filesystem directory: %v", err)
	}

	result := make([]byte, len(expected))
	file, err := fs.Open("test.txt")
	if err != nil {
		t.Fatalf("Failed to open file in pre-existing filesystem: %v", err)
	}

	file.Read(result)
	file.Close()
	if string(result) != expected {
		t.Fatalf("Expected file to contain \"%s\" but it contained \"%s\" intead", expected, string(result))
	}
	fs.Close()
}

func TestCreateFile(t *testing.T) {
	t.Logf("Creating file")
	err := createFile("test.txt", []byte{})

	if err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}
}

func TestOpenFile(t *testing.T) {
	expectedContent := []byte{0x41, 0x42, 0x43, 0x44}
	createFile("test.txt", expectedContent)
	file, err := fs.Open("test.txt")
	if err != nil {
		t.Fatalf("Failed to open existing file: %v", err)
	}

	content := make([]byte, 10)
	numRead, err := file.Read(content)
	file.Close()

	if bytes.Compare(content[0:numRead], expectedContent) != 0 {
		t.Fatalf("Failed to read file.  Expected %s got %s", string(expectedContent), string(content))
	}
}

func TestName(t *testing.T) {
	createFile("test.txt", []byte{})
	file, _ := fs.Open("test.txt")
	if file.Name() != "test.txt" {
		t.Fatalf("Expected name to be %s got %s instead", "test.txt", file.Name())
	}
}
