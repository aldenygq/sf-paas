package middleware

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//目录是否存在
func IsDir(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

/**
创建文件夹
*/
func CreateDir(dirName string) bool {
	err := os.Mkdir(dirName, 755)
	if err != nil {
		return false
	}
	return true
}

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}

func ToString(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

//remove file name trailing spaces
func ToTrimString(filePath string) (string, error) {
	str, err := ToString(filePath)
	if err != nil {
		return "", err
	}
	return str, nil
}

func ReadFileInfo(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(fd), nil
}

func ReadFileToStruct(file string, data interface{}) error {
	filecontent, err := ReadFileInfo(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(filecontent), &data)
	if err != nil {
		return err
	}

	return nil
}
