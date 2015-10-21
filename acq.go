package gobci

import (
	"encoding/binary"
	"fmt"
)

// GetHeader retrieves some useful meta-information from the current connection
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
	err = binary.Read(c.textConnection.R, ByteOrder, &h.NChannels)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &h.NSamples)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &h.NEvents)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &h.SamplingFrequency)
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

	if h.bufsize > 0 {
		buffy := make([]byte, h.bufsize)
		n, err := c.textConnection.R.Read(buffy)
		if err != nil {
			return nil, err
		}
		if uint32(n) != h.bufsize {
			return nil, fmt.Errorf("expected %d bytes, but received only %d", h.bufsize, n)
		}
		//fmt.Println("Received additional header information:", buffy)
	}

	return h, nil
}

// GetData retrieves a 2d-array with sample data; one []float64 per sample
func (c *Connection) GetData(begin, end uint32) ([][]float64, error) {
	def := &messageDefinition{1, CommandGetDat, 0}

	if end > 0 {
		def.Bufsize = 8
	}

	err := c.sendMessageDefinition(def)
	if err != nil {
		return nil, err
	}

	if end > 0 {
		err = binary.Write(c.textConnection.W, ByteOrder, begin)
		if err != nil {
			return nil, err
		}
		err = binary.Write(c.textConnection.W, ByteOrder, end)
		if err != nil {
			return nil, err
		}
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
		return nil, fmt.Errorf("expected GET_OK (0x204), but received (%v)", resp)
	}

	if resp.Bufsize < 16 {
		return nil, fmt.Errorf("expected a Bufsize >= 16")
	}

	// Read and parse sample definition
	sd := &sampleDefinition{}
	err = binary.Read(c.textConnection.R, ByteOrder, &sd.nchans)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &sd.nsamples)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &sd.data_type)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.textConnection.R, ByteOrder, &sd.bufsize)
	if err != nil {
		return nil, err
	}

	if sd.nsamples > 0 && sd.bufsize == 0 {
		return nil, fmt.Errorf("expected a Bufsize > 0 in the fixed definition")
	}

	// Read and parse the samples
	samples := make([][]float64, sd.nsamples)
	for sample := uint32(0); sample < sd.nsamples; sample++ {
		samples[sample] = make([]float64, sd.nchans)
		for ch := uint32(0); ch < sd.nchans; ch++ {
			switch sd.data_type {
			case DataTypeChar:
				r, size, err := c.textConnection.R.ReadRune()
				if err != nil {
					return nil, err
				}
				if size != 1 {
					return nil, fmt.Errorf("could not reliably read rune; size was %d", size)
				}
				samples[sample][ch] = float64(r)
			case DataTypeUint8:
				var placeholder uint8
				err = binary.Read(c.textConnection.R, ByteOrder, &placeholder)
				if err != nil {
					return nil, err
				}
				samples[sample][ch] = float64(placeholder)
			case DataTypeUint16:
				var placeholder uint16
				err = binary.Read(c.textConnection.R, ByteOrder, &placeholder)
				if err != nil {
					return nil, err
				}
				samples[sample][ch] = float64(placeholder)
			case DataTypeuint32:
				var placeholder uint32
				err = binary.Read(c.textConnection.R, ByteOrder, &placeholder)
				if err != nil {
					return nil, err
				}
				samples[sample][ch] = float64(placeholder)
			case DataTypeUint64:
				var placeholder uint64
				err = binary.Read(c.textConnection.R, ByteOrder, &placeholder)
				if err != nil {
					return nil, err
				}
				samples[sample][ch] = float64(placeholder)
			case DataTypeInt8:
				var placeholder int8
				err = binary.Read(c.textConnection.R, ByteOrder, &placeholder)
				if err != nil {
					return nil, err
				}
				samples[sample][ch] = float64(placeholder)
			case DataTypeInt16:
				var placeholder int16
				err = binary.Read(c.textConnection.R, ByteOrder, &placeholder)
				if err != nil {
					return nil, err
				}
				samples[sample][ch] = float64(placeholder)
			case DataTypeInt32:
				var placeholder int32
				err = binary.Read(c.textConnection.R, ByteOrder, &placeholder)
				if err != nil {
					return nil, err
				}
				samples[sample][ch] = float64(placeholder)
			case DataTypeInt64:
				var placeholder int64
				err = binary.Read(c.textConnection.R, ByteOrder, &placeholder)
				if err != nil {
					return nil, err
				}
				samples[sample][ch] = float64(placeholder)
			case DataTypeFloat32:
				var placeholder float32
				err = binary.Read(c.textConnection.R, ByteOrder, &placeholder)
				if err != nil {
					return nil, err
				}
				samples[sample][ch] = float64(placeholder)
			case DataTypeFloat64:
				var placeholder float64
				err = binary.Read(c.textConnection.R, ByteOrder, &placeholder)
				if err != nil {
					return nil, err
				}
				samples[sample][ch] = placeholder
			default:
				return nil, fmt.Errorf("unknown sample data type: %d", sd.data_type)
			}
		}
	}

	return samples, nil
}
