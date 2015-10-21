package bufferbci

import "encoding/binary"

func (c *Connection) putHeader(nChannels uint32, fSamp float32) error {
	def := &messageDefinition{1, CommandPutHdr, 24}
	err := c.sendMessageDefinition(def)

	h := Header{nChannels, 0, 0, fSamp, 0, 0, nil}

	err = binary.Write(c.textConnection.W, ByteOrder, h.nchans)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, ByteOrder, h.nsamples)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, ByteOrder, h.nevents)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, ByteOrder, h.fsamp)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, ByteOrder, h.data_type)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, ByteOrder, h.bufsize)
	if err != nil {
		return err
	}

	for range h.chunks {
		// TODO: write chunk
	}

	return c.textConnection.W.Flush()
}
