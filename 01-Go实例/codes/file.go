/**
* @Author: XJX
* @Description: 文件操作演示
* @File: file.go
* @Date: 2020/6/23 16:13
 */

package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

var testPath string = "D:/wamp64/www/fileTest" //测试路径

//使用 bufio.NewWriter 写入文件
func createFileByNewWriter(fileName, data string) error {
	fileName = testPath + "/" + fileName
	if fileIsExist(fileName) {
		return errors.New(fileName + " is Exist!!!")
	} else {
		f, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer f.Close()
		writer := bufio.NewWriter(f)
		_, err1 := writer.WriteString(data)
		if err1 == nil {
			writer.Flush()
		} else {
			err = err1
		}
		return err
	}
}

// 使用 io.WriteString 写入文件
func createFileByWriteString(fileName, data string) error {
	fileName = testPath + "/" + fileName
	if fileIsExist(fileName) {
		return errors.New(fileName + " is Exist!!!")
	} else {
		f, serr := os.Create(fileName)
		if serr != nil {
			return serr
		}
		defer f.Close()
		_, err := io.WriteString(f, data)
		return err
	}
}

//ioutil.WriteFile 写入文件
func createFileByWriteFile(fileName, data string) error {
	fileName = testPath + "/" + fileName
	if fileIsExist(fileName) {
		return errors.New(fileName + " is Exist!!!")
	} else {
		var sdata = []byte(data)
		err := ioutil.WriteFile(fileName, sdata, 0666)
		return err
	}
}

//f.Write和 f.WriteString 写入文件
func createFileByFileWrite(fileName, data string) error {
	fileName = testPath + "/" + fileName
	if fileIsExist(fileName) {
		return errors.New(fileName + " is Exist!!!")
	} else {
		f, serr := os.Create(fileName)
		if serr != nil {
			return serr
		}
		defer f.Close()
		var sdata = []byte(data)
		_, serr1 := f.Write(sdata)
		if serr1 != nil {
			return serr1
		}
		_, err := f.WriteString(data)
		f.Sync()
		return err
	}
}

func createJsonDataFile(fileName string, data []JsonData) error {
	fileName = testPath + "/" + fileName
	if fileIsExist(fileName) {
		return errors.New(fileName + " is Exist!!!")
	} else {
		f, cerr := os.Create(fileName)
		if cerr != nil {
			return cerr
		}
		defer f.Close()
		encoder := json.NewEncoder(f)
		err := encoder.Encode(data)
		return err
	}
}

func createXmlDataFile(fileName string, data []XmlData) error {
	fileName = testPath + "/" + fileName
	if fileIsExist(fileName) {
		return errors.New(fileName + " is Exist!!!")
	} else {
		f, cerr := os.Create(fileName)
		if cerr != nil {
			return cerr
		}
		defer f.Close()
		encoder := xml.NewEncoder(f)
		err := encoder.Encode(data)
		return err
	}
}

func createGobDataFile(fileName string, data []GobData) error {
	fileName = testPath + "/" + fileName
	if fileIsExist(fileName) {
		return errors.New(fileName + " is Exist!!!")
	} else {
		f, cerr := os.Create(fileName)
		if cerr != nil {
			return cerr
		}
		defer f.Close()
		encoder := gob.NewEncoder(f)
		err := encoder.Encode(data)
		return err
	}
}

//读取文件
func readByIoReadFile(fileName string) (data interface{}, err error) {
	fileName = testPath + "/" + fileName
	if !fileIsExist(fileName) {
		return nil, errors.New(fileName + " is Not Exist!!!")
	} else {
		data, err := ioutil.ReadFile(fileName)
		return string(data), err
	}
}

func readByIoReadAll(fileName string) (data interface{}, err error) {
	fileName = testPath + "/" + fileName
	if !fileIsExist(fileName) {
		return nil, errors.New(fileName + " is Not Exist!!!")
	} else {
		f, err1 := os.Open(fileName)
		if err1 != nil {
			return nil, err1
		}
		data, err := ioutil.ReadAll(f)
		return string(data), err
	}
}

func readByNewReader(fileName string) (data interface{}, err error) {
	fileName = testPath + "/" + fileName
	if !fileIsExist(fileName) {
		return nil, errors.New(fileName + " is Not Exist!!!")
	} else {
		f, err1 := os.Open(fileName)
		if err1 != nil {
			return nil, err1
		}
		r := bufio.NewReader(f)
		defer f.Close()
		var chunk []byte
		buf := make([]byte, 1024)
		for {
			n, err := r.Read(buf)
			fmt.Println(n)
			fmt.Println(string(buf))
			if err != nil && err != io.EOF {
				return nil, err
			}
			if n == 0 {
				break
			}
			chunk = append(chunk, buf[:n]...)
		}
		return string(chunk), err
	}
}

