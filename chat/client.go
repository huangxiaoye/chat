//package chat

//import (
//	"bufio"
//	"log"
//	"net"
//)

//type Client struct {
//	conn     net.Conn
//	incoming Message
//	outgoing Message
//	reader   *bufio.Reader
//	writer   *bufio.Writer
//	quiting  chan net.Conn
//	name     string
//}

//func (self *Client) GetName() string {
//	return self.name
//}

//func (self *Client) SetName(name string) {
//	self.name = name
//}

//func (self *Client) GetIncoming() string {
//	return <-self.incoming
//}

//func (self *Client) PutOutgoing(message string) {
//	self.outgoing <- message
//}

//func CreateClient(conn net.Conn) *Client {
//	reader := bufio.NewReader(conn)
//	writer := bufio.NewWriter(conn)

//	client := &Client{
//		conn:     conn,
//		incoming: make(Message),
//		outgoing: make(Message),
//		quiting:  make(chan net.Conn),
//		reader:   reader,
//		writer:   writer,
//	}

//	client.Listen()
//	return client
//}

//func (self *Client) Listen() {
//	go self.Read()
//	go self.Write()
//}

//func (self *Client) quit() {
//	self.quiting <- self.conn
//}

//func (self *Client) Read() {
//	for {
//		if line, _, err := self.reader.ReadLine(); err == nil {
//			self.incoming <- string(line)
//		} else {
//			log.Printf("Read error: %s\n", err)
//			self.quit()
//			return
//		}
//	}

//}

//func (self *Client) Write() {
//	for data := range self.outgoing {
//		if _, err := self.writer.WriteString(data + "\n"); err != nil {
//			self.quit()
//			return
//		}

//		if err := self.writer.Flush(); err != nil {
//			log.Printf("Write error: %s\n", err)
//			self.quit()
//			return
//		}
//	}
//}

//func (self *Client) Close() {
//	self.conn.Close()
//}

package chat

import (
	"bufio"
	"log"
	"net"
)

type Client struct {
	conn     net.Conn
	incoming Message
	outgoing Message
	name     string
	quiting  chan net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func (self *Client) GetIncoming() string {
	return <-self.incoming
}

func (self *Client) GetName() string {
	return self.name
}

func (self *Client) SetName(name string) {
	self.name = name
}

func (self *Client) PutOutgoing(message string) {
	self.outgoing <- message
}

func CreateClient(conn net.Conn) *Client {

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	c := &Client{
		conn:     conn,
		incoming: make(Message),
		outgoing: make(Message),
		quiting:  make(chan net.Conn),
		reader:   reader,
		writer:   writer,
	}

	c.Listen()
	return c
}

func (self *Client) Listen() {
	go self.Read()
	go self.Write()
}

func (self *Client) Read() {
	for {
		//func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
		if line, _, err := self.reader.ReadLine(); err == nil {
			self.incoming <- string(line)
		} else {
			log.Printf("Read error: %s\n", err)
			self.quit()
			return
		}
	}
}

func (self *Client) Write() {
	for data := range self.outgoing {
		//func (b *Writer) WriteString(s string) (int, error)
		if _, err := self.writer.WriteString(data + "\n"); err != nil {
			log.Printf("Write error: %s\n", err)
			self.quit()
			return
		}

		if err := self.writer.Flush(); err != nil {
			log.Printf("Flush error: %s\n", err)
			self.quit()
			return
		}
	}
}

func (self *Client) quit() {
	self.quiting <- self.conn
}

func (self *Client) Close() {
	self.conn.Close()
}
