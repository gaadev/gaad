package base

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gaad/models"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	upGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024 * 1024 * 10,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

//常数
const (
	wsMsgCmd    = "cmd"
	wsMsgResize = "resize"
)

//ShellClient
type ShellClient struct {
	StdinPipe       io.WriteCloser
	ComboOutput     *safeBuffer //ssh 终端混合输出
	LogBuff         *safeBuffer //保存session的日志
	InputFilterBuff *safeBuffer //用来过滤输入的命令和ssh_filter配置对比的
	Session         *ssh.Session
	WsConn          *websocket.Conn
	WaitGroup       waitGroupWrapper
	Exit            chan bool //退出标识
}

//安全缓存结构体
type safeBuffer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

type waitGroupWrapper struct {
	sync.WaitGroup
}

//处理webSocket和创建shell
func HandleWsAndShell(node *models.Node, cols, rows int, c *gin.Context) {
	//升级为websocket
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		controllers.Response(models.ParameterIllegal, "协议升级失败", nil)
		return
	}
	defer ws.Close()
	//创建shell
	shell, err := createShell(node)
	if err != nil {
		controllers.Response(models.ParameterIllegal, "创建shell失败", nil)
		return
	}
	defer shell.Close()
	shellClient, err := CreatShellClient(cols, rows, shell, ws)
	if err != nil {
		return
	}
	defer shellClient.Close()
	//开启读取,阻塞等待
	shellClient.WaitGroup.Wrap(func() { shellClient.read() })
	//开启写出
	shellClient.WaitGroup.Wrap(func() { shellClient.write() })
	shellClient.WaitGroup.Wrap(func() { shellClient.Wait() })
	shellClient.Wait()

}

//创建shell连接
func createShell(node *models.Node) (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:22", node.Ip) //拼接ip和端口
	shell, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            node.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(node.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}
	return shell, nil
}

//将ws和shell绑定,创建xterm模板
func CreatShellClient(cols, rows int, client *ssh.Client, wsConn *websocket.Conn) (*ShellClient, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	//获取读取的结果
	stdinP, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	comboWriter := &safeBuffer{}
	logBuf := &safeBuffer{}   //日志缓存
	inputBuf := &safeBuffer{} //输入缓存

	session.Stdout = comboWriter
	session.Stderr = comboWriter
	//创建Xterm模板
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // 是否回显
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	//使用xterm,rows行数cols格式。
	if err := session.RequestPty("xterm", rows, cols, modes); err != nil {
		return nil, err
	}
	//使用远程shell
	err = session.Shell()
	if err != nil {
		return nil, err
	}
	return &ShellClient{
		StdinPipe:       stdinP,
		ComboOutput:     comboWriter,
		LogBuff:         logBuf,
		InputFilterBuff: inputBuf,
		Session:         session,
		WsConn:          wsConn,
		Exit:            make(chan bool),
	}, nil
}

//读取消息
func (sc *ShellClient) read() {
	for {
		select {
		//表示退出
		case <-sc.Exit:
			return
		default:
			_, byteS, err := sc.WsConn.ReadMessage() //读取消息
			if err != nil {
				log.Error("webSocket读取消息失败", err)
				return
			}
			wsMsg := models.WsMsg{}
			if err := json.Unmarshal(byteS, &wsMsg); err != nil {
				log.Error("wsMsg JSON转换异常")
			}
			//处理消息
			switch wsMsg.Type {
			case wsMsgResize:
				//改变行数和行数
				if wsMsg.Cols > 0 && wsMsg.Rows > 0 {
					if err := sc.Session.WindowChange(wsMsg.Rows, wsMsg.Cols); err != nil {
						log.Error("shell 窗口大小改变失败")
					}
				}
			case wsMsgCmd:
				//处理执行命令
				byteS, err := base64.StdEncoding.DecodeString(wsMsg.Cmd)
				if err != nil {
					log.Error("webSocket执行命令转换失败", err)
				}
				if _, err := sc.StdinPipe.Write(byteS); err != nil {
					log.Error("StdinPipe处理命令失败", err)
				}
			}
		}
	}
}

//写出消息
func (sc *ShellClient) write() {
	tick := time.NewTicker(time.Millisecond * time.Duration(60))
	defer tick.Stop()
	for {
		select {
		case <-sc.Exit:
			return
		case <-tick.C:
			if sc.ComboOutput == nil {
				return
			}
			byteS := sc.ComboOutput.Bytes()
			if len(byteS) > 0 {
				err := sc.WsConn.WriteMessage(websocket.TextMessage, byteS) //写出消息
				if err != nil {
					log.Error("消息写出失败", err)
				}
				//日志缓存
				_, err = sc.LogBuff.Write(byteS)
				if err != nil {
					log.Error("写入日志缓存失败", err)
				}
				//缓存清除
				sc.ComboOutput.buffer.Reset()
			}
		}
	}

}

//Close 关闭
func (sc *ShellClient) Close() {
	if sc.Session != nil {
		sc.Session.Close()
	}
	if sc.LogBuff != nil {
		sc.LogBuff = nil
	}
	if sc.ComboOutput != nil {
		sc.ComboOutput = nil
	}
}

//等待
func (sc *ShellClient) Wait() {
	if err := sc.Session.Wait(); err != nil {
		log.Info(err.Error())
		// sws.exit <- true
		close(sc.Exit)
	}
	log.Debug("remote command to exit.")
	close(sc.Exit)
}

func (w *safeBuffer) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

//获取字节
func (w *safeBuffer) Bytes() []byte {
	w.mu.Lock()
	defer w.mu.Unlock() //该方法调用结束后，解锁
	return w.buffer.Bytes()
}

//重置
func (w *safeBuffer) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock() //该方法调用结束后，解锁
	w.buffer.Reset()
}

func (w *waitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
