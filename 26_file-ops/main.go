package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	directoryName        = "my_temp_files"
	originalFilename     = "original.txt"
	copiedFilename       = "copied.txt"
	fullOriginalFilePath = directoryName + "/" + originalFilename
	fullCopiedFilePath   = directoryName + "/" + copiedFilename
)

func main() {
	err := os.MkdirAll(directoryName, 0666)
	if err != nil {
		log.Fatalf("Error while creating directory %s: %v", directoryName, err)
	}
	defer func() {
		fmt.Println("Removing directory ", directoryName)
		err := os.RemoveAll(directoryName)
		if err != nil {
			log.Fatalf("Error while removing directory %s: %v", directoryName, err)
		}
	}()

	createAndWriteOriginalFile()

	readEntireOriginalFile()

	readOriginalFileLineByLine()

	getFileInformation()

	copyFile()

	verifyCopiedFileContent()
}

func verifyCopiedFileContent() {
	file, err := os.Open(fullCopiedFilePath)
	if err != nil {
		log.Fatalf("Error while opening %s: %v", fullCopiedFilePath, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error while reading %s: %v", fullCopiedFilePath, err)
	}
	fmt.Println("--- Copied File Content ---")
	fmt.Println(string(content))
}

func copyFile() {
	file, err := os.Open(fullOriginalFilePath)
	if err != nil {
		log.Fatalf("Error while trying to open %s: %v", fullOriginalFilePath, err)
	}
	defer file.Close()

	copyFile, err := os.Create(fullCopiedFilePath)
	if err != nil {
		log.Fatalf("Error while trying to create %s: %v", fullCopiedFilePath, err)
	}
	defer copyFile.Close()

	n, err := io.Copy(copyFile, file)
	if err != nil {
		log.Fatalf("Error while copying %s to %s: %v", fullOriginalFilePath, fullCopiedFilePath, err)
	}
	fmt.Printf("%d bytes of data successfully copied from %s to %s\n", n, fullOriginalFilePath, fullCopiedFilePath)
}

func getFileInformation() {
	fileInfo, err := os.Stat(fullOriginalFilePath)
	if err != nil {
		log.Fatalf("Error while getting %s file info: %v", fullOriginalFilePath, err)
	}

	fmt.Println("--- File Info for original.txt ---")
	fmt.Printf("File name: %s\n", fileInfo.Name())
	fmt.Printf("Size (bytes): %d\n", fileInfo.Size())
	fmt.Printf("Modification time: %s\n", fileInfo.ModTime().String())

}

func readOriginalFileLineByLine() {
	file, err := os.Open(fullOriginalFilePath)
	if err != nil {
		log.Fatalf("Error while open %s file: %v", fullOriginalFilePath, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	fmt.Println("--- Line-by-Line content of original.txt ---")
	for i := 1; scanner.Scan(); i++ {
		line := scanner.Text()
		fmt.Printf("%d %s\n", i, line)
	}

	err = scanner.Err()
	if err != nil {
		fmt.Printf("Scanner error occurred %v\n", err)
	}

}

func readEntireOriginalFile() {
	oriFile, err := os.Open(fullOriginalFilePath)
	if err != nil {
		log.Fatalf("Error while open file %s: %v", fullOriginalFilePath, err)
	}
	defer oriFile.Close()
	oriContent, err := io.ReadAll(oriFile)
	if err != nil {
		log.Fatalf("Error while reading %s: %v", fullOriginalFilePath, err)
	}
	fmt.Println("--- Full content of original.txt ---")
	fmt.Println(string(oriContent))
}

func createAndWriteOriginalFile() {
	originalFile, err := os.Create(fullOriginalFilePath)
	if err != nil {
		log.Fatalf("Error while creating file %s: %v", fullOriginalFilePath, err)
	}
	defer originalFile.Close()

	// write into file
	_, err = originalFile.WriteString("Hello from Go!\nThis is the second line.\nEng of original content.")
	if err != nil {
		log.Fatalf("Error while write to file %s: %v", fullOriginalFilePath, err)
	}
	fmt.Printf("Successful writing to file %s.\n", fullOriginalFilePath)
}
