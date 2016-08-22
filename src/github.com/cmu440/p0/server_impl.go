// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
    "net"
    "bufio"
    "fmt"
    "strconv"
)

type multiEchoServer struct {
    num int
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
    mes := multiEchoServer{num:0}
    return mes;
}

func (mes *multiEchoServer) Start(port int) error {
    if port > 65535 || port <= 0 {
        fmt.Println("port is invalid");
        panic("port is invalid");
    }
    listener, e := net.ListenTCP("tcp", ":" + port);
    if e != nil {
        panic("listen failed");
    }
    go func() {
        for {
            conn, err := listener.Accept();
            if err != nil {
                continue;
            }
            mes.num++;
            go func() {
                defer conn.Close();
                var buf [512]byte;
                for {
                    len, e := conn.Read(buf);
                    if e != nil {
                        break;
                    }
                    if len > 0 {
                        conn.Write(buf);
                    }
                }

            }()
        }
    }()
    return nil;
}

func (mes *multiEchoServer) Close() {
    // TODO: implement this!
}

func (mes *multiEchoServer) Count() int {
    return mes.num;
}

// TODO: add additional methods/functions below!
