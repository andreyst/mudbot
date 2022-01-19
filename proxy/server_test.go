package proxy

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"mudbot/botutil"
	"mudbot/telnet"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

var tearedDown bool
var mudListener net.Listener
var logger = botutil.NewLogger("test")

func tearUp() (clientConn net.Conn, mudConn net.Conn, server *server) {
	tearedDown = false

	// TODO: Run tests with random ports
	proxyAddr := "localhost:17844"
	mudAddr := "localhost:17845"

	mudWg := new(sync.WaitGroup)
	mudWg.Add(1)

	logger.Infof("Running mud listen")
	var listenErr error
	mudListener, listenErr = net.Listen("tcp", mudAddr)
	if listenErr != nil {
		logger.Fatalf("listenErr=%v", listenErr)
	}

	go func() {
		var acceptErr error
		mudConn, acceptErr = mudListener.Accept()
		if acceptErr != nil && !tearedDown {
			logger.Fatalf("acceptErr=%v", acceptErr)
		}
		logger.Infof("Mud accepted connection")
		mudWg.Done()
	}()

	server = NewServer(proxyAddr, mudAddr)
	go server.Start()

	logger.Infof("Connecting to proxy")
	clientConn, err := net.Dial("tcp", proxyAddr)
	if clientConn == nil {
		logger.Fatalf("proxy dial failed: %v", err)
	}

	logger.Infof("Waiting for Mud to accept")
	mudWg.Wait()

	return
}

func tearDown(clientConn net.Conn, mudConn net.Conn, s *server) {
	tearedDown = true

	logger.Infof("%s stopping server", time.Now())
	s.Stop(true)
	logger.Infof("closing client")
	clientConn.Close()
	logger.Infof("closing mud")
	mudConn.Close()
	logger.Infof("closing mud listener")
	mudListener.Close()
}

func readFromConn(conn net.Conn, t *testing.T) (out string, isTimeout bool) {
	buf := make([]byte, 2<<10)
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	n, readErr := conn.Read(buf)
	if netErr, ok := readErr.(net.Error); ok && netErr.Timeout() {
		isTimeout = true
		return
	}
	if readErr != nil {
		t.Errorf("Error reading from clientConn, err=%v\n", readErr)
	}
	out = string(buf[:n])
	return
}

func readFromConnUntilTimeout(conn net.Conn, t *testing.T) (out string) {
	var strOutBuilder strings.Builder
	for {
		strOut, isTimeout := readFromConn(conn, t)
		if isTimeout {
			break
		}
		strOutBuilder.WriteString(strOut)
	}
	return strOutBuilder.String()
}

func checkSentRcvdDiffers(strIn string, strOut string, t *testing.T) {
	if strIn != strOut {
		t.Errorf("Received does not equal sent\n(sent len=%v, rvcd len=%v):\n%v\n%v\n",
			len(strIn),
			len(strOut),
			strIn,
			strOut)
	}
}

func TestTearUp(t *testing.T) {
	clientConn, mudConn, server := tearUp()

	strIn := "Mud says yohohoho"
	mudConn.Write([]byte(strIn))
	strOut := readFromConnUntilTimeout(clientConn, t)
	checkSentRcvdDiffers(strIn, strOut, t)

	tearDown(clientConn, mudConn, server)
}

func TestMudToClient(t *testing.T) {
	clientConn, mudConn, server := tearUp()

	strIn := "Mud says yohohoho"
	mudConn.Write([]byte(strIn))
	strOut := readFromConnUntilTimeout(clientConn, t)
	checkSentRcvdDiffers(strIn, strOut, t)

	tearDown(clientConn, mudConn, server)
}

func TestClientToMud(t *testing.T) {
	clientConn, mudConn, server := tearUp()

	strIn := "Client says yohohoho"
	clientConn.Write([]byte(strIn))
	strOut := readFromConnUntilTimeout(mudConn, t)
	checkSentRcvdDiffers(strIn, strOut, t)

	tearDown(clientConn, mudConn, server)
}

