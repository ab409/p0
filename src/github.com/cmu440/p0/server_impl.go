// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
    "net"
    "bufio"
    "fmt"
)

type multiEchoServer struct {
    clients map[string] net.Conn
    clientMsgs map[string] chan([] byte)
    listener net.Listener
    read chan []byte
    stop bool
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
    mes := &multiEchoServer{
        clients: make(map[string] net.Conn),
        clientMsgs: make(map[string] chan([]byte)),
        read: make(chan([]byte)),
        stop: false,
    }
    return MultiEchoServer(mes);
}

func (mes *multiEchoServer) Start(port int) error {
    mes.stop = false;
    listener, e := net.Listen("tcp", fmt.Sprintf(":%d", port));
    if e != nil {
        fmt.Println("listen tcp failed");
        return e;
    }
    mes.listener = listener;
    go mes.handleListen();
    go mes.handleDistributeMsg();
    return nil;
}

func (mes *multiEchoServer) Close() {
    mes.stop = true;
    mes.listener.Close();
    for _, conn := range mes.clients{
        conn.Close();
    }
}

func (mes *multiEchoServer) Count() int {
    return len(mes.clients);
}

// TODO: add additional methods/functions below!
func (mes *multiEchoServer) handleConn(conn net.Conn) {
    go mes.handleWrite(conn);
    reader := bufio.NewReader(conn);
    for {
        line, e := reader.ReadBytes('\n');
        if  e != nil{
            fmt.Println("read close connetion, client : " + conn.RemoteAddr().String());
            delete(mes.clients, conn.RemoteAddr().String());
            conn.Close();
            return;
        }
        mes.read <- line;
    }
}

func (mes *multiEchoServer) handleWrite(conn net.Conn) {
    msgChan := mes.clientMsgs[conn.RemoteAddr().String()];
    for {
        msg := <- msgChan;
        _, e := conn.Write(msg);
        if e != nil {
            delete(mes.clients, conn.RemoteAddr().String());
        }
    }
}

func (mes *multiEchoServer) handleListen() {
    for {
        conn, e := mes.listener.Accept();
        if e != nil {
            if mes.stop == true {
                fmt.Println("listen stop")
                return
            }
            continue;
        }
        if _, ok := mes.clientMsgs[conn.RemoteAddr().String()]; !ok {
            mes.clientMsgs[conn.RemoteAddr().String()] = make(chan []byte, 100);
        }
        mes.clients[conn.RemoteAddr().String()] = conn;
        go mes.handleConn(conn);
    }
}

func (mes *multiEchoServer) handleDistributeMsg()  {
    for {
        line := <- mes.read;
        for key, _ := range mes.clients {
            msgChan := mes.clientMsgs[key];
            if len(msgChan) < 100 {
                msgChan <- line;
            }
        }
    }
}