// author: wongoo
// since: 2019-10-19

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	projectDir := flag.String("d", "", "project directory")
	flag.Parse()

	if *projectDir == "" {
		log.Fatal("project directory required")
	}

	err := setLicenseHeader(*projectDir)
	if err != nil {
		log.Fatal(err)
	}
}

func setLicenseHeader(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		path := filepath.Join(dir, f.Name())

		if f.IsDir() {
			if err := setLicenseHeader(path); err != nil {
				return err
			}
		}

		if strings.HasSuffix(f.Name(), ".go") || strings.HasSuffix(f.Name(), ".java") {
			if err := setFileLicenseHeader(path, f.Mode(), packageStart, apacheLicenseV2); err != nil {
				return err
			}
		}

		if f.Name() == "pom.xml" {
			if err := setFileLicenseHeader(path, f.Mode(), pomXMLStart, apacheLicenseV2XML); err != nil {
				return err
			}
		}
	}

	return nil
}

func setFileLicenseHeader(filePath string, mode os.FileMode, codeStart, licenseHeader []byte) error {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	idx := bytes.Index(b, codeStart)
	if idx < 0 {
		log.Printf("code start not found in file: %s", filePath)
		return nil
	}

	log.Printf("set license header: %s", filePath)

	return ioutil.WriteFile(filePath, append(licenseHeader, b[idx:]...), mode)
}
