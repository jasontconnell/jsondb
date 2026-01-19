package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func ensureJsonFile(f string) (*os.File, error) {
	file, err := os.Open(f)
	if err != nil {
		log.Println("ensuring file")
		file, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR, os.ModePerm)
		return file, err
	}
	return file, nil
}

func readJsonFileList[T any](f string) ([]T, error) {
	file, err := ensureJsonFile(f)
	if err != nil {
		return nil, fmt.Errorf("can't read file %s. %w", f, err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	dest := []T{}
	err = dec.Decode(&dest)
	if err != nil {
		log.Printf("problem parsing json %s. returning empty list.\n", f)
	}
	return dest, nil
}

func writeJsonFile(f string, data interface{}) error {
	file, err := os.OpenFile(f, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can't open json file for write %s. %w", f, err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent(" ", " ")
	err = enc.Encode(data)
	if err != nil {
		return fmt.Errorf("writing json to file %s. %w", f, err)
	}
	return nil
}