func readByFileRead(fileName string) (data interface{}, err error) {
	fileName = testPath + "/" + fileName
	if !fileIsExist(fileName) {
		return nil, errors.New(fileName + " is Not Exist!!!")
	} else {
		f, err1 := os.Open(fileName)
		if err1 != nil {
			return nil, err1
		}
		defer f.Close()
		var chunk []byte
		buf := make([]byte, 1024)
		for {
			n, err := f.Read(buf)
			fmt.Println(n)
			fmt.Println(string(buf))
			if err != nil && err != io.EOF {
				return nil, err
			}
			if n == 0 {
				break
			}
			chunk = append(chunk, buf[:n]...)
		}
		return string(chunk), err
	}
}

func writeFile(fileName, data string) error {
	fileName = testPath + "/" + fileName
	f, err1 := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0666)
	if err1 != nil {
		return err1
	}
	_, err := io.WriteString(f, data)
	return err
}

func readJsonDataFile(fileName string) (data interface{}, err error) {
	fileName = testPath + "/" + fileName
	if !fileIsExist(fileName) {
		return nil, errors.New(fileName + " is Not Exist!!!")
	} else {
		f, err1 := os.Open(fileName)
		if err1 != nil {
			return nil, err1
		}
		defer f.Close()
		var wdata []JsonData
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&wdata)
		return wdata, err
	}
}

//删除文件
func deleteFile(fileName string) error {
	fileName = testPath + "/" + fileName
	if !fileIsExist(fileName) {
		return errors.New(fileName + " is Not Exist!!!")
	}
	err := os.Remove(fileName)
	return err
}

func deleteDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

//遍历目录
func listFile(dir string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			s, err := listFile(dir+"/"+fi.Name(), s)
			if err != nil {
				return s, err
			}
		} else {
			s = append(s, dir+"/"+fi.Name())
		}
	}
	return s, nil
}

//文件是否存在
func fileIsExist(fileName string) bool {
	flag := true
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		flag = false
	}
	return flag
}

type JsonData struct {
	Name   string
	Url    string
	Course []string
}

type XmlData struct {
	Name   string
	Url    string
	Course []string
}

type GobData struct {
	Name   string
	Url    string
	Course []string
}

type FileLock struct {
	FileName string
	f        *os.File
}

func NewFileLock(FileName string) *FileLock {
	return &FileLock{
		FileName: FileName,
	}
}

func (l *FileLock) Lock() error {
	f, err := os.Open(l.FileName)
	if err != nil {
		return err
	}
	l.f = f
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return fmt.Errorf("cannot flock directory %s - %s", l.FileName, err)
	}
	return nil
}

func (l *FileLock) Unlock() error {
	defer l.f.Close()
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}

func main() {
	//l := NewFileLock("D:/wamp64/www/fileTest/FileWrite.txt")

	//jdata := []JsonData{{"php", "www.qq.com", []string{"22222", "333333", "4444444"}}, {"golang", "www.golang.com", []string{"342", "23", "242"}}}
	//xdata := []XmlData{{"php", "www.qq.com", []string{"22222", "333333", "4444444"}}, {"golang", "www.golang.com", []string{"342", "23", "242"}}}
	//gdata := []GobData{{"php", "www.qq.com", []string{"22222", "333333", "4444444"}}, {"golang", "www.golang.com", []string{"342", "23", "242"}}}
	/*err := createFileByNewWriter("NewWriter.txt", "NewWriter")
	if err != nil {
		fmt.Println("create file fail, error:" + err.Error())
	}
	err := createFileByWriteFile("WriteFile.txt", "WriteFile")
	if err != nil {
		fmt.Println("create file fail, error:" + err.Error())
	}

	err := createFileByWriteString("WriteString.txt", "WriteString")
	if err != nil {
		fmt.Println("create file fail, error:" + err.Error())
	}

	err := createFileByFileWrite("FileWrite.txt", "FileWrite")
	if err != nil {
		fmt.Println("create file fail, error:" + err.Error())
	}*/

	//createJsonDataFile("data.json", jdata)
	//createXmlDataFile("data.xml", xdata)
	//createGobDataFile("data.gob", gdata)
	//data, _ := readByIoReadFile("NewWriter.txt")
	//fmt.Println(data)
	//fmt.Println(err.Error())

	//data, err := readJsonDataFile("data.json")
	//fmt.Println(data)
	//fmt.Println(err.Error())

	//writeFile("WriteString.txt", "\rhello, xjx!!!")

	/*err := deleteFile("data.xml")
	if err == nil {
		fmt.Println("delete file succ!!!")
	}*/

	/*folderName := time.Now().Format("2006-01-02")
	folderPath := filepath.Join(testPath, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, 0777)
		os.Chmod(folderPath, 0777)
	}
	testPath = folderPath
	for i := 0; i < 10; i++ {
		createFileByWriteString("testfile_"+strconv.Itoa(i)+".txt", "hello "+strconv.Itoa(i))
	}*/

	//deleteDir(folderPath);
	/*var s []string
	s, err := listFile("D:/wamp64/www/static", s)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(s)
	}*/

	test_file_path, _ := os.Getwd()
	locked_file := test_file_path
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(num int) {
			flock := NewFileLock(locked_file)
			err := flock.Lock()
			if err != nil {
				wg.Done()
				fmt.Println(err.Error())
				return
			}
			fmt.Printf("output : %d\n", num)
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(2 * time.Second)

}
