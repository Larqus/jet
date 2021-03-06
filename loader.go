// Copyright 2016 José Santos <henrique_1609@me.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jet

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// Loader is a minimal interface required for loading templates.
type Loader interface {
	// Exists checks for template existence.
	Exists(templatePath string) (string, bool)
	// Open opens the underlying reader with template content.
	Open(templatePath string) (io.ReadCloser, error)
}

// OSFileSystemLoader implements Loader interface using OS file system (os.File).
type OSFileSystemLoader struct {
	dir string
}

// compile time check that we implement Loader
var _ Loader = (*OSFileSystemLoader)(nil)

// NewOSFileSystemLoader returns an initialized OSFileSystemLoader.
func NewOSFileSystemLoader(dirPath string) *OSFileSystemLoader {
	return &OSFileSystemLoader{
		dir: filepath.FromSlash(dirPath),
	}
}

// Open opens a file from OS file system.
func (l *OSFileSystemLoader) Open(templatePath string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(l.dir, filepath.FromSlash(templatePath)))
}

// Exists checks if the template name exists by walking the list of template paths
// returns true if the template file was found
func (l *OSFileSystemLoader) Exists(templatePath string) (string, bool) {
	templatePath = filepath.Join(l.dir, filepath.FromSlash(templatePath))
	stat, err := os.Stat(templatePath)
	if err == nil && !stat.IsDir() {
		return templatePath, true
	}
	return "", false
}

type InMemLoader struct {
	files map[string][]byte
}

// compile time check that we implement Loader
var _ Loader = (*InMemLoader)(nil)

func NewInMemLoader() *InMemLoader {
	return &InMemLoader{
		files: map[string][]byte{},
	}
}

func (l *InMemLoader) Open(templatePath string) (io.ReadCloser, error) {
	f, ok := l.files[templatePath]
	if !ok {
		return nil, fmt.Errorf("%s does not exist", templatePath)
	}

	return ioutil.NopCloser(bytes.NewReader(f)), nil
}

func (l *InMemLoader) Exists(templatePath string) (string, bool) {
	_, ok := l.files[templatePath]
	if !ok {
		return "", false
	}
	return templatePath, true
}

func (l *InMemLoader) Set(templatePath, contents string) {
	templatePath = filepath.ToSlash(templatePath)
	templatePath = path.Join("/", templatePath)
	l.files[templatePath] = []byte(contents)
}
