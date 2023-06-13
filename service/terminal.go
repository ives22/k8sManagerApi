package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"time"
)

var Terminal terminal

type terminal struct{}

// WebsocketHandler 定义websocket的handler方法
func (t *terminal) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// 解析form参数，获取namespace、podName、containerName、cluster参数
	if err := r.ParseForm(); err != nil {
		fmt.Printf("解析form参数失败, %s\n", err.Error())
		return
	}
	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	cluster := r.Form.Get("cluster")
	fmt.Printf("exec pod: %s, container: %s, namespace: %s, cluster: %s\n", podName, containerName, namespace, cluster)
	if namespace == "" || podName == "" || containerName == "" || cluster == "" {
		fmt.Printf("namespace、pod_name、container_name、cluster参数为空\n")
		return
	}
	client, err := K8s.GetClient(cluster)
	// 加载k8s配置
	conf, err := clientcmd.BuildConfigFromFlags("", K8s.GetClusterConf(cluster))
	if err != nil {
		fmt.Printf("加载k8s配置失败, %s\n", err.Error())
		return
	}
	//new一个TerminalSession类型的pty实例
	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		fmt.Printf("实例化TerminalSession失败, %s\n", err.Error())
		return
	}
	// 处理关闭，
	defer func() {
		fmt.Println("close session.")
		pty.Close()
	}()

	/*
		初始化pod所在的corev1资源组
		PodExecOptions struct 包括Container stdout stdout Command 等结构
		scheme.ParameterCodec 应该是pod 的GVK （GroupVersion & Kind）之类的
		URL长相: https://192.168.1.11:6443/api/v1/namespaces/default/pods/nginx-wf2-778d88d7c7rmsk/exec?command=%2Fbin%2Fbash&container=nginxwf2&stderr=true&stdin=true&stdout=true&tty=true
	*/
	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command:   []string{"/bin/bash"}, // TODO: 这里可以优化，从前端选择命令进入
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	//fmt.Println(req.URL())
	//fmt.Printf("exec post request url: %v\n", req)

	// remotecommand 主要实现了http 转 SPDY 添加X-Stream-Protocol-Version相关header 并发送请求
	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		fmt.Printf("建立SPDY连接失败: %v\n", err.Error())
		return
	}

	// 建立链接之后从请求的sream中发送、读取数据
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	})
	if err != nil {
		msg := fmt.Sprintf("Exec to pod error! err: %v", err)
		fmt.Println(msg)
		// 将报错返回出去
		pty.Write([]byte(msg))
		// 标记退出stream流
		pty.Down()
	}

}

// TerminalMessage 消息内容
/*
TerminalMessage定义了终端和容器shell交互内容的格式
Operation是操作类型
Data是具体数据内容
Rows和Cols可以理解为终端的行数和列数，也就是宽、高
*/
type TerminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

// TerminalSession 交互的结构体，接管输入和输出
/*
//定义TerminalSession结构体，实现PtyHandler接口
//wsConn是websocket连接
//sizeChan用来定义终端输入和输出的宽和高
//doneChan用于标记退出终端
*/
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

// 初始化一个websocket.Upgrader类型的对象，用于http协议升级为websocket协议
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader
}()

// NewTerminalSession 该方法用于升级http协议至websocket，并new一个TerminalSession类型的对象返回
func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeadler http.Header) (*TerminalSession, error) {
	// 升级ws协议
	conn, err := upgrader.Upgrade(w, r, responseHeadler)
	if err != nil {
		return nil, errors.New("升级websocket失败, " + err.Error())
	}
	// new
	terminalSession := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	return terminalSession, nil
}

// Read 读数据的方法, 用于读取web端的输入，接收web端输入的指令内容
func (t *TerminalSession) Read(p []byte) (int, error) {
	// 从websocket中读取消息
	_, message, err := t.wsConn.ReadMessage()

	if err != nil {
		fmt.Printf("read message err: %v\n", err)
		return 0, err
	}

	// 反序列化
	var msg TerminalMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		fmt.Printf("read parse message err: %v\n", err)
		return 0, err
	}

	// 逻辑判断
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		fmt.Printf("unknown message type '%s'\n", msg.Operation)
		return 0, fmt.Errorf("unknown message type '%s'\n", msg.Operation)
	}
}

// 写数据的方法，用于向web端输出，接收web端的指令后，将结果返回出去
func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		fmt.Printf("write parse message err: '%v'\n", err)
		return 0, err
	}
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		fmt.Printf("write message err: '%v'\n", err)
		return 0, err
	}
	return len(p), nil

}

// Down 标记关闭的方法，关闭后触发退出终端
func (t *TerminalSession) Down() {
	close(t.doneChan)
}

// Close 用于关闭websocket的连接
func (t *TerminalSession) Close() error {
	return t.wsConn.Close()

}

// Next 获取web端是否resize，以及是否退出终端
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}