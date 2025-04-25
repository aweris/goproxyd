package main

import (
	"context"
	"crypto/md5"
	"hash"
	"io"
	"os"
	"path"
	"strings"
)

type cacher struct {
	root     string
	readonly bool
}

// NewHash returns a new instance of the `hash.Hash` used to compute the
// checksums of the caches in the underlying cacher.
func (c *cacher) NewHash() hash.Hash {
	return md5.New()
}

// Put sets the content to the underlying cacher.
func (c *cacher) Put(ctx context.Context, name string, content io.ReadSeeker) (err error) {
	if c.readonly {
		return nil
	}
	filename := path.Join(c.root, name)
	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	return err
}

// Get returns the matched content for the name from the underlying
// cacher. It returns os.ErrNotExist if not found.
func (c *cacher) Get(ctx context.Context, name string) (content io.ReadCloser, err error) {
	filename := path.Join(c.root, name)
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, os.ErrNotExist
		}
		return nil, err
	}

	return file, nil
}

// MIMETypeByFilename returns the MIME type of the file by its name.
func MIMETypeByFilename(name string) string {
	switch ext := strings.ToLower(path.Ext(name)); ext {
	case ".info":
		return "application/json; charset=utf-8"
	case ".mod":
		return "text/plain; charset=utf-8"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}
