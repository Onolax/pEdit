package cmd

import (
	"fmt"
	"io"
	"os"
)

type Blob struct {
	filename string
	content  string
	file     *os.File
}

func InitBlob(filename string) *Blob {
	var blob *Blob
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err)
		return blob
	}
	fileInfo, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		file.Close()
		return blob
	}
	return &Blob{filename: filename, content: string(fileInfo), file: file}
}

func (blob *Blob) SetContent(content string) {
	blob.content = content
}
func (blob *Blob) GetContent() string {
	return blob.content
}

func (blob *Blob) Close() {
	blob.file.Close()
}

func (blob *Blob) Save(content string) {
	blob.content = content
	err := os.WriteFile(blob.filename, []byte(blob.content), 0777)
	if err != nil {
		fmt.Println(err)
	}
}
