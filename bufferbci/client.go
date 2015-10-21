package bufferbci

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/textproto"
)

const (
	connType = "tcp"
)

var (
	byteOrder = binary.LittleEndian
)

type Connection struct {
	rawConnection  *net.TCPConn
	textConnection *textproto.Conn
}

func Connect(host string) (*Connection, error) {
	tcpaddr, err := net.ResolveTCPAddr(connType, host)
	if err != nil {
		return nil, err
	}

	rawConn, err := net.DialTCP(connType, nil, tcpaddr)
	if err != nil {
		return nil, err
	}

	conn := &Connection{rawConnection: rawConn}
	conn.textConnection = textproto.NewConn(rawConn)

	return conn, nil
}

func (c *Connection) SendDebugCommand() ([]byte, error) {
	header, err := c.getHeader()
	if err != nil {
		return nil, err
	}
	fmt.Println(header)

	return nil, nil
}

func (c *Connection) sendMessageDefinition(def *messageDefinition) error {
	var err error

	err = binary.Write(c.textConnection.W, byteOrder, def.Version)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, byteOrder, def.Cmd)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, byteOrder, def.Bufsize)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) putHeader(nChannels uint32, fSamp float32) error {
	def := &messageDefinition{1, CommandPutHdr, 24}
	err := c.sendMessageDefinition(def)

	h := header{nChannels, 0, 0, fSamp, 0, 0, nil}

	err = binary.Write(c.textConnection.W, byteOrder, h.nchans)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, byteOrder, h.nsamples)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, byteOrder, h.nevents)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, byteOrder, h.fsamp)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, byteOrder, h.data_type)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, byteOrder, h.bufsize)
	if err != nil {
		return err
	}

	for range h.chunks {
		// TODO: write chunk
	}

	return c.textConnection.W.Flush()
}

func (c *Connection) readResponse() (*messageDefinition, error) {
	def := &messageDefinition{}
	var err error

	err = binary.Read(c.textConnection.R, byteOrder, &def.Version)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, byteOrder, &def.Cmd)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, byteOrder, &def.Bufsize)
	if err != nil {
		return nil, err
	}

	return def, nil
}

func (c *Connection) getHeader() (*header, error) {
	// Request header information

	def := &messageDefinition{1, CommandGetHdr, 0}
	err := c.sendMessageDefinition(def)
	if err != nil {
		return nil, err
	}

	err = c.textConnection.W.Flush()
	if err != nil {
		return nil, err
	}

	resp, err := c.readResponse()
	if err != nil {
		return nil, err
	}
	if resp.Cmd != CommandGetOk {
		return nil, fmt.Errorf("expected GET_OK (0x204), but received (%#x)", resp.Cmd)
	}

	if resp.Bufsize == 0 {
		return nil, fmt.Errorf("expected header information")
	}

	// Read and parse header information

	h := &header{}
	err = binary.Read(c.textConnection.R, byteOrder, &h.nchans)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, byteOrder, &h.nsamples)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, byteOrder, &h.nevents)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, byteOrder, &h.fsamp)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, byteOrder, &h.data_type)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, byteOrder, &h.bufsize)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (c *Connection) getDat() error {
	def := &messageDefinition{1, CommandGetDat, 0}
	err := c.sendMessageDefinition(def)
	if err != nil {
		return err
	}

	err = c.textConnection.W.Flush()
	if err != nil {
		return err
	}

	resp, err := c.readResponse()
	if err != nil {
		return err
	}
	if resp.Cmd != CommandGetOk {
		return fmt.Errorf("expected GET_OK (0x204), but received (%#x)", resp.Cmd)
	}

	if resp.Bufsize > 0 {
		additionalResponse := make([]byte, resp.Bufsize)
		fmt.Println("Reading the additional", resp.Bufsize, "bytes")
		n, err := c.textConnection.R.Read(additionalResponse)
		if err != nil {
			return err
		}
		if uint32(n) != resp.Bufsize {
			return fmt.Errorf("expected %d bytes, but received only %d", resp.Bufsize, n)
		}
		fmt.Println("Received:", additionalResponse)
	}

	return nil
}

func (c *Connection) Close() error {
	return c.textConnection.Close()
}

type messageDefinition struct {
	Version uint16
	Cmd     uint16
	Bufsize uint32
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

type header struct {
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

type Request struct {
	definition messageDefinition
	header     header
}

type Response struct {
	definition messageDefinition
}
