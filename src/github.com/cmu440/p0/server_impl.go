// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
    "net"
    "bufio"
    "fmt"
)

type multiEchoServer struct {
    clients          map[string] net.Conn
    clientMsgs       map[string] chan([] byte)
    listener         net.Listener
    readChan         chan []byte
    closeCliChan     chan string
    addCliChan       chan net.Conn
    countRequestChan chan chan int
    closeChan        chan bool
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
    mes := &multiEchoServer{
        clients: make(map[string] net.Conn),
        clientMsgs: make(map[string] chan([]byte)),
        readChan: make(chan([]byte)),
        closeCliChan: make(chan string),
        addCliChan: make(chan net.Conn),
        countRequestChan: make(chan chan int),
        closeChan: make(chan bool),
    }
    return MultiEchoServer(mes);
}

func (mes *multiEchoServer) Start(port int) error {
    listener, e := net.Listen("tcp", fmt.Sprintf(":%d", port));
    if e != nil {
        fmt.Println("listen tcp failed");
        return e;
    }
    mes.listener = listener;
    go mes.handleEvent();
    go mes.handleAccept();
    return nil;
}

func (mes *multiEchoServer) Close() {
    mes.closeChan <- true;
}

func (mes *multiEchoServer) Count() int {
    countChan := make(chan int);
    mes.countRequestChan<-countChan;
    return <-countChan;
}

// TODO: add additional methods/functions below!
func (mes *multiEchoServer) handleConn(conn net.Conn) {
    go mes.handleWrite(conn);
    reader := bufio.NewReader(conn);
    for {
        line, e := reader.ReadBytes('\n');
        if  e != nil{
            mes.closeCliChan <-conn.RemoteAddr().String();
            return;
        }
        mes.readChan <- line;
    }
}

func (mes *multiEchoServer) handleWrite(conn net.Conn) {
    msgChan := mes.clientMsgs[conn.RemoteAddr().String()];
    for {
        msg := <- msgChan;
        _, e := conn.Write(msg);
        if e != nil {
            mes.closeCliChan <-conn.RemoteAddr().String();
            return;
        }
    }
}

func (mes *multiEchoServer) handleAccept() {
    for {
        conn, e := mes.listener.Accept();
        if e != nil {
            fmt.Println("listen stop")
            return
        }
        mes.addCliChan<-conn;
    }
}

func (mes *multiEchoServer) handleEvent() {
    for {
        select {
        case line := <-mes.readChan:
            for key, _ := range mes.clients {
                msgChan := mes.clientMsgs[key];
                if len(msgChan) < 100 {
                    msgChan <- line;
                }
            }
        case addr := <-mes.closeCliChan:
            if conn, ok := mes.clients[addr]; ok {
                conn.Close();
                delete(mes.clients, addr);
            }
        case countChan := <-mes.countRequestChan:
            cnt := len(mes.clients);
            countChan<-cnt;
        case conn := <-mes.addCliChan:
            if _, ok := mes.clientMsgs[conn.RemoteAddr().String()]; !ok {
                mes.clientMsgs[conn.RemoteAddr().String()] = make(chan []byte, 100);
            }
            mes.clients[conn.RemoteAddr().String()] = conn;
            go mes.handleConn(conn);
        case <-mes.closeChan:
            mes.listener.Close();
            for _, conn := range mes.clients {
                conn.Close();
            }
            return;
        }
    }
}