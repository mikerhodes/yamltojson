package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// main accepts a file path as single program argument. It can be:
// - a dir, in which case yamltojson converts all `.yaml` files in the folder.
// - a file, in which case yamltojson assumes a YAML file and converts.
// The output for any input files is a file alongside the original with a
// `.json` suffix.
func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: yamltojson [file or directory]")
		log.Fatal("Error: missing path to convert")
	}

	pathArg := os.Args[1]

	info, err := os.Stat(pathArg)
	if err != nil {
		log.Fatalf("Could not stat %s: %s", pathArg, err.Error())
	}

	if info.IsDir() {
		// Find .yaml files in the folder, and convert them.
		dir := pathArg
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) != ".yaml" {
				continue
			}

			fullPath := filepath.Join(dir, file.Name())
			if err := convertPath(fullPath); err != nil {
				// Print an error if there is one, then continue
				log.Print(err)
			}
		}
	} else {
		// Assume path is valid YAML regardless of extension, and convert it.
		if err := convertPath(pathArg); err != nil {
			log.Fatal(err)
		}
	}
}

func convertPath(p string) error {
	dir := path.Dir(p)
	fn := path.Base(p)
	ofn := strings.TrimSuffix(fn, filepath.Ext(fn)) + ".json"

	log.Printf("%s -> %s", fn, ofn)

	if err := convertFile(
		filepath.Join(dir, fn),
		filepath.Join(dir, ofn),
	); err != nil {
		return err
	}
	return nil
}

// convertFile reads an input YAML file, i, and converts it into
// and output JSON file, o.
func convertFile(i, o string) error {
	data, err := ioutil.ReadFile(i)
	if err != nil {
		return err
	}
	var body interface{}
	if err := yaml.Unmarshal(data, &body); err != nil {
		panic(err)
	}

	body, err = convert(body)
	if err != nil {
		return err
	}

	if b, err := json.MarshalIndent(body, "", "    "); err != nil {
		return err
	} else {
		ioutil.WriteFile(o, b, 0644)
	}

	return nil
}

// convert remaps a YAML parsed into map[interface{}]interface{} into
// map[string](another map[string]map[string] etc) recursively so
// it can be written by the JSON marshaller.
func convert(i interface{}) (interface{}, error) {
	switch x := i.(type) {

	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			k2, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("Could not convert field name '%v' to string", k)
			}

			v2, err := convert(v)
			if err != nil {
				return nil, fmt.Errorf("Could not convert %v to JSON: %v", v, err)
			}

			m2[k2] = v2
		}
		return m2, nil

	case []interface{}:
		for i, v := range x {
			v2, err := convert(v)
			if err != nil {
				return nil, fmt.Errorf("Could not convert %v to JSON: %v", v, err)
			}
			x[i] = v2
		}
	}

	return i, nil
}
