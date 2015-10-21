package bufferbci

import (
	"encoding/binary"
	"net"
	"net/textproto"
)

var (
	ConnType  = "tcp"
	ByteOrder = binary.LittleEndian
)

type Connection struct {
	rawConnection  *net.TCPConn
	textConnection *textproto.Conn
}

func Connect(host string) (*Connection, error) {
	tcpaddr, err := net.ResolveTCPAddr(ConnType, host)
	if err != nil {
		return nil, err
	}

	rawConn, err := net.DialTCP(ConnType, nil, tcpaddr)
	if err != nil {
		return nil, err
	}

	conn := &Connection{rawConnection: rawConn}
	conn.textConnection = textproto.NewConn(rawConn)

	return conn, nil
}

func (c *Connection) sendMessageDefinition(def *messageDefinition) error {
	var err error

	err = binary.Write(c.textConnection.W, ByteOrder, def.Version)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, ByteOrder, def.Cmd)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, ByteOrder, def.Bufsize)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) readMessageDefinition() (*messageDefinition, error) {
	def := &messageDefinition{}
	var err error

	err = binary.Read(c.textConnection.R, ByteOrder, &def.Version)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &def.Cmd)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &def.Bufsize)
	if err != nil {
		return nil, err
	}

	return def, nil
}

func (c *Connection) Close() error {
	return c.textConnection.Close()
}

const (
	CommandPutHdr uint16 = 0x101
	CommandPutDat uint16 = 0x102
	CommandPutOk  uint16 = 0x104
	CommandPutErr uint16 = 0x105

	CommandGetHdr uint16 = 0x201
	CommandGetDat uint16 = 0x202
	CommandGetOk  uint16 = 0x204
	CommandGetErr uint16 = 0x205
)

type messageDefinition struct {
	Version uint16
	Cmd     uint16
	Bufsize uint32
}

type Header struct {
	nchans    uint32  // number of channels
	nsamples  uint32  // number of samples, must be 0 for PUT_HDR
	nevents   uint32  // number of events, must be 0 for PUT_HDR
	fsamp     float32 // sampling frequency (Hz)
	data_type uint32  // type of the sample data (see table above)
	bufsize   uint32  // size of remaining parts of the message in bytes (total size of all chunks)

	chunks []HeaderChunk
}

type HeaderChunk struct {
}

type sampleDefinition struct {
	nchans    uint32 // number of channels
	nsamples  uint32 // number of samples
	data_type uint32 // type of the samples
	bufsize   uint32 // number of remaining bytes in the message, i.e. size of the data samples
}
