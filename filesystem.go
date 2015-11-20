package tagfs

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

type Tags struct {
	data map[string]tagset
}

func (t Tags) add(tag string, file *File) error {
	if t[tag] == nil {
		t[tag] = make(tagset)
	}
	t[tag][file.name] = file.id
	return nil
}

const (
	DB_DIR   string = "db"
	FILE_DIR string = "files"
	TAG_SEP  string = "|"
)

type FileSystem struct {
	basedir   string
	dbdir     string
	filedir   string
	tags      Tags
	numShards int
}

func Create(basedir string, mode os.FileMode, numShards int) (*FileSystem, error) {
	err := os.MkdirAll(basedir, mode)
	if err == nil {
		err = os.Mkdir(basedir+"/"+DB_DIR, mode)
	}

	if err == nil {
		err = os.Mkdir(basedir+"/"+FILE_DIR, mode)
	}

	for i := 0; i < numShards; i++ {
		if err == nil {
			err = os.Mkdir(fmt.Sprintf("%s/%s/%d", basedir, FILE_DIR, i), mode)
		}
	}

	if err == nil {
		return Open(basedir)
	}
	return nil, err
}

func Open(basedir string) (*FileSystem, error) {
	fs := new(FileSystem)
	var err error
	fs.basedir = basedir
	fs.dbdir = basedir + "/" + DB_DIR
	fs.filedir = basedir + "/" + FILE_DIR

	fs.tags = make(Tags)

	files, err := ioutil.ReadDir(fs.filedir)
	if err == nil {
		for _, fi := range files {
			n, e := strconv.ParseInt(fi.Name(), 10, 32)
			num := int(n)
			if e == nil && num+1 > fs.numShards {
				fs.numShards = num + 1
			}
		}
		files, err = ioutil.ReadDir(fs.dbdir)
	}

	for _, fi := range files {
		if err == nil {
			name := fi.Name()
			ext := path.Ext(name)
			tags := strings.Split(fi.Name()[:len(name)-len(ext)], TAG_SEP)
			sort.Strings(tags)
			fs.tags[strings.Join(tags, TAG_SEP)], err = loadTagSet(fs.dbdir + "/" + fi.Name())
		}
	}

	if err == nil {
		return fs, nil
	}
	return nil, err
}

func (fs *FileSystem) Close() error {
	return nil
}

type File struct {
	*os.File
	name string
	id   uuid.UUID
}

func (fs *FileSystem) openOrCreate(name string, id uuid.UUID, action func(string) (*os.File, error)) (*File, error) {
	//host := binary.LittleEndian.Uint64(id[8:16]) % fs.numHosts
	shard := binary.LittleEndian.Uint64(id[:8]) % uint64(fs.numShards)

	f, err := action(fmt.Sprintf("%s/%d/%s", fs.filedir, shard, id.String()))
	if err == nil {
		file := new(File)
		file.File = f
		file.id = id
		file.name = name
		return file, nil
	}
	return nil, err
}

func parseFilename(name string) (string, []string) {
	tags := strings.Split(path.Dir(name), "/")
	return path.Base(name), tags
}

func (fs *FileSystem) Open(name string) (*File, error) {
	/* Convert the path to a list of tags */
	filename, tags := parseFilename(name)

	/* Lookup the file in the database:
	 * where name=? and tags=?
	 */
	tagsets := make([]tagset, len(tags)+1)
	tagsets[0] = fs.tags[""]
	for i, t := range tags {
		/* If the tag doesn't exist then */
		if fs.tags[t] == nil {
			return nil, &os.PathError{"open", name, errors.New("No such file or directory")}
		}
		tagsets[i+1] = fs.tags[t]
	}

	return fs.openOrCreate(name, intersection(tagsets...)[filename], os.Open)
}

func (fs *FileSystem) Create(name string) (*File, error) {
	filename, tags := parseFilename(name)
	if len(filename) > 255 {
		return nil, fmt.Errorf("%s: file name too long")
	}
	id := uuid.NewRandom()
	file, err := fs.openOrCreate(name, id, os.Create)
	file.name = filename
	if err == nil {
		fs.tags.add("", file)
		for _, tag := range tags {
			fs.tags.add(tag, file)
		}
	}

	return file, err
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Tags() []string {
	return nil
}

func (f *File) AddTag(tag string) {
}

func (f *File) RemoveTag(tag string) {
}

func (f *File) SetTags(tags []string) {
}

/*
func (f *File) Chmod(mode os.FileMode) error {
}

func (f *File) Chown(uid, gid int) error {
}

func (f *File) Close() error {
}

func (f *File) Read(b []byte) (n int, err error) {
}

func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
}

func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
}

func (f *File) Stat() (FileInfo, error) {
}

func (f *File) Sync() error {
}

func (f *File) Truncate(size int64) error {
}

func (f *File) Write(b []byte) (n int, err error) {
}

func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
}

func (f *File) WriteString(s string) (n int, err error) {
}
*/
