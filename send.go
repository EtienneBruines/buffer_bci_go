package gobci

import (
	"encoding/binary"
	"fmt"
	"log"
)

func (c *Connection) putHeader(nChannels uint32, fSamp float32) error {
	def := &messageDefinition{1, CommandPutHdr, 24}
	err := c.sendMessageDefinition(def)

	h := Header{nChannels, 0, 0, fSamp, 0, 0, nil}

	err = binary.Write(c.textConnection.W, ByteOrder, h.NChannels)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, ByteOrder, h.NSamples)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, ByteOrder, h.NEvents)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, ByteOrder, h.SamplingFrequency)
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

func (c *Connection) FlushData() error {
	def := &messageDefinition{1, CommandFlushDat, 0}
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

	if resp.Cmd != CommandFlushOk {
		return fmt.Errorf("expected FLUSH_OK (0x304), but received %x", resp.Cmd)
	}

	return nil
}

func (c *Connection) PutEvent(key, value string) error {
	sampleCount, _, err := c.WaitData(0, 0, 0)
	if err != nil {
		return err
	}

	/*
		Fixed part
			type_type		uint32	0
			type_numel		uint32	6
			value_type		uint32	0
			value_numel		uint32	4
			sample			int32	10
			offset			int32	0
			duration		int32	0
			bufsize			uint32	10
	*/

	var (
		type_type   uint32 = DataTypeChar
		type_numel  uint32 = uint32(len(key))
		value_type  uint32 = DataTypeChar
		value_numel uint32 = uint32(len(value))
		sample      int32  = int32(sampleCount) - 1
		offset      int32  = 0
		duration    int32  = 0
		bufsize     uint32 = type_numel + value_numel
	)

	def := &messageDefinition{1, CommandPutEvt, 32 + bufsize}
	err = c.sendMessageDefinition(def)
	if err != nil {
		return err
	}

	err = binary.Write(c.textConnection.W, ByteOrder, type_type)
	if err != nil {
		return err
	}
	err = binary.Write(c.textConnection.W, ByteOrder, type_numel)
	if err != nil {
		return err
	}
	err = binary.Write(c.textConnection.W, ByteOrder, value_type)
	if err != nil {
		return err
	}
	err = binary.Write(c.textConnection.W, ByteOrder, value_numel)
	if err != nil {
		return err
	}
	err = binary.Write(c.textConnection.W, ByteOrder, sample)
	if err != nil {
		return err
	}
	err = binary.Write(c.textConnection.W, ByteOrder, offset)
	if err != nil {
		return err
	}
	err = binary.Write(c.textConnection.W, ByteOrder, duration)
	if err != nil {
		return err
	}
	err = binary.Write(c.textConnection.W, ByteOrder, bufsize)
	if err != nil {
		return err
	}

	_, err = c.textConnection.W.WriteString(key)
	if err != nil {
		return err
	}
	_, err = c.textConnection.W.WriteString(value)
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

	if resp.Cmd != CommandPutOk {
		return fmt.Errorf("expected PUT_OK (0x104), but received %x", resp.Cmd)
	}

	return nil
}

