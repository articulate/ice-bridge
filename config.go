package main

import (
	"os"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
)

type ConfigFile struct {
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Token string `json:"token"`
	LocalPath string `json:"local_path"`
	DropboxPath string `json:"dropbox_path"`
	changed bool `json:"-"`
}

func (conf *ConfigFile) Read(fname string) error {
	var err error
	var buf []byte

	file := openFile(fname)
	if buf, err = ioutil.ReadFile(file); err == nil {
		err = json.Unmarshal(buf, conf)
	}

	return err
}

func (conf *ConfigFile) Write(fname string) error {
	var err error
	var buf[]byte

	file := openFile(fname)
	if buf, err = json.MarshalIndent(conf, "", ""); err == nil {
		err = ioutil.WriteFile(file, buf, 0600)
	}

	return err
}

func openFile(fname string) string {
	return filepath.Join(os.Getenv("HOME"), fname)
}
