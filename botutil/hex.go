package botutil

import (
	"bytes"
	"log"
	"os/exec"
)

func ByteToHex(buf []byte) string {
	cmd := exec.Command("xxd")
	cmd.Stdin = bytes.NewReader(buf)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	return out.String()
}

func StrToHex(s string) string {
	return ByteToHex([]byte(s))
}
