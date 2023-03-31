package main

import (
	"fmt"
	"log"
	"os"
)

func saveFile(fileName string, content string) {
	f, err := os.Create(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(content)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("Created : " + fileName)
}