func compressStr(zlibWriter *zlib.Writer, str string, compressedBuffer *bytes.Buffer, t *testing.T) []byte {
	_, zlibWriteErr := zlibWriter.Write([]byte(str))
	if zlibWriteErr != nil {
		t.Errorf("zlibWriteErr=%v\n", zlibWriteErr)
	}

	zlibFlushErr := zlibWriter.Flush()
	if zlibFlushErr != nil {
		t.Errorf("zlibFlushErr=%v\n", zlibFlushErr)
	}

	compressedBytes := make([]byte, 100000)
	n, readErr := compressedBuffer.Read(compressedBytes)
	if readErr != nil {
		t.Errorf("readErr=%v\n", readErr)
	}

	return compressedBytes[:n]
}

func TestAccumulatorFlush(t *testing.T) {
	clientConn, mudConn, server := tearUp()

	mudConn.Write([]byte("THIS "))
	time.Sleep(10 * time.Millisecond)
	mudConn.Write([]byte("IS "))
	time.Sleep(10 * time.Millisecond)
	mudConn.Write([]byte("A "))
	time.Sleep(10 * time.Millisecond)
	mudConn.Write(append([]byte("STRING "), telnet.GaSequence...))

	// expect fully flushed string

	readFromConnUntilTimeout(clientConn, t)

	mudConn.Write(botutil.Multiappend(
		[]byte("STRING1 "),
		telnet.GaSequence,
		[]byte("STRING2 "),
		telnet.GaSequence,
	))

	// expect two separate flushes

	readFromConnUntilTimeout(clientConn, t)

	// case with IAC GA split between sends

	// case with uncompressed + IAC GA in compressed

	tearDown(clientConn, mudConn, server)
}

func TestCompressed(t *testing.T) {
	clientConn, mudConn, server := tearUp()

	strInUncompressed := "UncompressedHohoho"

	compressedBuffer := new(bytes.Buffer)
	zlibWriter := zlib.NewWriter(compressedBuffer)

	strInToCompress := "Yohohoho"

	compressedBytes := compressStr(zlibWriter, strInToCompress, compressedBuffer, t)
	in := append([]byte(strInUncompressed), telnet.CompressionStartSequence...)
	in = append(in, compressedBytes...)
	mudConn.Write(in)

	strOut := readFromConnUntilTimeout(clientConn, t)
	checkSentRcvdDiffers(strInUncompressed+strInToCompress, strOut, t)

	strIn2 := "Hohoyo"

	compressedBytes2 := compressStr(zlibWriter, strIn2, compressedBuffer, t)
	mudConn.Write(compressedBytes2)

	strOut2 := readFromConnUntilTimeout(clientConn, t)
	checkSentRcvdDiffers(strIn2, strOut2, t)

	tearDown(clientConn, mudConn, server)
}

func TestCompressedMany(t *testing.T) {
	clientConn, mudConn, server := tearUp()

	compressedBuffer := new(bytes.Buffer)
	zlibWriter := zlib.NewWriter(compressedBuffer)

	strIn := "start"

	compressedBytes := compressStr(zlibWriter, strIn, compressedBuffer, t)
	mudConn.Write(append(telnet.CompressionStartSequence, compressedBytes...))

	strOut := readFromConnUntilTimeout(clientConn, t)
	checkSentRcvdDiffers(strIn, strOut, t)

	var strInBuilder strings.Builder
	for i := 0; i < 10; i++ {
		strIn := "Hohoyo" + fmt.Sprintf("%d", i)
		strInBuilder.WriteString(strIn)

		compressedBytes2 := compressStr(zlibWriter, strIn, compressedBuffer, t)
		n, mudWriteErr := mudConn.Write(compressedBytes2)
		if mudWriteErr != nil {
			logger.Infof("mudWriteErr=%v", mudWriteErr)
		}
		logger.Infof("Wrote \"%s\" (n=%v) to mudConn", strIn, n)
	}

	strOut2 := readFromConnUntilTimeout(clientConn, t)
	checkSentRcvdDiffers(strInBuilder.String(), strOut2, t)

	tearDown(clientConn, mudConn, server)
}
