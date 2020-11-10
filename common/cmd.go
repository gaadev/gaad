package common

import (
	"bufio"
	"fmt"
	"gaad/models"
	//    "golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"os/exec"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

//封装一个函数来执行命令
func ExecCommand(commandName string, params []string) error {

	//执行命令
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()
	errReader, err := cmd.StderrPipe()

	if err != nil {
		fmt.Println("err:" + err.Error())
		return err
	}

	//开启错误处理
	go handlerErr(errReader)

	if err != nil {
		fmt.Println(err)
		return err
	}

	cmd.Start()
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		cmdRe := ConvertByte2String(in.Bytes(), "UTF8")
		fmt.Println("->", cmdRe)
	}

	return cmd.Wait()
}

//封装一个函数来执行命令
func DeployCommand(service *models.Service, logFilePath string, commandName string, params []string) error {

	//执行命令
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()
	errReader, err := cmd.StderrPipe()

	if err != nil {
		fmt.Println("err:" + err.Error())
		return err
	}

	//开启错误处理
	go handlerErr(errReader)

	if err != nil {
		fmt.Println(err)
		return err
	}

	f, err := os.Create(logFilePath)
	defer f.Close()

	cmd.Start()
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		cmdRe := ConvertByte2String(in.Bytes(), "UTF8")
		fmt.Println("->", cmdRe)
		f.WriteString(cmdRe + "\n")
	}
	err = cmd.Wait()

	if err != nil {
		fmt.Println(err)
	}
	return err
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//开启一个协程来输出错误
func handlerErr(errReader io.ReadCloser) {
	in := bufio.NewScanner(errReader)
	for in.Scan() {
		cmdRe := ConvertByte2String(in.Bytes(), "UTF8")
		fmt.Errorf(cmdRe)
	}
}

//对字符进行转码
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	//    case GB18030:
	//        var decodeBytes,_=simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
	//        str= string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
