package main

import (
	"encoding/json"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

const configFilename = ".icebridge"

type ConfigFile struct {
	ClientId      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	Token         string `json:"token"`
	LocalPath     string `json:"local_path"`
	DropboxPath   string `json:"dropbox_path"`
	Cursor        string `json:"cursor"`
	MaxFileAge    int    `json:"max_file_age"`
	FileAgeMethod string `json:"file_age_method"`
	changed       bool   `json:"-"`
}

func (conf *ConfigFile) Read(fname string) error {
	var buf []byte
	var file string
	var err error

	if file, err = openFile(fname); err == nil {
		if buf, err = ioutil.ReadFile(file); err == nil {
			err = json.Unmarshal(buf, conf)
		}
	}

	return err
}

func (conf *ConfigFile) Write(fname string) error {
	var buf []byte
	var file string
	var err error

	if file, err = openFile(fname); err == nil {
		if buf, err = json.MarshalIndent(conf, "", ""); err == nil {
			err = ioutil.WriteFile(file, buf, 0600)
		}
	}

	return err
}

func openFile(fname string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		// this should return nil perhaps
		return "", err
	}
	return filepath.Join(usr.HomeDir, fname), nil
}
