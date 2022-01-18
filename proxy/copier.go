package proxy

import (
	"bytes"
	"compress/zlib"
	"io"
	"mudbot/botutil"
	"mudbot/parser"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
)

const bufSize = 65536

type accumulationPolicy int

const (
	ACCUMULATION_POLICY_DO accumulationPolicy = iota
	ACCUMULATION_POLICY_DONT
)

type copierState int

const (
	STATE_START copierState = iota
	STATE_TELNET_IAC
	STATE_TELNET_SB
	STATE_TELNET_MCCPV2          // SE option was MCCPv2
	STATE_TELNET_MCCPV2_IAC      // IAC after IAC SE MCCPv2
	STATE_TELNET_COMMAND_WITH_OP // Other unimportant commands with op
)

type Copier struct {
	logger *zap.SugaredLogger

	state copierState

	parser *parser.Parser

	accumulationPolicy accumulationPolicy

	done chan struct{}
	wg   sync.WaitGroup

	buf         []byte
	readBytes   []byte
	accumulator []byte
	idx         int

	compressed      bool
	compressedBytes []byte
	zlibInBuffer    *bytes.Buffer
	zlibOutBuffer   []byte
	zlibReader      io.ReadCloser
}

func NewCopier(accumulate accumulationPolicy, parser *parser.Parser, logger *zap.SugaredLogger) *Copier {
	c := Copier{}
	c.logger = logger
	c.parser = parser
	c.accumulationPolicy = accumulate
	c.zlibInBuffer = new(bytes.Buffer)
	c.zlibOutBuffer = make([]byte, bufSize*5)
	c.done = make(chan struct{})
	c.wg.Add(1)

	return &c
}

func (c *Copier) Run(dst net.Conn, src net.Conn, workerDone chan struct{}) {
	// copier.copy(dst, src, workerDone, logger)
	c.copy(dst, src, workerDone)
	c.logger.Debugf("stopping copy")
	close(c.done)
	c.wg.Done()
}

func (c *Copier) copy(dst net.Conn, src net.Conn, workerDone chan struct{}) {
	c.buf = make([]byte, bufSize)

	for {
		select {
		case <-workerDone:
			c.logger.Infof("Worker done")
			return
		default:
		}

		c.idx = 0

		n, isEof, isStop, readErr := c.read(src, c.buf)
		if readErr != nil {
			c.logger.Fatalf("Got error from src: %v", readErr)
		}
		if isStop {
			break
		}
		if n == 0 {
			continue
		}

		c.readBytes = c.buf[:n]

		if c.compressed {
			decompressErr := c.decompress()
			if decompressErr != nil {
				c.logger.Fatalf("decompressErr=%v", decompressErr)
			}
		}

		afterFlushIdx := 0
		for c.idx < len(c.readBytes) {
			b := c.readBytes[c.idx]
			switch c.state {
			case STATE_START:
				switch b {
				case TELNET_IAC:
					c.state = STATE_TELNET_IAC
				}
			case STATE_TELNET_IAC:
				switch b {
				case TELNET_CMD_GO_AHEAD:
					if c.accumulationPolicy == ACCUMULATION_POLICY_DO {
						appendBytes := c.readBytes[afterFlushIdx : c.idx+1-len(gaSequence)]
						c.accumulator = append(c.accumulator, appendBytes...)
						c.flushAccumulator()
						afterFlushIdx = c.idx + 1
					}
					c.state = STATE_START

				case TELNET_CMD_SB:
					c.state = STATE_TELNET_SB
				}
			case STATE_TELNET_SB:
				switch b {
				case TELNET_OPT_MCCPV2:
					c.state = STATE_TELNET_MCCPV2
				}
			case STATE_TELNET_MCCPV2:
				switch b {
				case TELNET_IAC:
					c.state = STATE_TELNET_MCCPV2_IAC
				}
			case STATE_TELNET_MCCPV2_IAC:
				switch b {
				case TELNET_CMD_SE:
					c.compressedBytes = c.readBytes[c.idx+1:]
					// Cut out compression start sequence
					// since we are always sending decompressed stream to client
					c.readBytes = c.readBytes[:c.idx+1-len(compressionStartSequence)]
					c.compressed = true
					c.idx = len(c.readBytes)
					c.state = STATE_START
				}
			}

			c.idx++
		}

		c.logger.Debugf("Data str:\n%v", string(c.readBytes))
		c.logger.Debugf("Data hex:\n%v", botutil.PrintHex(c.readBytes))

		telnetCmds := GetTelnetCommandsStrings(c.readBytes)
		if len(telnetCmds) > 0 {
			c.logger.Debugf("Telnet commands: %v", telnetCmds)
		}

		c.writeToConn(c.readBytes, dst)

		if c.accumulationPolicy == ACCUMULATION_POLICY_DO {
			c.accumulator = append(c.accumulator, c.readBytes[afterFlushIdx:]...)
			c.logger.Debugf("Accumulator str:\n%v", string(c.accumulator))
		}

		if isEof {
			return
		}
	}
}

func (c *Copier) decompress() (err error) {
	err = nil

	c.zlibInBuffer.Write(c.readBytes)

	if c.zlibReader == nil {
		c.zlibReader, err = zlib.NewReader(c.zlibInBuffer)
		if err != nil {
			return
		}
	}

	c.readBytes = make([]byte, 0)

	for c.zlibInBuffer.Len() > 0 {
		var n int
		n, err = c.zlibReader.Read(c.zlibOutBuffer)
		if err != nil {
			return
		}
		c.readBytes = append(c.readBytes, c.zlibOutBuffer[:n]...)
	}

	return
}

func (c *Copier) read(conn net.Conn, buf []byte) (n int, isEof bool, isStop bool, err error) {
	if len(c.compressedBytes) > 0 {
		n = copy(buf, c.compressedBytes)
		c.compressedBytes = []byte{}
		return
	}

	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	n, readErr := conn.Read(buf)

	isEof = readErr == io.EOF
	if readErr != nil && !isEof {
		if readErr == io.ErrClosedPipe {
			c.logger.Infof("Src pipe closed")
			isStop = true
			return
		}

		if netErr, ok := readErr.(net.Error); ok && netErr.Timeout() {
			n = 0
			return
		}

		err = readErr
		return
	}

	return
}

func (c *Copier) writeToConn(buf []byte, conn net.Conn) {
	n, writeErr := conn.Write(buf)
	if writeErr != nil {
		c.logger.Fatalf("Error while writing to conn: %v", writeErr)
	}
	c.logger.Debugf("Wrote string of len %v to conn", n)
}

func (c *Copier) flushAccumulator() {
	c.logger.Debugf("Flushing accumulator content:%v", string(c.accumulator))
	c.logger.Debugf("Flushing accumulator hex:%v", botutil.PrintHex(c.accumulator))
	c.parser.Parse(c.accumulator)
	c.accumulator = []byte{}
}

//go:generate stringer -type=accumulationPolicy,copierState -output copier_string.go
