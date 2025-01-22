package fileutils

import (
	"bufio"
	"fmt"
	"os"
)

// WriteLinesWithBufferFlag writes the lines to the given file.
func WriteLinesWithBufferFlag(lines []string, fileName string, flag int) error {
	var (
		file *os.File
		err  error
	)
	if FileExists(fileName) {
		file, err = os.OpenFile(fileName, flag, 0666)
		if err != nil {
			fmt.Println("Open file err =", err)
			return err
		}
	} else {
		file, err = os.Create(fileName) //创建文件
		if err != nil {
			fmt.Println("file create fail")
			return err
		}
	}

	defer file.Close()
	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
