package chat

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Message chan string
type ClientTable map[net.Conn]*Client
type Token chan int

const (
	MAXCLIENTS = 50
)

type Server struct {
	listener net.Listener
	incoming Message
	outgoing Message
	quiting  chan net.Conn
	clients  ClientTable
	tokens   Token
	pending  chan net.Conn
}

func (self *Server) generateToken() {
	self.tokens <- 0
}

func (self *Server) takeToken() {
	<-self.tokens
}

func CreateServer() *Server {
	s := &Server{
		incoming: make(Message),
		outgoing: make(Message),
		clients:  make(ClientTable, MAXCLIENTS),
		tokens:   make(Token, MAXCLIENTS),
		pending:  make(chan net.Conn),
		quiting:  make(chan net.Conn),
	}

	s.listen()
	return s
}

func (self *Server) listen() {
	go func() {
		for {
			select {
			case message := <-self.incoming:
				self.broadcast(message)
			case conn := <-self.pending:
				self.join(conn)
			case conn := <-self.quiting:
				self.leave(conn)
			}
		}
	}()
}

func (self *Server) join(conn net.Conn) {
	client := CreateClient(conn)
	//name := fmt.Sprintf("hax %d", time.Now().Unix())
	name := getUniqName()
	client.SetName(name)
	self.clients[conn] = client
	log.Printf("Auto assigned name for conn %p: %s\n", conn, name)
	go func() {
		for {
			msg := <-client.incoming
			log.Printf("Accepted message from conn %p: %s\n", conn, name)

			if strings.HasPrefix(msg, ":") {
				if cmd, err := parseCommand(msg); err == nil {
					if err = self.executeCommand(client, cmd); err == nil {
						continue
					} else {
						log.Println(err.Error())
					}
				} else {
					log.Println(err.Error())
				}
			}
			self.incoming <- fmt.Sprintf("%s says: %s", client.GetName(), msg)
		}
	}()

	go func() {
		for {
			conn := <-client.quiting
			log.Printf("Client %s is quiting\n", client.GetName())
			self.quiting <- conn
		}
	}()

	//go func() {
	//	for {
	//		select {
	//		case msg := <-client.incoming:
	//			log.Printf("Accepted message from conn %p: %s\n", conn, name)

	//			if strings.HasPrefix(msg, ":") {
	//				if cmd, err := parseCommand(msg); err == nil {
	//					if err = self.executeCommand(client, cmd); err == nil {
	//						continue
	//					} else {
	//						log.Println(err.Error())
	//					}
	//				} else {
	//					log.Println(err.Error())
	//				}
	//			}
	//			self.incoming <- fmt.Sprintf("%s says: %s", client.GetName(), msg)
	//		case conn := <-client.quiting:
	//			log.Printf("Client %s is quiting\n", client.GetName())
	//			self.quiting <- conn
	//		}

	//	}
	//}()

}

func (self *Server) broadcast(message string) {
	log.Printf("Broadcasting message: %s\n", message)
	for _, client := range self.clients {
		client.outgoing <- message
	}
}

func (self *Server) leave(conn net.Conn) {
	log.Printf("Somebody leaving now\n")
	conn.Close()
	delete(self.clients, conn)
	self.generateToken()
}

func (self *Server) Start(connString string) {
	self.listener, _ = net.Listen("tcp", connString)

	log.Printf("Server %p starts\n", self)

	for i := 0; i < MAXCLIENTS; i++ {
		self.generateToken()
	}

	for {
		conn, err := self.listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("A new connection %v kicks\n", conn)
		self.takeToken()
		self.pending <- conn
	}
}

func (self *Server) Stop() {
	self.listener.Close()
}
