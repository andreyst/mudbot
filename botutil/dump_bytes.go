package botutil

import "os"

func DumpBytes(b []byte) {
	fo, err := os.Create("/tmp/bytes")
	if err != nil {
		panic(err)
	}
	fo.Write(b)
	fo.Close()
}
