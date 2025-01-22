package fileutils

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// FileExists 判断文件或是文件夹是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// SplitPath 将路径按分隔符分隔成字符串数组。比如：
//  /a/b/c  ==>  []string{"a", "b", "c"}
func SplitPath(path string) []string {
	vol := filepath.VolumeName(path)
	ret := make([]string, 0, 10)

	index := 0
	if len(vol) > 0 {
		ret = append(ret, vol)
		path = path[len(vol)+1:]
	}
	for i := 0; i < len(path); i++ {
		if os.IsPathSeparator(path[i]) {
			if i > index {
				ret = append(ret, path[index:i])
			}
			index = i + 1 // 过滤掉此符号
		}
	}

	if len(path) > index {
		ret = append(ret, path[index:])
	}

	return ret
}

func genericIsPathExists(pth string) (os.FileInfo, bool, error) {
	if pth == "" {
		return nil, false, errors.New("no path provided")
	}
	fileInf, err := os.Lstat(pth)
	if err == nil {
		return fileInf, true, nil
	}
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	return fileInf, false, err
}

// IsPathExists ...
func IsPathExists(pth string) (bool, error) {
	_, isExists, err := genericIsPathExists(pth)
	return isExists, err
}

// IsDirExists ...
func IsDirExists(pth string) (bool, error) {
	fileInf, isExists, err := genericIsPathExists(pth)
	if err != nil {
		return false, err
	}
	if !isExists {
		return false, nil
	}
	if fileInf == nil {
		return false, errors.New("no file info available")
	}
	return fileInf.IsDir(), nil
}

// IsRelativePath ...
func IsRelativePath(pth string) bool {
	if strings.HasPrefix(pth, "./") {
		return true
	}

	if strings.HasPrefix(pth, "/") {
		return false
	}

	if strings.HasPrefix(pth, "$") {
		return false
	}

	return true
}

// ListDir 按照文件后缀匹配目录文件
func ListDir(dirPath, suffix string) ([]string, error) {
	var files []string
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	suffix = strings.ToLower(suffix)
	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(strings.ToLower(file.Name()), suffix) {
			files = append(files, path.Join(dirPath, file.Name()))
		}
	}

	return files, nil
}