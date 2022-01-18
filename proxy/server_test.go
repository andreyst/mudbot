package proxy

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"log"
	"mudbot/botutil"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

var tearedDown bool
var mudListener net.Listener

func tearUp() (clientConn net.Conn, mudConn net.Conn, server *server) {
	tearedDown = false

	// TODO: Run tests with random ports
	proxyAddr := "localhost:17844"
	mudAddr := "localhost:17845"

	mudWg := new(sync.WaitGroup)
	mudWg.Add(1)

	fmt.Printf("Running mud listen\n")
	var listenErr error
	mudListener, listenErr = net.Listen("tcp", mudAddr)
	if listenErr != nil {
		log.Fatalf("listenErr=%v\n", listenErr)
	}

	go func() {
		var acceptErr error
		mudConn, acceptErr = mudListener.Accept()
		if acceptErr != nil && !tearedDown {
			log.Fatalf("acceptErr=%v\n", acceptErr)
		}
		fmt.Printf("Mud accepted connection\n")
		mudWg.Done()
	}()

	server = NewServer()
	go server.Start(proxyAddr, mudAddr)

	fmt.Printf("Connecting to proxy\n")
	clientConn, err := net.Dial("tcp", proxyAddr)
	if clientConn == nil {
		log.Fatalf("proxy dial failed: %v\n", err)
	}

	fmt.Printf("Waiting for Mud to accept\n")
	mudWg.Wait()

	return
}

func tearDown(clientConn net.Conn, mudConn net.Conn, s *server) {
	tearedDown = true

	fmt.Printf("%s stopping server\n", time.Now())
	s.Stop(true)
	fmt.Printf("closing client\n")
	clientConn.Close()
	fmt.Printf("closing mud\n")
	mudConn.Close()
	fmt.Printf("closing mud listener\n")
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
	mudConn.Write(append([]byte("STRING "), []byte{TELNET_IAC, TELNET_CMD_GO_AHEAD}...))

	// expect fully flushed string

	readFromConnUntilTimeout(clientConn, t)

	mudConn.Write(botutil.Multiappend(
		[]byte("STRING1 "),
		[]byte{TELNET_IAC, TELNET_CMD_GO_AHEAD},
		[]byte("STRING2 "),
		[]byte{TELNET_IAC, TELNET_CMD_GO_AHEAD},
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
	in := append([]byte(strInUncompressed), compressionStartSequence...)
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
	mudConn.Write(append(compressionStartSequence, compressedBytes...))

	strOut := readFromConnUntilTimeout(clientConn, t)
	checkSentRcvdDiffers(strIn, strOut, t)

	var strInBuilder strings.Builder
	for i := 0; i < 10; i++ {
		strIn := "Hohoyo" + fmt.Sprintf("%d", i)
		strInBuilder.WriteString(strIn)

		compressedBytes2 := compressStr(zlibWriter, strIn, compressedBuffer, t)
		n, mudWriteErr := mudConn.Write(compressedBytes2)
		if mudWriteErr != nil {
			fmt.Printf("mudWriteErr=%v\n", mudWriteErr)
		}
		fmt.Printf("Wrote \"%s\" (n=%v) to mudConn\n", strIn, n)
	}

	strOut2 := readFromConnUntilTimeout(clientConn, t)
	checkSentRcvdDiffers(strInBuilder.String(), strOut2, t)

	tearDown(clientConn, mudConn, server)
}
