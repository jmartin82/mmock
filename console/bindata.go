// Code generated by go-bindata.
// sources:
// tmpl/css/style.css
// tmpl/index.html
// tmpl/js/script.js
// DO NOT EDIT!

package console

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _tmplCssStyleCss = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xc4\x55\x6f\x6b\xdb\x3e\x10\x7e\x5d\x7f\x0a\xfd\xe8\xbb\xf2\x73\x9c\x74\x29\xc9\x6c\x18\x64\x5b\x03\x83\xd1\x8d\x7d\x81\x21\x4b\x27\x5b\x8b\xac\x33\xb2\x9c\x74\x84\x7e\xf7\x49\x8e\xf3\xaf\x8d\xbb\x78\x0b\x2c\x21\x46\xe8\xf4\x3c\xf7\xdc\x73\x67\x25\x8e\xc3\x15\xa4\x0b\x69\xc3\x8a\x19\x54\x2a\xa5\x86\xac\x03\xe2\x3e\x2b\xc9\x6d\x1e\x93\xd1\x6d\xf9\x98\x04\x4f\xc1\x89\x83\xa1\x35\x94\x2d\xda\xe3\xa9\x5b\x66\x06\x6b\xcd\x43\x86\x0a\x4d\x4c\xae\x3f\xcc\xee\xc7\xf3\x59\xb2\x09\xa3\xe1\x60\x42\x05\xc2\x3a\xce\xf2\x91\x54\xa8\x24\x27\xd7\x8c\xb1\x4e\xf6\xbc\x2e\xd2\x6e\xf6\xf1\xed\x34\x65\xf4\x55\x70\x9c\xe3\x12\x4c\x37\xc5\xe8\xe3\x64\xfc\x7e\xee\x29\x82\x60\x50\x81\x02\x66\x81\x7f\x37\xb8\x22\x6b\x22\x50\x5b\x47\x2b\xb3\xdc\x09\x4e\x51\xf1\xe4\xc9\x9d\xb1\x46\xea\xcc\x45\x5b\x86\xcc\x00\xe8\x84\xb8\x88\x76\xe9\x7c\xaa\x6d\x84\x53\xb3\x40\x43\x75\x06\x4d\x38\x45\x54\x40\xf5\x3e\x9e\xaa\x1a\x5a\xa0\x52\xfb\xed\x82\x66\xa0\x2d\x6d\x22\x0b\xf8\xb9\x0f\x18\xe0\x7e\xd3\xe9\xcc\x25\xe7\xa0\xbf\x79\x91\x4d\x5d\x25\xe5\xdc\x89\x8a\xc9\x90\xfc\x27\x8b\x12\x8d\xa5\xda\x26\x5d\x25\x8b\x3b\xff\x6d\x4b\x2e\xa9\x06\x35\x58\x49\x9b\x87\x9a\x2e\x43\x4b\xd3\x8a\x6c\x36\xc3\x1c\xa8\x67\xf5\x29\x8e\xb3\xdc\xb9\xde\xf9\xdf\xd0\x3f\x3d\xcf\x69\x96\xed\x6a\x1d\x5c\xb5\xad\x4f\xd1\x5a\x2c\x62\xa2\x51\xc3\xeb\xb8\x1f\x75\x65\xa5\x90\xc0\x1d\xb8\xa0\x26\x93\x7a\x07\x0e\x47\x9b\xa4\x41\x74\x73\x73\x43\xbe\xce\x1e\xee\x3f\x93\x4f\x0f\xf3\x2f\x24\xb7\xb6\x8c\xa3\xc8\xf9\x6c\x2b\x2d\xcb\x72\xc0\xb0\x88\x9a\x15\xd8\x2a\x12\x40\x6d\xed\x2c\x8c\x9a\x9c\x55\x78\x9c\xd4\x51\x45\xc1\xb1\x90\xd6\x05\xa9\x05\xee\x6b\x21\xef\x88\x92\xee\x41\xff\xef\x75\x7a\x33\x85\x3d\x31\x02\x59\x5d\xb9\x0e\x5f\x6d\x1b\xf7\x66\x34\x19\x4e\x45\x63\xdc\x59\x3c\x03\x2c\x41\xf7\x11\xbb\x03\xf4\xd3\xbb\x87\x35\x92\xff\xa5\x35\x57\x27\xe6\x3d\x65\x30\x15\xa3\x64\x37\x85\xed\xbe\xbb\xb9\x74\x55\x52\x03\xfe\x55\x39\xd7\x52\x25\x07\x94\x59\xb9\x84\x7e\x33\x70\x00\xea\x5b\xef\x21\xb4\x57\xd9\x42\x88\x17\x35\xbf\xf0\x62\xf3\x52\xfd\xad\x25\xdc\x60\xc9\x71\xa5\xc9\x6e\x15\x16\xa0\xeb\xee\x5b\x97\xbf\x05\x2e\x26\x47\x7f\x0b\xcf\x15\x5e\x20\xfd\x76\x5c\x5a\x1d\xcf\x4c\x73\x3b\x97\x4c\xd2\xbb\xaf\xbf\x25\xdc\x76\xbb\xc3\xc3\x8b\x1a\xf5\x67\x53\x7d\x16\xdf\x05\x8d\x39\xf5\x2e\x1c\x76\xb6\x19\xf9\x0e\xbb\x76\x97\xe7\xaf\x00\x00\x00\xff\xff\xf9\xa4\x9d\x54\xef\x08\x00\x00")

