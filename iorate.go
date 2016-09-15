// Provides io.Reader and io.Writer wrappers that artificially limit
// the transmission rate.
package iorate

/*
	We assume that data transfer occurs in small time slices, each
	lasting for time 'tau'. The time then can be represented as index
	't' = 0, 1, 2, ... The amount of data transferred during slice 't'
	is 'd_t'. If the transfer speed is limited to 'L', then for any 't'
	the following holds:

		d_t/tau <= L.

	Then the Write and Read functions below simply split their time in
	fixed parts of length 'tau' and take care not to pass more than
	L*tau data during each part.
*/

import (
	"io"
	"time"
)

const tau = 100 // ms

type Rate int64

// Rate units. 1 * Bps is one byte per second.
// Note that there are two kinds of semantics. Units with small "b"
// assume bits and SI units for "K", "M" and "G", whereas "B" versions
// assume bytes and power-of-two units for "K", "M" and "G".
const (
	// bytes and data semantics
	Bps  Rate = 1
	KBps      = 1024 * Bps
	MBps      = 1024 * KBps
	GBps      = 1024 * MBps
	// bits and rate semantics
	Kbps = 8 * 1000 * Bps
	Mbps = 1000 * Kbps
	Gbps = 1000 * Mbps
)

type writer struct {
	out         io.Writer
	maxSendSize int
}

type reader struct {
	in          io.Reader
	maxReadSize int
}

// Returns a writer limited to 'maxSpeed' bytes per second.
func NewWriter(out io.Writer, maxSpeed Rate) *writer {
	t := new(writer)
	t.out = out
	t.maxSendSize = int(int64(maxSpeed) * int64(tau) / 1000)
	return t
}

// Returns a reader limited to 'maxSpeed' bytes per second.
func NewReader(in io.Reader, maxSpeed Rate) *reader {
	t := new(reader)
	t.in = in
	t.maxReadSize = int(int64(maxSpeed) * int64(tau) / 1000)
	return t
}

// Implements the io.Read function.
func (t *reader) Read(b []byte) (n int, err error) {
	max := cap(b)

	// Maximum receive size we can do in 'tau' time
	readSize := max
	if readSize > t.maxReadSize {
		readSize = t.maxReadSize
	}

	dt := time.Duration(tau) * time.Millisecond
	err = nil
	n = 0
	end := 0

	for n < max {
		time.Sleep(dt)

		end = n + readSize
		if end > max {
			end = max
		}
		read, err := t.in.Read(b[n:end])
		n += read
		if err != nil {
			break
		}
	}
	return n, err
}

// Implements the io.Write function.
func (t *writer) Write(b []byte) (n int, err error) {
	total := len(b)

	// Maximum send size we can do in 'tau' time
	sendSize := total
	if sendSize > t.maxSendSize {
		sendSize = t.maxSendSize
	}

	dt := time.Duration(tau) * time.Millisecond
	err = nil

	// Bounds of the data portion being sent
	pos := 0
	end := 0

	for pos < total {
		time.Sleep(dt)

		end = pos + sendSize
		if end > total {
			end = total
		}

		sent, err := t.out.Write(b[pos:end])
		pos += sent
		if err != nil {
			break
		}
	}

	return pos, err
}
