// Code generated by go-bindata.
// sources:
// wski18n/resources/de_DE.all.json
// wski18n/resources/en_US.all.json
// wski18n/resources/es_ES.all.json
// wski18n/resources/fr_FR.all.json
// wski18n/resources/it_IT.all.json
// wski18n/resources/ja_JA.all.json
// wski18n/resources/ko_KR.all.json
// wski18n/resources/pt_BR.all.json
// wski18n/resources/zh_Hans.all.json
// wski18n/resources/zh_Hant.all.json
// DO NOT EDIT!

package wski18n

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

var _wski18nResourcesDe_deAllJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func wski18nResourcesDe_deAllJsonBytes() ([]byte, error) {
	return bindataRead(
		_wski18nResourcesDe_deAllJson,
		"wski18n/resources/de_DE.all.json",
	)
}

func wski18nResourcesDe_deAllJson() (*asset, error) {
	bytes, err := wski18nResourcesDe_deAllJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "wski18n/resources/de_DE.all.json", size: 0, mode: os.FileMode(420), modTime: time.Unix(1520374115, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _wski18nResourcesEn_usAllJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x5b\x5f\x8f\xdb\x38\x92\x7f\xcf\xa7\x28\x0c\x0e\x98\x19\xc0\x71\xcf\xec\xe1\x80\xc3\x00\x79\xc8\x25\x99\xd9\xbe\x4d\x26\x41\x27\xd9\x60\x91\x0b\x14\x5a\x2a\xdb\x5c\x4b\xa4\x96\xa4\xec\x78\x1a\xfe\xee\x87\x2a\x92\x92\xec\x36\x25\xd9\xc9\xdc\xe5\x25\x6e\x93\xac\xfa\x55\xb1\x58\xac\x3f\xf4\xc7\x47\x00\xf7\x8f\x00\x00\xbe\x93\xc5\x77\xbf\xc0\x77\x95\x5d\x65\xb5\xc1\xa5\xfc\x92\xa1\x31\xda\x7c\x37\xf3\xa3\xce\x08\x65\x4b\xe1\xa4\x56\x34\xed\x05\x8f\x3d\x02\x38\xcc\x06\x28\x48\xb5\xd4\x09\x02\xb7\x34\x34\xb6\xde\x36\x79\x8e\xd6\x26\x48\xbc\x0d\xa3\x63\x54\x76\xc2\x28\xa9\x56\x09\x2a\x1f\xc2\x68\x92\x4a\x5e\x15\x59\x81\x36\xcf\x4a\xad\x56\x99\xc1\x5a\x1b\x97\xa0\x75\xc7\x83\x16\xb4\x82\x02\xeb\x52\xef\xb1\x00\x54\x4e\x3a\x89\x16\x7e\x90\x73\x9c\xcf\xe0\x8d\xc8\x37\x62\x85\x76\x06\x4f\x73\x5a\x67\x67\xf0\xce\xc8\xd5\x0a\x8d\x9d\xc1\x5d\x53\xd2\x08\xba\x7c\xfe\x23\x08\x0b\x3b\x2c\x4b\xfa\xdf\x60\x8e\xca\xf1\x8a\x2d\x73\xb3\x20\x15\xb8\x35\x82\xad\x31\x97\x4b\x89\x05\x28\x51\xa1\xad\x45\x8e\xf3\xc9\xb2\x68\x9d\x92\xe4\x29\x38\xad\x4b\x70\x3a\x08\x32\x83\x46\xf9\x4f\x20\x54\x01\x76\xaf\x72\xd0\x35\xaa\xdd\x5a\xda\x0d\xd4\x41\x26\x68\xac\x54\x2b\x10\x50\x09\x25\x97\x68\x1d\x4f\xd6\x35\x51\x15\x65\x20\x55\x91\x24\x4b\x59\xb6\xd3\xff\xf1\xf4\xd5\xcb\x29\x98\xed\x5a\x1b\x37\xbc\x01\x6f\x8c\xde\xca\x02\x2d\x08\xb0\x4d\x55\x09\xb3\x07\x3f\x1f\xf4\x12\x76\x6b\xe1\xbe\xb7\xb0\x40\xec\x6d\xcf\xd7\xa9\x31\x40\x1a\xd5\xa3\x45\x47\xba\x5c\x63\x59\x07\xd6\xb0\xd7\x8d\x99\xa4\x42\x52\xd5\x74\x2c\x5b\x34\x96\x78\xa7\xf4\x23\x95\x63\x81\xc3\x3c\x50\x4d\xb5\x40\xc3\xea\xb1\x1b\x0f\x6d\x32\x2f\xb2\x82\x51\xfb\x61\x53\x61\x61\x5f\xd7\xa8\x3e\x1c\x0b\xbb\x40\xb7\xa3\xed\xc8\x4b\x49\x56\xc1\xa6\x85\x66\x8b\x66\xb2\x0d\x4f\xc7\xd0\xb3\x3e\xe2\x13\xed\x99\xbf\xd0\xcb\xff\x4b\x6b\x5e\x96\x62\x95\x89\x5a\x66\x6b\x6d\x53\x86\xe3\xa1\x3c\x7d\x73\x0b\x9f\xff\xfa\xfa\xed\xbb\xcf\x13\x29\x0e\x6f\x7f\x8f\xe8\xdf\x5f\xdc\xbd\xbd\x7d\xfd\xfb\x24\xba\x8d\x5b\x67\x1b\xdc\x27\x88\xd2\xb0\x36\xf2\x0f\xfe\x02\x3e\xff\xed\xc5\x3f\xa6\x10\xcd\xd1\xb8\x8c\xf4\x96\xa0\x5a\x0b\xb7\xa6\x6d\x21\x5b\x9d\xd3\x64\x56\xf2\x14\xc2\x5a\x2d\x65\xca\xd9\xfb\x41\x26\x05\x3f\x14\xb8\x14\x4d\xe9\x40\x5a\xf8\xb7\xbf\xbe\x7e\xf5\xe2\x66\xbe\xb3\x9b\xda\xe8\xda\xfe\x38\x45\x2b\x65\xa9\x77\x59\xa0\x91\xba\xa2\x78\x12\xb4\x93\xc6\xa9\x76\x46\x35\xa4\x97\xd6\x2d\xb7\xd6\x37\x81\x74\x6d\x70\x2b\x71\x97\xa0\x6b\xd7\x0c\x34\x12\xbd\x39\x3a\x1e\x75\x29\xd4\x04\x0e\x1b\xdc\x4f\xde\xd2\x0d\xee\xa7\x02\xf7\x9a\xae\x84\x12\x2b\x2c\x06\x15\x5d\x1b\xfd\x4f\xcc\x5d\x77\xe7\x3a\x0d\x0b\x84\x4a\x98\x0d\x16\x10\x29\x4c\x51\x15\xd3\xc9\xe8\x2e\x48\x09\x13\x58\xf1\x94\x71\x8a\xd1\x85\x8c\xec\xea\x91\xd3\x9f\x40\xb6\xbd\xac\x12\x74\xbb\xf1\xc9\x42\x8f\x20\xf4\xee\xb9\x44\x6b\xa3\xb6\x27\x90\xb6\xce\xc8\x24\x65\xbf\x75\x8d\x45\x43\x07\x45\x2a\x2c\xc0\x34\xca\xc9\xaa\xbd\xa4\x26\x70\x70\x26\xad\x04\x1e\x03\xdd\xb8\xba\x99\x02\xd6\x9b\xdb\x16\xcd\x42\xdb\x14\xc9\x30\x7a\x29\xd1\x5a\x18\x51\x25\x15\x6c\x44\x85\x0e\x0d\x6c\x45\xd9\x20\x07\x78\xe4\x4c\xe1\xef\x4f\x5f\xbe\x7f\xf1\x19\x96\xda\x54\xe2\x42\x56\x43\xa7\xf1\xf3\xaf\xb7\x2f\x5f\x7c\x86\x5c\x2b\x27\x24\x45\xc0\x70\x0e\xc1\x7f\xbf\x7d\xfd\xfb\x38\x6b\xf6\xaa\x59\x25\x2d\xdd\x80\x7c\x5f\xa4\xaf\x8b\x77\x6b\x04\x9a\x41\x67\x34\xf7\x77\x06\xf9\x02\x69\x41\x69\x07\x9e\x54\x63\xb0\x98\xff\xcf\xd0\xbe\x9f\x70\xac\xe5\xc0\x55\x4a\x1c\xe9\xce\xa3\x29\x5f\xc7\x67\xec\xb8\x11\xa7\x76\xce\x75\xac\x82\x28\x43\xa9\xd3\xa9\x3c\x1f\xef\xef\xe7\xf4\xf9\x70\xf8\x34\x83\xa5\xd1\x15\xdc\xdf\xcf\xad\x6e\x4c\x8e\x87\xc3\x24\x9e\x7e\xc3\xc6\x78\xd2\xb4\xb8\x57\x16\xdd\x75\xbc\x5a\xf5\x8c\x71\x3b\xd2\x23\x89\xd8\x7e\x71\xbd\x9c\xb5\x5c\xed\x32\xc1\x59\x63\xe6\xf4\x06\xd5\xa8\xc8\xb4\x02\xfc\x0a\xe0\x15\xd7\x09\xdf\xa8\x4a\x18\xbb\x16\x65\x56\xea\x5c\x94\x09\x8e\xef\xe3\xac\x5e\xa8\x1c\x5c\xa1\xf5\xfc\x78\x75\x38\x9e\x13\x19\x2a\x74\x3b\x6d\x36\x57\xb3\x94\xca\xa1\x51\xe8\x40\x38\x12\xb7\x31\xe5\x88\xac\x5d\xdc\x90\xe5\x42\xe5\x58\x96\xc9\x5b\xfb\xf5\xdf\xe6\xf0\xcc\xcf\x21\x07\xd4\xad\x9c\xca\x60\x29\x64\x9a\xfa\xf3\x2e\x80\x29\x64\x11\xce\x62\x55\x97\xe8\x10\x6c\x43\x5b\xba\x6c\xca\x72\x3f\x87\xbb\x46\xc1\xe7\x36\xfd\x69\x33\x83\xcf\x74\xdf\x19\xac\xf4\x16\xc9\x37\x3a\x29\xca\x72\xdf\x65\x8e\xc2\x5a\x74\xc3\xbb\xd0\x43\xea\xd3\xd0\xcc\x3a\xe1\x9a\x54\xb4\xf8\xf8\xf1\xe3\xc7\x4f\x9e\x3c\x79\xd2\xdb\x8b\x9e\x0c\x6f\x79\x29\xd0\x04\x9a\x38\x89\x2b\x17\x50\xb0\x98\xa2\xa2\xa8\x9a\x02\x42\xd5\xc5\x2b\x67\xd8\xc8\xae\xdf\xeb\xfe\xda\xe9\x4c\x06\xf7\xfb\x7d\x3f\x64\x1d\xdc\xf1\xc9\xfc\xc6\xf4\x77\xc4\xf2\x0a\x0d\xe6\xba\xaa\x84\x2a\x32\x4e\x1d\xf9\xb6\x26\x2f\x97\x09\x97\x51\xbc\x95\x60\x7a\x7f\x3f\xcf\xab\xe2\x70\x08\x09\xe7\xfd\xfd\x9c\x16\xba\x7d\x8d\x87\x03\x7b\x4a\x5a\x7b\x38\x7c\x9a\xcf\x07\x79\x73\x90\xbc\xcf\xa2\x3d\x8f\x14\xdb\xee\xef\x29\x64\x0f\x0c\x08\xe4\xe1\xf0\x09\xd6\x22\x94\x53\xfa\x02\xb7\x27\x64\x3a\xf7\x74\x75\xee\x79\x1c\x87\xb3\x00\xe6\xf3\x81\x54\x3b\xb0\x88\x1b\xfa\x2d\x45\xec\x68\x4e\x11\x32\xce\x4e\x8b\xf9\xbe\x9b\x71\x56\xd0\x41\x39\x0b\xac\x51\x15\xa8\xf2\x4b\xd4\xd9\x2d\xba\x9e\x4f\x77\x44\x92\x3a\x7d\x7e\x96\xcd\xd7\x18\xce\x79\x14\xe4\x18\x1a\x93\x8a\xcb\x9e\x1f\x55\x7a\xce\x8b\xfe\xff\x78\x47\x44\x79\x2e\xb3\x93\xaf\xdb\xc1\x87\x6e\xee\xdb\xec\xe1\xc4\x93\x91\x42\x32\xbc\x8f\xef\x4f\x6a\x76\xd7\xec\xe4\x10\xaa\x50\x21\xb8\xf6\xce\x61\x44\xfe\x06\x68\x2b\x10\x43\x58\xa0\x68\x0c\xed\x64\x60\xdb\x8f\x7f\xfe\x3c\x7b\x8b\x32\x2e\x75\xa3\x8a\x2c\xe0\x0d\x9e\x2a\x69\x00\x25\xba\xa4\x0f\xde\xad\x65\xbe\x86\x1d\x77\x29\x08\x57\xe1\xe3\x46\xb7\x46\xc8\x1b\x63\x48\x31\x51\xc0\x58\x34\xe1\x4b\xca\x7f\x26\x0a\xc2\xb2\x2c\xa4\xbf\xc9\x61\x41\xa8\xa9\x65\xa1\x58\x9b\xaa\x77\xfb\x51\x4e\x26\xa0\x57\xef\x33\xc8\x75\x8c\x62\x06\xa2\xec\x87\xbe\xed\xb6\x11\x0e\xd3\xae\x08\x4c\x40\x18\x6c\x75\x7d\xd3\x59\x3a\x14\xd2\x60\xee\x82\xf5\x1b\x5f\xed\x1e\xeb\x23\xbc\xb8\xbb\x7b\x7d\xf7\x36\x81\xfb\xc9\xe9\x3f\xf0\xd3\xe1\xc1\xc0\x93\x27\x03\xd7\x8f\x31\xc7\x07\x6d\xa3\xf4\x4e\x65\x14\x29\x8c\x1f\x75\x9a\x45\xaa\x0a\xab\xe6\xd0\x35\x08\x40\xab\x72\x0f\xb6\xa9\x7d\xb7\xeb\x86\xcb\xca\x73\xbb\xb7\x0e\x2b\x58\x48\x55\x48\xb5\xb2\xa0\x0d\xac\xa4\x5b\x37\x8b\x79\xae\xab\xb6\xa8\x3e\x7c\x5f\x1a\x13\xef\xcc\xdc\xa0\x70\x29\x98\xdc\x7d\x04\x9e\x72\x64\x96\x3b\xe9\xd6\xc0\x6d\x4b\xa8\xd0\x5a\xb1\xc2\x5f\x68\x10\x8d\x39\x1c\xb8\x78\xef\xc7\x72\x5d\xf8\x01\xfa\x30\x92\xcd\xf4\x20\xf9\xb3\x32\x08\xa9\x78\x70\x52\xfe\x24\x48\x4b\xc4\x22\x93\x6a\xab\x37\x29\x40\xbf\xb2\xdb\x22\x77\xe1\xa7\xf1\x81\xa4\x65\xb0\x5b\x73\x03\x2c\x20\x75\xbe\xf9\x18\x86\xfe\x1c\xb4\x1b\xdc\xb7\x35\x14\x8a\x77\x85\xd3\x66\xa8\x3e\xd4\xce\xe1\x72\xc3\xc7\xa8\xcc\x4f\x64\x8f\x81\xce\x28\xcf\x58\x4a\xcd\x94\x76\xde\xd9\x25\x18\xbe\xea\xd7\x5c\xd9\x57\xf3\x6c\xca\x77\xb9\xe8\xd9\x8f\xa8\xc7\x98\x72\xf4\x5e\x49\x5b\x09\x97\xa7\xc2\x77\x12\xb0\x35\x0f\x5a\x50\x30\x8b\x22\xfa\x53\xa9\x4e\x8b\xfb\x7e\x3c\x60\x80\x42\xa3\x2f\x2c\x31\x13\xde\x56\x76\x6f\x34\xa9\xea\x11\x39\xaa\x25\xfb\xd1\x28\xc6\xb0\x10\x21\xff\x27\xf3\x12\xa5\x4c\xa9\xed\xd6\x8f\xd2\x31\x0f\x5b\xd2\x96\x6d\x89\x57\xf8\x4c\x58\xba\xde\xea\x11\x2a\x6d\x18\xbb\xe0\x2e\x38\xaf\xf1\x1f\xa7\xe8\x39\x42\x1c\x51\xf5\xdd\x25\x80\x4e\xf4\xca\x47\xc1\x23\xfa\xde\x82\xaf\xf2\x78\x55\xe2\x17\x87\xca\x46\xd0\xf8\x85\xef\x30\x12\xe7\x6b\x44\xb1\xd9\x0a\x53\x05\xcc\xee\x28\xaf\xd0\x77\x6f\x83\xef\xed\x4a\xe5\xa1\x58\xd3\xdd\x64\x74\xbf\xc9\xbc\x77\x7c\x27\xeb\xd4\x43\xcf\xbc\xc4\x7c\x7a\x5a\x6e\x09\x7c\x47\x02\x73\x5c\x48\x6a\xec\xb4\x2c\xd4\xbe\xb5\x0d\x72\x22\xbd\x6d\x1f\xd5\x6b\x28\xa2\xb6\x10\x46\xc5\x68\x4c\x79\xb9\xe5\xfa\xc2\x56\x48\xa1\xdf\xdf\xbd\x64\x04\x5c\xea\xe2\xa3\xf4\xf1\x28\xc7\xfe\xe4\x5b\xf2\x53\x80\x54\xa2\x5c\x6a\x53\x25\x35\xf7\x2a\x8e\x0f\x21\x98\xc3\x3b\xb3\x07\xb1\x12\x52\x8d\xa5\xf4\xc6\x64\xff\xb4\x5a\xb5\xce\x36\xaf\x8a\x81\xce\x2d\x17\xf7\xa5\xaa\x1b\x07\x85\x70\x02\x5e\x05\x6d\x7c\x9f\x57\xc5\xf7\xe4\x7a\x87\x39\x89\x5a\x76\x15\x78\x6f\x34\xda\x64\x16\xff\xd5\xa0\x4a\x96\xc8\xfd\xa3\x97\x9b\xb7\x61\xd6\xf1\x61\xe9\xf9\x77\x6f\xcf\x47\x3e\x6c\xc6\x55\x6f\x5e\x50\x4b\x9a\x9d\x0b\xe5\x43\x91\x05\xfa\x60\x00\x0b\x58\x08\x8b\x05\x68\xd5\x33\xb2\x9b\x08\xe9\x0c\xcd\x39\xbc\x29\x51\x58\x84\xa6\x2e\x84\xc3\x13\xa7\xc9\x97\x67\x5e\x36\xc5\x29\x4e\x61\x41\xc0\x0e\x17\xa7\x1c\x46\x77\x27\xe8\x69\xd8\x40\x9f\x9e\xf1\x23\xa4\x9a\xb0\x6a\x0e\xb7\xce\x67\x5f\xda\xad\xf9\x2e\xe6\x53\xb5\x6c\x54\x38\x53\xf1\xe0\xcd\xbc\x76\xb4\xc2\xd0\x76\xad\x88\x0a\x7e\xa9\x31\x9f\x72\x92\x02\xd6\xb8\xc5\xd1\x3f\x90\x63\xcc\x88\xeb\x57\xa2\x67\xe0\x9d\x93\x20\xb2\xba\x71\x7d\x67\x31\x87\x0f\x9d\x13\x8e\xae\x82\x96\xcd\x5a\x77\x42\x06\x13\x83\x85\x91\x6b\x2d\x88\x13\xd5\x94\x51\xb6\xe2\x30\x2b\xa4\x99\xe4\xe4\xce\x8a\x45\x72\xb4\x7a\xaf\xb5\x54\x3e\xa4\xf2\x29\x9a\xc3\x90\x18\x50\x20\xd3\x1d\xe7\x19\xa5\x80\x51\x2a\xcb\x39\xc5\xb1\x87\x1b\x16\x23\x17\x94\xb0\x8b\x2d\x66\x85\xce\x37\x98\x7a\xa0\xf7\x4c\x28\xa6\x2a\xb6\x08\xcf\x79\x22\xc8\x8a\x03\xf0\x91\xc0\x52\x96\x98\x89\xd2\xa0\x28\xf6\x19\x7e\x91\x36\xf9\xb6\xe1\x57\x3a\x21\x61\x26\xf8\x99\x23\xb4\x0b\xb9\x5c\x22\x25\x84\x5d\x56\x22\xd1\x7a\x83\xb2\x14\x39\x95\x62\x81\xa9\xe6\xc8\x6b\x85\x40\x76\x58\xe2\x69\xda\xdf\xfd\x19\xb7\xc4\xed\x34\xb4\xcc\xb8\x69\xe2\x75\x4d\xb3\xe3\x5f\xde\xb1\xae\xa5\x85\x8d\x54\x05\x1d\x90\x60\x8b\xa1\x29\xf9\xe0\xe2\x39\xf1\x14\xe4\x5f\x7a\x40\x18\xfa\x19\x38\xe1\x7d\xd9\x03\xbf\xc2\xc6\xc2\x0d\x75\x8a\xdd\x22\x28\x88\x69\x0d\xb2\x0c\x16\x6b\x61\xe8\x0f\xa6\x6e\x7d\xcc\x74\x5e\xb6\x69\xc6\x1f\x0e\x59\x46\x22\x5f\x6a\xe7\x4a\x7b\x4d\x59\x74\x97\x31\xbb\xd4\x57\x04\x66\xbd\xf3\x3e\xc2\x2f\x7a\xdf\x6c\x2d\xb6\xe4\xa9\xd8\x96\x7c\x21\xdd\x06\x30\xa9\x27\xa4\xfd\x6b\x28\x92\x09\xfe\x2a\x9a\x76\x7c\x94\x40\x3e\x5f\x45\x67\xe4\x13\x7d\x0e\xc5\x68\xff\x42\x76\x3b\x8f\x6f\x3a\xc3\x4b\x36\x4f\xcf\xf2\x45\x45\xc6\xb4\xa6\xd3\xc8\x0b\x38\x62\x97\x0a\x44\xb4\xe9\x48\x61\xe4\xf0\x6b\xb5\x2c\x65\x4e\x5e\x26\x0b\x89\x1b\x49\x68\xb4\xb5\xb1\x12\x92\x3a\xae\xbd\xf3\x13\x53\x3e\x12\x3a\x7c\x0e\x32\x47\x59\x39\xf8\xad\x9a\xd2\xc9\xba\xf4\x59\xa3\x3f\x3c\xf4\x29\x44\x24\x9e\x39\xbb\xaf\x78\xf7\x9e\x94\x41\x5c\xbf\x8b\x3b\x03\xe9\xfc\x89\xaa\xb5\xb5\x72\xe1\x4f\x01\x2b\x24\x0a\xe2\xb9\x76\xea\x59\x50\x5c\xd2\x5a\x3a\x83\x78\x70\x08\x83\x24\xcc\xe6\x41\xd2\x73\x81\x32\x4d\x53\xe2\x15\x9a\xa4\x65\x21\xbb\x28\xf1\x9c\x0e\x3b\xfc\xd1\xdf\x9f\x04\x12\xaa\xb8\xa1\x43\x1d\x55\x70\xbc\x25\x73\xff\x20\xf8\x5b\x28\x99\x05\x3c\xa7\x61\x61\xad\xce\x25\x93\x3e\x8f\xf8\x26\x82\x3b\x55\x3e\x0b\x7f\x95\xe6\x85\xe9\xde\x54\x70\x33\x3b\xf9\x82\x33\x34\xc8\xa0\x94\x0a\x41\x98\x55\xc3\x49\x31\xa9\xd0\xac\x0e\x87\x7e\xbc\xc8\x74\x66\x50\x7b\x88\xde\x97\xef\x59\x1f\x3c\x72\x01\xa2\x0d\xee\xbf\x19\xaa\x0d\xee\x6f\x98\x16\xd4\x42\x9a\x07\xf0\x8e\x87\xd9\xbf\xe3\x17\x51\xd5\x14\xec\xb6\xe4\x36\xb8\x9f\x24\x43\x08\xb0\xc6\x9f\xfe\xa4\x04\xf8\x21\xb2\xfc\x91\x7d\x70\xa0\xe7\xdf\x05\xf9\x8b\xab\x2d\x85\xcc\x7c\x41\xb2\x97\x5e\x46\xe3\x88\xa2\x09\xf0\xab\x39\xc9\xe8\x48\x8c\xd5\x1e\xf0\x5f\x8d\x34\x5c\xdb\xaa\x1b\x67\x27\x59\xc9\x5d\x58\xe3\x53\x19\x7f\x5a\x8e\xac\xc2\x02\x6e\x51\x81\x58\x3a\x34\x20\xea\xba\xe4\xfe\x09\x3f\x6c\xa8\xb5\xa7\x13\x7a\xa9\xa8\xb6\x73\xd8\x0a\x23\xc5\xa2\xc4\xce\xe0\x2d\xba\x96\xe2\xf1\x94\x78\x80\x7d\x16\xd5\xbd\x9b\x8a\xa7\xe1\xe6\xb4\x94\xa3\x0d\x25\x27\xcf\x5e\xde\xf2\x66\x2f\x75\x59\xea\x9d\x47\x43\xd8\x59\x9f\xfe\xe3\xe1\x30\x9e\x7d\xad\x84\xc3\x9d\xd8\x67\x94\xf4\x70\xc7\x78\x2c\xb1\x78\x73\x0b\xbf\xf9\x35\x9c\x28\x75\x05\x2e\x51\x4b\xfa\x22\xd6\x98\xce\x84\xeb\x3c\xb5\x7d\x22\x66\x43\xd9\xff\x34\x4a\x0a\x29\x87\x41\x62\xba\x0d\x0c\xda\x4a\xf1\x09\x8d\x68\x0b\xa7\x22\x7e\x78\x7a\xf7\xfb\xed\xef\xbf\x4d\x2f\x8e\xc7\x05\x97\x95\xc7\x77\xc2\xa8\xb6\x03\x6f\xd0\x25\x4b\x92\x77\x34\x46\x9b\xf4\x31\xb6\xde\x3f\x05\x63\xe2\x72\xe8\x2f\xbe\x5e\x41\xa7\xe0\xd3\x50\x4e\x15\xf8\xf1\x53\xa4\x8b\x2b\x14\xfd\x97\xcb\xbd\x8a\x24\x14\xe8\xc6\xb3\x39\xe6\x4c\x6e\xad\xc0\xda\x60\x4e\xce\x3e\x33\x58\x97\x22\x4f\xa6\x3b\xef\xd6\x9e\x8f\x2e\x8b\x50\x7b\xe5\x97\x5f\x3e\x9a\x3d\x7e\x72\xb0\x93\x65\x09\x56\x6b\x45\x79\x78\xc7\xa1\x75\x76\x8d\xf5\xd1\x32\x37\x8d\x70\x77\x44\xce\x3a\x14\x13\xb1\x07\x4d\x5c\x53\x36\xb6\x6b\xdd\x94\x05\xc1\xa3\xe0\x15\xde\x5b\xdf\x3f\xf5\xcd\x1d\xef\x7f\x69\x36\x7f\x1a\x7f\x38\xd1\x22\xe2\xf9\x23\x5b\x49\xb8\x3c\x07\x3a\xef\x0f\xcb\xd9\x74\x7a\xfc\x41\xbb\x80\x25\xe7\xab\x62\x3b\xb8\x79\x63\x4c\x79\x7d\xdc\xd0\xd8\xa8\x8b\xbf\x0a\xe9\xff\x1c\x64\x1c\x58\x29\x2b\xe9\x32\xb9\x52\xda\x24\x21\x45\x93\x0e\xf1\x33\x2f\xf1\xf9\x18\x7d\x3a\x2d\x59\x93\xff\xf1\xe4\xa6\x72\xcf\xd7\x42\xad\x90\x7c\x72\x02\xc0\xcb\x96\x63\x5b\x23\xb7\x51\xee\x72\xef\x7b\xb4\x2d\x8d\x39\xdc\x12\x7b\xa9\x56\x53\x6c\x81\x11\xd8\xac\xd4\xab\xcc\xca\x3f\x52\x00\x4a\xbd\x7a\x2b\xff\xe0\x52\x8c\x5f\x70\x24\x71\x67\xa2\x42\xf1\xd5\x44\x61\x6d\xfc\x79\xcc\x4f\x9c\x4f\xfc\xfc\xd3\x64\x28\x15\x56\xda\xec\x87\xd0\xf8\x19\xd7\x02\xfa\xf9\x2f\xff\xc9\x90\xfe\xe3\xe7\xbf\x4c\xc6\xe4\x64\x85\xba\x49\xd5\xb8\xc3\xe8\x55\x60\x7e\xf2\xfa\xf9\xf7\x9f\xe8\xdf\x38\x1e\x6e\x57\x66\xb5\xd1\x35\x1a\x27\x93\x61\x7e\xf4\x80\x3d\x7f\xe5\x9b\xdc\xce\x48\x6c\xdb\xdc\xbe\xf7\xd9\x11\x8b\xed\xf0\xf3\x3e\x31\xba\xc4\x42\xb3\xc1\x91\x67\x94\x0e\x74\xe3\xac\x2c\x78\x23\xde\x19\xb1\x95\x16\x16\x8d\x2c\x8b\xe1\x5e\x29\x8b\xe2\xdd\x81\x21\xb3\x9d\xe4\x0a\x5a\xeb\x3f\x72\x08\xea\xc4\xa1\x07\x6d\x73\x07\x98\xf2\x11\xff\x6d\x54\xf7\xfd\xfd\xbc\x92\x2a\xf4\x03\xe9\x0f\x91\x8f\x74\x17\x18\x6a\x2c\x1f\xfa\x43\x96\x72\x13\xb1\x63\x13\x66\x51\xfa\x72\xd2\xbc\x39\x53\xe0\x4d\xf6\x67\xae\x6a\xca\x30\xda\xd0\xf2\xe5\x22\xc2\x60\x15\xec\x41\x37\xef\xc8\xc5\x9c\x94\xc7\xba\x78\xb2\xc4\xdc\x81\x50\xda\xad\x43\xf6\x3a\x0e\x29\x66\xa5\xa3\x0d\xcd\x77\x0f\xea\x4d\xfd\x80\x21\x3c\xf8\xc7\x02\x94\x9e\xd6\x95\x67\xee\xbd\x07\x31\xac\x94\x29\x20\xce\x3e\x17\x09\x37\xce\x69\x5c\xbc\x0b\x5d\x23\xdf\x7b\x3d\x57\x35\x9b\xa0\xa1\xde\xcf\x76\x32\xbd\x45\x63\x64\x51\x60\xaa\xf6\x43\x08\xfb\xbf\xe2\xe9\x1e\x34\x75\x4b\x63\xac\xd0\x7f\xaf\x32\x75\xa3\x32\x69\xb3\xba\x59\x94\x32\xf5\xfb\x44\xbf\x2b\x3c\x37\xf6\x3e\xfc\x0f\x95\x28\xda\xe6\x85\x0f\xf2\x6a\x4a\xf0\xbd\x6f\x59\x20\x6c\xa5\x4f\xf1\xe9\x1c\xe6\x82\x3d\x8d\x7f\xa9\x8e\x05\x2c\xf6\x20\xd4\x5e\xab\x81\x1f\xfe\x30\xd6\x58\xaa\xc3\x45\x86\x5f\xf8\x85\xf2\xf0\x35\xfe\xb0\x52\xc7\x4d\x08\x6e\x85\xa8\x82\xfe\x7f\xec\xe9\x3c\xe8\x42\xd0\x41\x20\x55\xee\x70\x31\xf3\x97\x7b\xf8\x2b\x2c\x18\x48\x0c\x3d\xd2\x5e\xb7\x89\xe0\x0e\xd6\xf5\x52\x3d\x08\xb2\xb0\x7e\xe3\x66\x52\x4b\xc9\xa7\x86\xdd\xa2\x39\x3c\xd3\x6a\x4b\xee\x3e\xa4\x04\x1d\x0b\xa7\x8f\xc8\x8f\x9b\xec\xa9\x54\x23\xdd\xb3\xa1\x7a\x65\x27\x5b\x1c\xb8\x50\xba\xb6\x89\x75\x2a\x5f\x9f\x51\x2b\xe1\xa4\x96\x57\x2b\x63\xac\x4d\x18\xb4\xb5\x56\x16\x87\x1e\x21\x9d\x80\xe6\xaa\xd4\x69\xf6\x19\xc6\x63\x9e\xd9\xcb\x5b\x63\x05\xa1\xad\x7c\xad\x9d\xab\xfd\x6f\xe8\x3d\x6b\xbe\xd7\xe6\xf0\x8c\x6e\x18\x7e\xb5\xd0\xff\xde\x5f\xea\x7c\xe5\x84\xaf\x83\xd0\x4c\x85\xee\x93\x0e\x59\xc2\x62\x9f\xbf\xf8\xaf\xf7\xbf\x4d\x4e\x5d\x79\xf6\x65\x79\x6b\xb1\x58\x65\x16\x85\xc9\xd7\x64\x35\xd1\xe9\xb5\xad\xa6\xa4\xe9\x84\x15\xad\xd3\x3b\x6e\x4e\x45\x15\x46\x19\x7d\x70\x30\x12\xfe\x12\x94\xd3\x9b\xe1\x5b\xdf\x0a\x57\xde\x08\x04\xad\xbd\x32\xfd\x63\xc7\x81\xdf\xe9\x3f\x3f\xf3\xe2\x26\x68\xe4\x17\xf8\x95\x11\x74\x3f\x0b\xe7\xc2\x2b\x11\xbb\x14\xc0\xf0\x4f\x2c\x2f\xc7\xd0\x7f\x4f\x19\xdf\xff\x06\x48\x8f\x3e\x3d\xfa\xdf\x00\x00\x00\xff\xff\x8c\xf0\xc2\xc6\xdc\x43\x00\x00")

func wski18nResourcesEn_usAllJsonBytes() ([]byte, error) {
	return bindataRead(
		_wski18nResourcesEn_usAllJson,
		"wski18n/resources/en_US.all.json",
	)
}

func wski18nResourcesEn_usAllJson() (*asset, error) {
	bytes, err := wski18nResourcesEn_usAllJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "wski18n/resources/en_US.all.json", size: 17372, mode: os.FileMode(420), modTime: time.Unix(1528403123, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _wski18nResourcesEs_esAllJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func wski18nResourcesEs_esAllJsonBytes() ([]byte, error) {
	return bindataRead(
		_wski18nResourcesEs_esAllJson,
		"wski18n/resources/es_ES.all.json",
	)
}

func wski18nResourcesEs_esAllJson() (*asset, error) {
	bytes, err := wski18nResourcesEs_esAllJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "wski18n/resources/es_ES.all.json", size: 0, mode: os.FileMode(420), modTime: time.Unix(1520374115, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _wski18nResourcesFr_frAllJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8a\xe6\x52\x50\xa8\xe6\x52\x50\x50\x50\x50\xca\x4c\x51\xb2\x52\x50\x4a\xaa\x2c\x48\x2c\x2e\x56\x48\x4e\x2d\x2a\xc9\x4c\xcb\x4c\x4e\x2c\x49\x55\x48\xce\x48\x4d\xce\xce\xcc\x4b\x57\xd2\x81\x28\x2c\x29\x4a\xcc\x2b\xce\x49\x2c\xc9\xcc\xcf\x03\xe9\x08\xce\xcf\x4d\x55\x40\x12\x53\xc8\xcc\x53\x70\x2b\x4a\xcd\x4b\xce\x50\xe2\x52\x50\xa8\xe5\x8a\xe5\x02\x04\x00\x00\xff\xff\x45\xa4\xe9\x62\x65\x00\x00\x00")

func wski18nResourcesFr_frAllJsonBytes() ([]byte, error) {
	return bindataRead(
		_wski18nResourcesFr_frAllJson,
		"wski18n/resources/fr_FR.all.json",
	)
}

func wski18nResourcesFr_frAllJson() (*asset, error) {
	bytes, err := wski18nResourcesFr_frAllJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "wski18n/resources/fr_FR.all.json", size: 101, mode: os.FileMode(420), modTime: time.Unix(1520374115, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _wski18nResourcesIt_itAllJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func wski18nResourcesIt_itAllJsonBytes() ([]byte, error) {
	return bindataRead(
		_wski18nResourcesIt_itAllJson,
		"wski18n/resources/it_IT.all.json",
	)
}

func wski18nResourcesIt_itAllJson() (*asset, error) {
	bytes, err := wski18nResourcesIt_itAllJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "wski18n/resources/it_IT.all.json", size: 0, mode: os.FileMode(420), modTime: time.Unix(1520374115, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _wski18nResourcesJa_jaAllJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func wski18nResourcesJa_jaAllJsonBytes() ([]byte, error) {
	return bindataRead(
		_wski18nResourcesJa_jaAllJson,
		"wski18n/resources/ja_JA.all.json",
	)
}

func wski18nResourcesJa_jaAllJson() (*asset, error) {
	bytes, err := wski18nResourcesJa_jaAllJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "wski18n/resources/ja_JA.all.json", size: 0, mode: os.FileMode(420), modTime: time.Unix(1520374115, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _wski18nResourcesKo_krAllJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func wski18nResourcesKo_krAllJsonBytes() ([]byte, error) {
	return bindataRead(
		_wski18nResourcesKo_krAllJson,
		"wski18n/resources/ko_KR.all.json",
	)
}

func wski18nResourcesKo_krAllJson() (*asset, error) {
	bytes, err := wski18nResourcesKo_krAllJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "wski18n/resources/ko_KR.all.json", size: 0, mode: os.FileMode(420), modTime: time.Unix(1520374115, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _wski18nResourcesPt_brAllJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func wski18nResourcesPt_brAllJsonBytes() ([]byte, error) {
	return bindataRead(
		_wski18nResourcesPt_brAllJson,
		"wski18n/resources/pt_BR.all.json",
	)
}

func wski18nResourcesPt_brAllJson() (*asset, error) {
	bytes, err := wski18nResourcesPt_brAllJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "wski18n/resources/pt_BR.all.json", size: 0, mode: os.FileMode(420), modTime: time.Unix(1520374115, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _wski18nResourcesZh_hansAllJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func wski18nResourcesZh_hansAllJsonBytes() ([]byte, error) {
	return bindataRead(
		_wski18nResourcesZh_hansAllJson,
		"wski18n/resources/zh_Hans.all.json",
	)
}

func wski18nResourcesZh_hansAllJson() (*asset, error) {
	bytes, err := wski18nResourcesZh_hansAllJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "wski18n/resources/zh_Hans.all.json", size: 0, mode: os.FileMode(420), modTime: time.Unix(1520374115, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _wski18nResourcesZh_hantAllJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func wski18nResourcesZh_hantAllJsonBytes() ([]byte, error) {
	return bindataRead(
		_wski18nResourcesZh_hantAllJson,
		"wski18n/resources/zh_Hant.all.json",
	)
}

func wski18nResourcesZh_hantAllJson() (*asset, error) {
	bytes, err := wski18nResourcesZh_hantAllJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "wski18n/resources/zh_Hant.all.json", size: 0, mode: os.FileMode(420), modTime: time.Unix(1520374115, 0)}
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
	"wski18n/resources/de_DE.all.json":   wski18nResourcesDe_deAllJson,
	"wski18n/resources/en_US.all.json":   wski18nResourcesEn_usAllJson,
	"wski18n/resources/es_ES.all.json":   wski18nResourcesEs_esAllJson,
	"wski18n/resources/fr_FR.all.json":   wski18nResourcesFr_frAllJson,
	"wski18n/resources/it_IT.all.json":   wski18nResourcesIt_itAllJson,
	"wski18n/resources/ja_JA.all.json":   wski18nResourcesJa_jaAllJson,
	"wski18n/resources/ko_KR.all.json":   wski18nResourcesKo_krAllJson,
	"wski18n/resources/pt_BR.all.json":   wski18nResourcesPt_brAllJson,
	"wski18n/resources/zh_Hans.all.json": wski18nResourcesZh_hansAllJson,
	"wski18n/resources/zh_Hant.all.json": wski18nResourcesZh_hantAllJson,
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
	"wski18n": &bintree{nil, map[string]*bintree{
		"resources": &bintree{nil, map[string]*bintree{
			"de_DE.all.json":   &bintree{wski18nResourcesDe_deAllJson, map[string]*bintree{}},
			"en_US.all.json":   &bintree{wski18nResourcesEn_usAllJson, map[string]*bintree{}},
			"es_ES.all.json":   &bintree{wski18nResourcesEs_esAllJson, map[string]*bintree{}},
			"fr_FR.all.json":   &bintree{wski18nResourcesFr_frAllJson, map[string]*bintree{}},
			"it_IT.all.json":   &bintree{wski18nResourcesIt_itAllJson, map[string]*bintree{}},
			"ja_JA.all.json":   &bintree{wski18nResourcesJa_jaAllJson, map[string]*bintree{}},
			"ko_KR.all.json":   &bintree{wski18nResourcesKo_krAllJson, map[string]*bintree{}},
			"pt_BR.all.json":   &bintree{wski18nResourcesPt_brAllJson, map[string]*bintree{}},
			"zh_Hans.all.json": &bintree{wski18nResourcesZh_hansAllJson, map[string]*bintree{}},
			"zh_Hant.all.json": &bintree{wski18nResourcesZh_hantAllJson, map[string]*bintree{}},
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
