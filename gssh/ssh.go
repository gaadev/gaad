package gssh

import (
	"bufio"
	"errors"
	"fmt"
	"gaad/common"
	"gaad/gssh/comp"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//src 可以是文件或者文件夾
func ScpCompress(client *Client, src string, dist string) {

	Client, err := ssh.Dial(client.Network, client.HostPort, client.ClientConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer Client.Close()
	if session, err := Client.NewSession(); err == nil {
		uid := uuid.NewV4()
		compFileName := uid.String() + ".tar.gz"
		//沒有/tmp文件夾自动生成
		common.CreateFile("./tmp")
		compFile := "./tmp/" + compFileName
		defer session.Close()
		go func() {
			Buf := make([]byte, 1024)
			w, _ := session.StdinPipe()
			var f *os.File
			defer w.Close()
			defer func() {
				f.Close()
				err := os.Remove(compFile)

				if err != nil {
					fmt.Println(err)
				}
			}()

			comp.Compress(src, compFile)

			f, _ = os.Open(compFile)
			info, _ := f.Stat()
			fmt.Fprintln(w, "C0644", info.Size(), compFileName)
			for {
				n, err := f.Read(Buf)
				fmt.Fprint(w, string(Buf[:n]))
				if err != nil {
					if err == io.EOF {
						return
					} else {
						panic(err)
					}
				}
			}
		}()

		stdout, err := session.StdoutPipe()

		cmd := "/usr/bin/scp -qrt " + dist + " >/dev/null 2>&1;cd " + dist + "; tar -zxvf " + compFileName + ";rm -rf " + compFileName + ";"

		err = session.Start(cmd)
		if err != nil {
			log.Fatal(err)
		}
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(string(in.Bytes()))
		}
		session.Wait()

	}

}

//src 可以是文件或者文件夾
func Scp(client *Client, src string, dist string, isDelete bool) error {

	Client, err := ssh.Dial(client.Network, client.HostPort, client.ClientConfig)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer Client.Close()
	if session, err := Client.NewSession(); err == nil {
		defer session.Close()

		if strings.HasSuffix(src, "/") {
			err = errors.New("scp 不用使用文件夹")
			return err
		}
		pos := strings.LastIndex(src, "/")
		filename := src[pos+1:]

		go func() {

			Buf := make([]byte, 1024)
			w, _ := session.StdinPipe()
			var f *os.File
			defer w.Close()
			defer func() {
				f.Close()
				if isDelete {
					err = os.Remove(src)
				}

				if err != nil {
					fmt.Println(err)
				}
			}()

			f, _ = os.Open(src)
			info, _ := f.Stat()

			fmt.Fprintln(w, "C0644", info.Size(), filename)
			for {
				n, err := f.Read(Buf)
				fmt.Fprint(w, string(Buf[:n]))
				if err != nil {
					if err == io.EOF {
						return
					} else {
						panic(err)
					}
				}
			}
		}()

		stdout, err := session.StdoutPipe()

		cmd := "/usr/bin/scp -qrt " + dist + " >/dev/null 2>&1;cd " + dist + "; tar -zxvf " + filename + ";rm -rf " + filename + ";"

		err = session.Start(cmd)
		if err != nil {
			log.Fatal(err)
		}
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(string(in.Bytes()))
		}
		session.Wait()

	}
	return err

}
func RunShell(client *Client, shell string) {

	Client, err := ssh.Dial(client.Network, client.HostPort, client.ClientConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer Client.Close()
	if session, err := Client.NewSession(); err == nil {
		defer session.Close()

		stdout, err := session.StdoutPipe()

		err = session.Start(shell)
		if err != nil {
			log.Fatal(err)
		}
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(string(in.Bytes()))
		}
		session.Wait()

	}

}

//src 可以是文件或者文件夾
func RunShellFile(client *Client, shellFile string) {
	f, err := ioutil.ReadFile(shellFile)
	if err == nil {
		fmt.Println(string(f))
	}
	RunShell(client, string(f))
}
