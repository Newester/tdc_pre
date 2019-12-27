package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"runtime"
)

const file = "./test.gob"

type User struct {
	Name, Pass string
}

func main() {
	var datato = &User{"Donald", "DuckPass"}
	var datafrom = new(User)

	err := Save(file, datato)
	Check(err)
	err = Append(file)
	Check(err)
	err = Load(file, datafrom)
	Check(err)
	fmt.Println(datafrom)
}

// Encode via Gob to file
func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Decode Gob file
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

// Append UnusedBytes
func Append(path string) error {
	fileInfo, _ := os.Stat(path)
	used := fileInfo.Size()
	fmt.Printf("file size: %d\n", used)
	file, err := os.OpenFile(path, os.O_APPEND, os.ModeAppend)
	if err == nil {
		b := make([]byte, 1024-used)

		_, err = file.Write(b)
		if err != nil {
			file.Close()
			log.Fatalf("write file error %v", err)
		}
	}
	file.Close()
	fileInfo, _ = os.Stat(path)
	fmt.Printf("file size: %d\n", fileInfo.Size())

	return err
}

func Check(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", e)
		os.Exit(1)
	}
}
