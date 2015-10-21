package bufferbci

import (
	"encoding/binary"
	"fmt"
)

func (c *Connection) GetHeader() (*Header, error) {
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

	resp, err := c.readMessageDefinition()
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

	h := &Header{}
	err = binary.Read(c.textConnection.R, ByteOrder, &h.nchans)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &h.nsamples)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &h.nevents)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &h.fsamp)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &h.data_type)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &h.bufsize)
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

	resp, err := c.readMessageDefinition()
	if err != nil {
		return err
	}
	if resp.Cmd != CommandGetOk {
		return fmt.Errorf("expected GET_OK (0x204), but received (%#x)", resp.Cmd)
	}

	if resp.Bufsize > 0 {
		additionalResponse := make([]byte, resp.Bufsize)
		n, err := c.textConnection.R.Read(additionalResponse)
		if err != nil {
			return err
		}
		if uint32(n) != resp.Bufsize {
			return fmt.Errorf("expected %d bytes, but received only %d", resp.Bufsize, n)
		}
	}
	// TODO: process samples

	return nil
}