func tmplCssStyleCssBytes() ([]byte, error) {
	return bindataRead(
		_tmplCssStyleCss,
		"tmpl/css/style.css",
	)
}

func tmplCssStyleCss() (*asset, error) {
	bytes, err := tmplCssStyleCssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tmpl/css/style.css", size: 2287, mode: os.FileMode(420), modTime: time.Unix(1484959285, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _tmplIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xbc\x58\x6d\x53\xe3\x38\x12\xfe\x3c\xfc\x0a\xad\xf7\xaa\x96\xa9\x39\xd9\x24\x19\x66\x18\x48\x52\xc5\x40\xf6\x38\x8a\x0c\x19\xc2\xc1\xed\x7e\xd9\x52\x6c\x25\x56\x90\x2d\x23\xc9\x79\x19\x8e\xff\x7e\x2d\xdb\x89\x1d\xc7\x0e\xb0\xb7\x75\xa9\x99\xc4\x96\xbb\x9f\x6e\xf5\xcb\xe3\x16\xed\x9f\xce\xaf\xcf\x6e\x7f\x1b\xf4\xd0\xc5\x6d\xff\xaa\xbb\xd7\xf6\x75\xc0\xbb\x7b\xf0\x4b\x89\xd7\xdd\x43\xf0\x69\x07\x54\x13\xe4\x6b\x1d\x61\xfa\x18\xb3\x59\xc7\x3a\x13\xa1\xa6\xa1\xc6\xb7\xcb\x88\x5a\xc8\x4d\xef\x3a\x96\xa6\x0b\xed\x18\xfd\x13\xe4\xfa\x44\x2a\xaa\x3b\xff\xba\xfd\x15\x1f\x59\xc8\xc9\x90\x34\xd3\x9c\x76\xfb\x7d\xe1\x3e\x20\x00\x51\x82\xd3\xb6\x93\x2e\xa6\x02\xca\x95\x2c\xd2\x48\x49\xb7\x63\x19\x8b\xea\xd8\x71\x5c\xe1\x51\x7b\xfa\x18\x53\xb9\xb4\x5d\x11\x38\xe9\x25\x6e\xd8\x8d\xa6\x7d\x60\x07\x2c\xb4\xa7\xca\xea\xb6\x9d\x54\xf7\xed\x40\x01\x9b\x48\xa2\x29\x00\x36\xed\xc6\x9b\xf0\xbc\x70\xaa\x6c\x97\x8b\xd8\x1b\x73\x22\x69\x02\x4a\xa6\x64\xe1\x70\x36\x52\x8e\x4f\x42\x8f\xd3\x11\x04\x02\x00\x1d\xe3\xeb\xc1\xe6\xda\x96\x91\x9f\x30\x46\x57\xe0\x8a\xd2\x10\xd4\x20\x62\x9c\x7a\x08\x14\x10\xf8\xc4\xc6\x0c\x6e\xce\x86\x43\x84\x71\x26\xcd\x59\xf8\x80\x24\xe5\x1d\x4b\xe9\x25\xa7\xca\xa7\x54\x5b\xc8\x97\x74\x9c\xbb\x18\x90\x05\x78\x69\x8f\x84\xd0\x4a\x4b\x12\x99\x1b\xe3\xe5\x7a\xc1\x69\xd9\x2d\xfb\x93\xe3\x2a\x95\xaf\x25\x31\x80\x15\x0b\x31\xc8\xec\x44\x32\xbd\x04\x1b\x3e\x69\x1d\x7d\xc4\x8d\xc7\xa3\xe0\xf6\xf2\xfa\x74\xb8\x38\x9a\x36\x4e\xe3\x0f\xe4\xf0\xfe\xfc\x2e\x1c\xb0\x26\x7f\xf8\x75\x3c\x9f\xf7\x4e\xc9\x91\x7f\x7e\xee\x4d\x7f\xe7\xd1\x15\x9d\x2c\xfc\xe9\x5d\xbf\xd7\x18\x4f\xa6\xf7\x83\x7f\x04\x0f\x3f\xd4\x67\x2b\xf1\xdd\x7c\x5c\x29\x94\x12\x92\x4d\x58\xd8\xb1\x48\x28\xc2\x65\x20\x62\x88\x49\x1e\x8a\xeb\x48\x33\x11\x12\x8e\xb4\x4f\x03\xfa\x7f\xd8\x38\x4e\x0c\xed\xda\xfe\xf8\xea\xbe\xf9\xed\xa0\xc1\xfb\x8f\x53\xf2\xf0\xf5\x61\xd1\xe2\x4e\xff\x4b\x8f\xf8\xf1\x3c\x1a\x8e\xe9\xb7\xd9\xdd\xa7\xd6\xe5\x21\xfd\x11\xb6\xe2\xdf\x7f\x90\xe8\xf6\x20\xfe\xdc\xfb\x4d\xfd\xbb\x3f\xfd\x7e\xf7\xe1\xa0\x17\x1e\xca\xb7\x6c\x7f\x67\x25\x5c\x92\x19\x19\xa6\x35\xb9\x8e\x4b\x55\x8d\xbe\x36\x0e\xd3\x72\xfe\xa7\x95\xfb\x3f\x08\x86\xa3\xcb\xf3\xde\x05\x23\x7c\x1c\xc4\x5f\xbf\x7e\x1f\x7c\x3a\xfd\xf8\x5d\x46\xf2\xf1\xf0\xfa\x6e\x7c\xdf\xfa\x3c\xb8\xb9\x69\x4d\x0f\x7b\x57\x8f\x0b\xa5\x1a\xcb\xbb\xc7\x6b\x1d\xd2\x28\xbc\xb8\x1b\x7c\x21\x97\x9f\x17\xc3\x97\xf7\x5f\xdd\x77\x1a\x98\x26\x23\x98\x7c\xe7\x56\xba\x55\x70\x3d\x95\xaa\x6a\xa9\xea\x5a\x29\xc0\x25\x89\x4e\x6b\xc7\xd4\x42\x22\x95\x64\xbf\x48\x7e\x21\x09\x40\x7e\xc6\xe8\x3c\x12\x52\x17\x28\x6f\xce\x3c\xed\x77\x3c\x3a\x63\x2e\xc5\xc9\xcd\xdf\x21\x6a\x4c\x43\x7c\xb0\x72\x09\xa7\x9d\x86\xf5\x9a\x9d\x74\xb3\xb8\xfc\x6d\x1f\x79\xc2\x8d\x03\x00\x47\xef\x6d\x09\x1c\xbc\xdc\x1f\xc7\xa1\x6b\x1a\x61\xff\x3d\x7a\x5a\x87\x0f\xa1\x19\x91\xb0\x31\xa0\x30\xa5\xaf\xc4\x64\x42\x25\xea\xa0\x90\xce\xd1\x4d\x71\x6d\xff\xfd\x49\xa6\xf2\x2e\xd5\x98\xab\x4c\xec\x9e\x8e\x86\x40\xc3\x54\xef\x5b\x73\x53\x28\x16\xfa\x80\xb8\x70\x89\xb1\x64\xfb\x02\x0a\xef\x03\xb2\x1c\xea\xfa\xc2\x2a\x62\xcc\x95\x2d\xc2\x80\x2a\x45\x26\x14\x90\xd6\xbe\xd1\xdc\xb9\x77\x6b\x63\xb9\xdc\xe5\xf0\xfa\x9b\x1d\x99\x77\xc2\x3e\x9d\xc1\xe6\x6c\x8f\x68\x52\xc4\x85\x7f\x1b\x9b\xb1\xb9\x98\xf4\x42\x2d\x97\xfb\x19\xc8\xfb\x93\x77\xa9\x70\xae\xf3\xbc\xd2\x7f\xce\x91\xd2\x58\xaf\x6b\xa0\xed\xa4\xef\xb1\xbd\xf6\x48\x78\xcb\x2c\x13\xbe\xcc\x2e\x3c\x36\x43\x2e\x27\x4a\x41\xee\x21\xa3\x84\x85\x54\xae\x53\x61\x04\x1b\xe5\x77\x15\xac\xe4\x8f\x0b\xea\x52\xcc\x0b\x8a\xc9\x53\xa2\x98\x47\x73\x78\x8e\x03\x0f\x37\x51\x44\x3c\x0c\x45\xef\x6b\x7c\x50\x52\x48\x94\x62\xbe\xd2\x08\xc9\x0c\xca\x6e\x86\xa1\xf7\xb9\x4a\xae\x94\x26\x90\x2f\xaf\x42\x2d\xab\xf4\x95\x2a\x81\x94\xcc\x28\xb4\x02\xc9\x0a\x7b\x0a\x85\x96\x46\xe4\xf8\x04\x96\x55\x44\x42\x14\x64\xc2\x23\xe2\x41\x86\xa2\x98\xf3\xd4\x2f\xab\xdb\x84\xf0\x81\x44\x77\x55\x49\x0a\xe2\x49\xa0\xaf\x38\xab\xb5\x5c\x67\xea\x74\x24\x62\x5d\xaf\xdd\x76\x62\x5e\x8a\x9a\x93\x84\xad\xb4\xb8\x91\xa7\x24\x90\x8d\xca\xe8\x01\x6b\xb2\x70\x24\x16\x48\x0b\xc1\xe1\x25\xbb\x62\xc6\x3a\xb0\xed\xac\xd5\x99\x5c\x28\xdc\x68\xd6\xc8\x96\xe5\x8b\xb1\xac\x55\x48\x94\xc8\x4a\x65\xa4\x43\x04\xff\xb1\x47\xc7\x24\xe6\x3a\xb9\x5e\x00\x2f\x25\xa3\x51\xc7\xba\xa1\x63\x09\xd4\x05\x84\xec\x25\xb2\x67\x9c\x12\x09\x7d\xf2\x02\x7e\x62\x23\x49\x76\x66\x66\xc2\x97\x91\xcf\xa0\xd4\xd1\xfa\x0a\x4b\x1a\x88\x19\xc5\x8a\x4d\x42\xab\x8b\xb2\xd4\xef\x76\x1b\xf2\x59\x1f\x08\x07\x22\x51\x13\xd3\x1d\x8f\x20\x73\x8e\x89\x73\x45\xc2\x92\xe7\x55\xe9\x40\xe0\xaa\x4b\xe5\xa1\xe1\xfc\x4a\xe4\xba\x65\x63\x0c\x52\x5f\x59\x1d\xe6\xd9\x9b\x6a\x08\xc2\x45\x39\x4a\xbe\xd7\xe9\xcb\x99\x2c\x4d\x59\x7e\x3f\x30\x72\x75\x35\x07\xb6\x57\x8c\xc9\x99\xd2\xaf\x09\x86\x26\x23\x4e\x21\x87\x2a\x02\x86\x4a\x9a\xbe\x3e\x33\x89\x6c\xd1\x1f\x9c\xac\x58\x1b\x58\x28\x45\xf4\xa1\x26\xcc\xeb\x25\xa9\x3b\xbc\x26\xc7\x94\x27\x3c\x31\x0f\x5f\x2a\x6d\x6d\x78\x77\x85\xec\x33\xcf\xa3\x49\x45\xbf\xa2\x60\xb5\x7c\x59\x28\x15\xf4\x8a\x25\xa1\x02\xdc\x34\xdc\xa6\xa5\x08\x27\xdd\x73\x98\x9d\xa0\x96\xd3\x1b\x38\x62\x78\x7f\x16\xb3\x91\x63\x9e\xc1\xf1\xe1\xaf\xc6\xec\x53\xed\x0b\xef\x2f\x41\x3d\xca\x51\x07\x44\xfb\x6f\xc4\x04\xa9\x17\xe2\x0e\x12\xf9\x99\xb0\xde\xad\xfc\x3d\xbb\x03\x68\xb7\x10\x08\x98\x22\x7c\x13\x8b\xec\xea\xf5\xb4\x9f\xd3\x16\x7d\xa1\x9b\xe7\x94\x73\x64\xbe\xb0\x42\x66\x46\x5b\xf1\x78\x9b\x06\xdd\x7f\x26\x30\x20\xa7\x51\x1c\xc1\xfc\x42\xbd\xe3\x8c\x5e\x4d\x4f\x99\xf5\x3f\xb2\x75\xab\xfb\xf4\x64\x2e\xd0\x7f\x90\xf9\x39\xfe\xa5\xdf\xc7\x9e\x87\x97\xf0\x41\x17\x17\xc7\x41\x70\xac\xd4\x2f\xcf\xcf\x19\xd9\xb6\x1d\x00\xaf\xf0\xbf\x6a\x69\xc5\x93\x5f\xca\x1b\x69\xbf\xa4\x5e\xb8\x2d\x5e\x9a\x83\xc6\x6d\xaf\x3f\xb8\x3a\xbd\xed\x0d\xb7\x8e\x12\x45\xb2\xa0\x66\x1a\xdb\x18\x9d\x17\x38\x3f\xce\x62\x4d\x83\x88\xc3\x66\x8b\xf3\x93\x96\x66\xff\x04\x6b\x18\xe9\xcc\x7b\x0c\x3c\xe7\x24\x52\x40\x39\xe9\x32\x91\x13\x0a\x43\xf4\xcf\x2b\x13\x1e\xcc\xda\x8c\x2b\xfc\xf4\x94\xad\xfc\x11\xc6\xc1\xf3\xb3\x95\x0f\x36\xae\x90\x1e\x8c\x9b\x19\x22\xca\x05\x01\x5a\x48\x10\x2d\xc5\x60\xab\x4b\x5a\x26\x37\x2b\x25\x93\x1c\x93\x86\x72\x83\x54\xb6\x6c\xd1\x96\xf7\x67\xd4\x82\xa4\xd7\x5f\xa7\xf8\xb9\xa8\x18\x41\x3b\x97\xd5\x36\xfb\x75\x8b\x34\x13\x44\xc1\x4d\x7d\x75\xac\x8f\xd6\x26\x11\xdf\x54\x8e\x3d\xc5\x36\xc8\x02\x4d\x42\x6c\x3a\x15\xe5\x79\x2b\x16\x44\x5d\xb6\x60\xda\x2e\x63\x97\xf1\x57\x53\x5c\x13\x8e\x70\xe6\xc0\xd5\xb1\x02\xa8\x05\x66\xf2\x1a\x1d\xa3\xc6\x41\xb4\x80\xe1\x71\x6f\xfb\x35\x3b\x67\xda\xc7\x66\x16\x06\x82\x50\x59\x47\xb3\x70\x2c\xac\xed\x81\x3e\x7d\x68\x18\x8b\x85\xc5\x41\xa9\x62\xc4\x36\x60\xe5\xc2\xd9\x31\x4d\xff\x0c\xf2\xd9\x68\xbc\x5d\xa9\x1b\xf5\x0e\x82\x56\x37\x13\xad\x9e\x82\x37\x66\xe7\x14\x38\x79\x95\xd3\xd7\x21\xa7\xb2\xaf\x84\x86\xf1\xe3\x35\xa8\x20\xb6\x0d\x98\x4f\xea\x45\xf2\xd8\x0a\xb7\xa9\x96\x75\x4a\xcd\x08\x31\xe6\x62\x7e\x8c\x80\x4d\xa0\x82\x4e\xac\xea\xb3\x13\x58\xc5\xd9\x81\xba\x9c\x85\x92\x94\x31\x82\xc6\x04\x8e\x55\x2c\x44\x59\x5e\x50\x52\x93\x3b\x52\xd2\x6d\x47\x92\xe6\xdd\x64\x1a\xc9\x2c\x54\xb1\x6b\x9d\x39\x2b\xb7\x51\x93\x9d\xdc\x48\x2a\xf0\x3f\x59\xa9\x4a\xd4\xda\x00\x9c\x8b\x4b\xe0\x7b\x65\x8e\x2f\xe6\xa9\xe6\xc7\x50\x49\x4a\x21\xc5\xb3\x72\xfa\x5a\x36\x87\x66\xf3\x47\xe0\xff\x06\x00\x00\xff\xff\xf1\xc3\x8c\xfa\x1b\x16\x00\x00")

func tmplIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_tmplIndexHtml,
		"tmpl/index.html",
	)
}

func tmplIndexHtml() (*asset, error) {
	bytes, err := tmplIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tmpl/index.html", size: 5659, mode: os.FileMode(420), modTime: time.Unix(1484959620, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _tmplJsScriptJs = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x9c\x56\x6d\x73\xe3\x34\x10\xfe\xde\x5f\xa1\x31\x9d\x8b\xdc\x34\x71\x8e\xe9\x17\xdc\x97\x00\xe5\x98\x0e\xc3\xc1\x0c\xe5\x0b\xc4\xe9\x8d\x62\x6f\x62\x73\xb6\x64\x24\x99\xbb\xd0\xe4\xbf\xb3\xb2\xec\xf8\xa5\x6e\xda\x23\x33\x69\x92\xdd\x67\x1f\xef\x9b\x1e\x75\x5d\xf0\x50\x27\x82\x93\xdf\xe0\xef\x02\x94\xfe\x59\x6c\x36\x20\xa9\x4b\x1e\x4f\x08\xbe\x3c\x2f\x14\x59\x9e\xa4\x40\x74\x8c\x6f\xc8\xf2\x94\x69\x28\x5d\xff\x30\x49\x94\x28\x64\x08\xe4\x9a\x9c\x52\xe7\x2b\x69\x19\x26\xc0\xb5\xdc\x3a\xee\x34\xd6\x59\x4a\xdd\xcb\x03\xb8\x8e\x46\xf8\x1d\xe3\x51\x0a\x2b\x26\xd5\xb4\xe2\xa7\x96\x0a\xe1\x25\x5e\xc7\x89\x9a\xf2\x22\x43\xec\xec\xb2\xb1\xe8\x6d\x6e\xc2\x9d\x8c\x85\x09\xd7\x42\xc5\x4e\xcb\x19\x8a\x54\x48\xe3\x95\x10\x39\x15\xcf\xa1\xbc\x0d\xe8\xdb\x42\x4a\xcc\xed\xf7\x24\x83\x43\x7d\x75\x6e\xa1\xf5\x45\x36\x3d\x0e\x9f\xc8\x0f\xf8\xb5\xce\xbe\x46\x19\xb7\xc6\x70\x84\xb4\x02\xa6\xc8\x6d\xd1\x64\x4c\x1c\xcf\x21\xe3\x43\x90\x79\xd1\x1e\xf4\xbd\xe0\x3a\x2e\xb1\x6f\x87\x03\x7a\xf8\x1f\x8b\x34\xfd\x03\x98\xb4\xf4\xe4\x5b\xf2\x02\xfe\x0e\x1b\xa9\x2c\xd8\x7f\x01\xfa\x3e\xe1\x85\x86\x57\x82\xef\x21\x14\x3c\x52\xed\x9e\x48\xd0\x85\xe4\x87\xb6\x58\xc7\x7e\xa0\xf3\x66\x32\xdf\x6f\xef\x35\xd3\x85\xa2\xaa\xfc\xb8\x15\x11\xb4\xa7\x90\xac\x49\xcb\x43\xae\xaf\xc9\xd7\xb3\x19\xd9\xed\x48\xdf\xf8\xb6\x1d\xd5\xca\xc2\x51\x45\x18\x82\x52\x4e\x93\xdf\x9e\x40\xaa\x60\x80\xfa\x62\x76\xf1\x1c\x4b\xc4\x38\x1e\x80\xa7\x24\xc3\xe8\x4f\x4c\xf2\x84\x6f\xda\xf0\xba\x09\x03\x5d\xe0\x1a\x3e\x6b\x8a\x6b\x7d\x8e\x2d\x63\xfd\x1d\xb4\x29\xe2\x6e\x19\xe7\x54\x82\xca\x05\x57\x30\x6d\x32\xef\x2e\x63\x06\x3a\x16\x51\x03\x2f\x0f\xdf\xd4\x5a\xbb\xc8\x9c\xe9\xb8\x8f\x33\xb6\x2e\xaa\xf2\x20\xf0\xa7\xfb\x5f\x7f\xc1\xc7\x4a\x2c\x2c\x59\x6f\x69\x3b\xee\x9c\x14\x3c\x82\x75\xc2\x21\x3a\x27\x17\x6e\x9f\xc1\xa6\xfc\x2c\x85\x75\x1f\xe5\x48\xc5\xe6\xf9\xf0\x22\x3d\x9e\x40\xad\x00\xcf\xac\x9c\x7b\xd2\xdf\xdc\xc7\xba\xea\x0f\x38\x14\x9f\x98\xc9\xd4\x16\xbf\xae\xb8\x4e\xdb\x3f\xe4\x2f\x31\x49\x1f\xdf\xe7\x75\xb0\x39\x00\xfe\x13\x89\x39\x50\x7d\x08\x71\x78\x7e\x35\xdf\xc6\x6a\x47\xe5\x57\x83\x6c\xec\x66\x34\x7e\x39\xb4\x36\x03\x96\xe3\x57\x05\xee\x3b\x07\xad\xd4\x3e\xcc\xe6\x9d\x91\x5d\x2c\xfe\xb0\x72\xb4\xbf\x64\xb5\xa8\x8e\xc7\xfd\xb6\x95\x8b\x59\x37\xce\x6e\x69\x0d\xb6\xab\xda\x0d\x30\xca\x8e\xe8\x5a\xd0\x69\x45\xd0\x42\xb5\xef\x03\xcd\x56\xe6\xfe\x58\x89\xc8\xdc\x0a\xb9\x84\x1c\x78\x44\x0d\x47\x15\x80\xf5\x60\x2d\x2f\x49\xf5\xcb\x32\xfd\xbf\x24\xfa\x8b\xe4\xf9\x4b\xa4\xf9\x95\xb2\xfc\x6a\x49\x3e\x2e\xc7\x4f\xa4\xb8\xdd\x51\xb5\xe5\x9a\x7d\xbe\x4b\x36\x71\x8a\x6f\x4d\xff\x52\x82\xd7\x6d\x35\xdf\xb1\x59\xe6\x03\x8f\x18\xce\x33\x04\xea\xbd\xf1\x36\xe7\x64\xf4\x86\x65\xf9\xe5\xc8\x6d\xcc\x57\xd6\x9c\xea\x8e\xf5\xc6\x5a\x37\xc6\xda\x49\xa6\xcb\x49\x1d\x1a\x04\xc5\x82\x4d\xfe\xfd\x6e\xf2\xe7\x6c\xf2\xcd\xf2\xf1\x62\xbf\x0b\x82\xc5\x43\xb1\xdc\x2d\x1e\x82\xc0\x59\xba\x67\x08\x51\x67\xbe\x3b\xdf\x05\x2b\xaa\x65\x01\xbb\x35\x43\xed\xdd\x71\xec\xb3\x1b\xac\x76\x93\x79\x10\x8d\xe9\xdc\x0f\xa6\x41\x74\xe6\xce\xf1\xdb\x02\xde\x2d\x17\xe3\x60\xb2\x34\x1e\x77\xee\x9a\x54\x9a\x03\x90\x31\x1d\xc6\x4f\xae\xfa\xd4\x68\xec\x08\x57\x7b\x05\x72\x74\xd9\xb9\x80\xbc\x07\xc7\x9b\x62\xfb\x75\x15\xda\xbf\x24\x4a\x8c\x7f\x7a\x14\x53\xce\xca\x3e\xe3\x23\x6c\x5b\x0f\x28\x57\x7d\xe8\x32\x69\x45\x58\xcd\xeb\x07\x0d\x5d\x68\x5e\xd3\x9f\xa3\xe9\x54\xc4\x2b\x21\x52\x60\x7c\x34\x78\x39\x7a\xa6\xc1\xaf\x61\x31\xb8\x51\xff\xae\x6b\x4d\x7c\x74\xa5\x72\xc6\x11\xcd\x94\xba\x76\x46\xb8\xc3\x26\x70\x4c\x46\xce\x8d\xf9\x51\x72\x9b\x9f\x57\x9e\xc1\xdd\x54\x54\x7b\xb7\x5c\xd7\x93\x53\x1a\x89\xb0\xc8\x70\xc9\xcd\x76\xb1\x68\x4b\x9b\x51\xb6\x13\x32\xe2\xb2\xd2\xfc\x16\x0b\x92\xf8\xff\x2a\x8a\x4a\x98\x26\xe1\xc7\x23\xe8\x61\x29\x42\xf9\xd2\xdb\xfa\x00\x99\x24\x4e\xcc\x9f\xff\x02\x00\x00\xff\xff\x5f\x8d\xbf\x9f\x12\x0b\x00\x00")

func tmplJsScriptJsBytes() ([]byte, error) {
	return bindataRead(
		_tmplJsScriptJs,
		"tmpl/js/script.js",
	)
}

func tmplJsScriptJs() (*asset, error) {
	bytes, err := tmplJsScriptJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tmpl/js/script.js", size: 2834, mode: os.FileMode(420), modTime: time.Unix(1484958529, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"tmpl/css/style.css": tmplCssStyleCss,
	"tmpl/index.html": tmplIndexHtml,
	"tmpl/js/script.js": tmplJsScriptJs,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"tmpl": &bintree{nil, map[string]*bintree{
		"css": &bintree{nil, map[string]*bintree{
			"style.css": &bintree{tmplCssStyleCss, map[string]*bintree{}},
		}},
		"index.html": &bintree{tmplIndexHtml, map[string]*bintree{}},
		"js": &bintree{nil, map[string]*bintree{
			"script.js": &bintree{tmplJsScriptJs, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

