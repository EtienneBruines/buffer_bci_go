package gobci

import (
	"encoding/binary"
	"net"
	"net/textproto"
)

const (
	DefaultHost = "localhost:1972"
)

var (
	ConnType  = "tcp"
	ByteOrder = binary.BigEndian
)

type Connection struct {
	rawConnection  *net.TCPConn
	textConnection *textproto.Conn
}

// Connect attempts a connection to the Buffer Server at given host
// Host can be in the format of "localhost:1972", and may be empty to use the default hostname.
func Connect(host string) (*Connection, error) {
	if len(host) == 0 {
		host = DefaultHost
	}

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

// Close closes the TCP connection with the host
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

	CommandFlushDat uint16 = 0x302
	CommandFlushOk  uint16 = 0x304
	CommandFlushErr uint16 = 0x305

	CommandWaitDat uint16 = 0x402
	CommandWaitOk  uint16 = 0x404
	CommandWaitErr uint16 = 0x405
)

const (
	DataTypeChar uint32 = iota
	DataTypeUint8
	DataTypeUint16
	DataTypeuint32
	DataTypeUint64
	DataTypeInt8
	DataTypeInt16
	DataTypeInt32
	DataTypeInt64
	DataTypeFloat32
	DataTypeFloat64
)

type messageDefinition struct {
	Version uint16
	Cmd     uint16
	Bufsize uint32
}

type Header struct {
	NChannels         uint32  // number of channels
	NSamples          uint32  // number of samples, must be 0 for PUT_HDR
	NEvents           uint32  // number of events, must be 0 for PUT_HDR
	SamplingFrequency float32 // sampling frequency (Hz)
	data_type         uint32  // type of the sample data (see table above)
	bufsize           uint32  // size of remaining parts of the message in bytes (total size of all chunks)

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
