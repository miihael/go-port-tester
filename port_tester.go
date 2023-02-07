package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Server defines the minimum contract our
// TCP and UDP server implementations must satisfy.
type Server interface {
	Run(q chan bool) error
	Close() error
}

// NewServer creates a new Server using given protocol
// and addr.
func NewServer(protocol, addr string) (Server, error) {
	switch strings.ToLower(protocol) {
	case "tcp":
		return &TCPServer{
			addr: addr,
		}, nil
	case "udp":
		return &UDPServer{
			addr: addr,
		}, nil
	}
	return nil, errors.New("Invalid protocol given")
}

// TCPServer holds the structure of our TCP
// implementation.
type TCPServer struct {
	addr   string
	server net.Listener
}

// Run starts the TCP Server.
func (t *TCPServer) Run(quit chan bool) (err error) {
	t.server, err = net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	defer t.Close()

	return t.handleConnections(quit)
}

// Close shuts down the TCP Server
func (t *TCPServer) Close() (err error) {
	return t.server.Close()
}

// handleConnections is used to accept connections on
// the TCPServer and handle each of them in separate
// goroutines.
func (t *TCPServer) handleConnections(quit chan bool) (err error) {
	for {
		select {
		case <-quit:
			return nil
		default:
			conn, err := t.server.Accept()
			if err != nil || conn == nil {
				err = errors.New("could not accept connection")
				break
			}
			log.Printf("# incoming connection from %s", conn.RemoteAddr())

			go t.handleConnection(conn)
		}
	}
	return nil
}

// handleConnections deals with the business logic of
// each connection and their requests.
func (t *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		req, err := rw.ReadString('\n')
		if err != nil {
			rw.WriteString("failed to read input")
			rw.Flush()
			return
		}

		rw.WriteString(fmt.Sprintf("Request received: %s", req))
		rw.Flush()
	}
}

// UDPServer holds the necessary structure for our
// UDP server.
type UDPServer struct {
	addr   string
	server *net.UDPConn
}

// Run starts the UDP server.
func (u *UDPServer) Run(q chan bool) (err error) {
	laddr, err := net.ResolveUDPAddr("udp", u.addr)
	if err != nil {
		return errors.New("could not resolve UDP addr")
	}

	u.server, err = net.ListenUDP("udp", laddr)
	if err != nil {
		return errors.New("could not listen on UDP")
	}

	return u.handleConnections(q)
}

func (u *UDPServer) handleConnections(quit chan bool) error {
	var err error
	for {
		select {
		case <-quit:
			return nil
		default:
			buf := make([]byte, 2048)
			n, conn, err := u.server.ReadFromUDP(buf)
			if err != nil {
				log.Println(err)
				break
			}
			if conn == nil {
				continue
			}

			go u.handleConnection(conn, buf[:n])
		}
	}
	return err
}

func (u *UDPServer) handleConnection(addr *net.UDPAddr, cmd []byte) {
	u.server.WriteToUDP([]byte(fmt.Sprintf("Request recieved: %s", cmd)), addr)
}

// Close ensures that the UDPServer is shut down gracefully.
func (u *UDPServer) Close() error {
	return u.server.Close()
}

var waitgroup sync.WaitGroup

func checkPort(ip string, port int, proto string, timeoutSecs int) {
	defer waitgroup.Done()

	host := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout(proto, host, time.Duration(timeoutSecs)*time.Second)

	if err != nil {
		fmt.Printf("%s: Error %s\n", host, err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "TEST\n")

	reply := make([]byte, 4096)
	n, err := conn.Read(reply)

	if err != nil {
		fmt.Printf("%s: Error %s\n", host, err)
		return
	}

	fmt.Printf("%s: OK %d\n", host, n)
}

func main() {
	var port = flag.Int("port", 0, "Port to listen and check")
	var timeout = flag.Int("timeout", 30, "Timeout for reply")
	var proto = flag.String("proto", "tcp", "Protocol (tcp/udp)")
	var bind = flag.String("bind", "", "Bind to specific IP, empty by default")
	var sleep = flag.Int("sleep", 30, "Time to sleep after checks")
	var delay = flag.Int("delay", 15, "Time to sleep before starting checks")
	var noListen = flag.Bool("no-listen", false, "Do not start local servers, only check remote")
	flag.Parse()
	if *port == 0 {
		log.Fatal("Please specify port")
	}
	if flag.NArg() < 1 {
		log.Fatal("Please specify at least one IP to check")
	}

	var s Server
	var err error
	q := make(chan bool)
	if !*noListen {
		go func() {
			s, err = NewServer(*proto, fmt.Sprintf("%s:%d", *bind, *port))
			if err != nil {
				log.Fatal("Error starting server on %s:%d", *bind, *port)
			}
			s.Run(q)
		}()

		time.Sleep(time.Duration(*delay) * time.Second)
	}
	waitgroup.Add(flag.NArg())
	for _, ip := range flag.Args() {
		go checkPort(ip, *port, *proto, *timeout)
	}
	waitgroup.Wait()
	//log.Printf("sleeping %d sec", *sleep)
	time.Sleep(time.Duration(*sleep) * time.Second)
	close(q)
	//log.Printf("quit")
}
