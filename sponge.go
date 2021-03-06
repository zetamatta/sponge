package main

import (
	"fmt"
	"io"
	"os"
)

type OutputT struct {
	Fd      *os.File
	TmpName string
	Fname   string
}

func main() {
	outputList := make([]*OutputT, 0, len(os.Args)-1)
	for _, fname := range os.Args[1:] {
		tmpName := fname + ".sponge"
		fd, err := os.Create(tmpName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", tmpName, err.Error())
		} else {
			outputList = append(outputList, &OutputT{
				Fd:      fd,
				TmpName: tmpName,
				Fname:   fname,
			})
		}
	}

	buffer := make([]byte, 256)
	for {
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		if n == 0 {
			break
		}
		for _, p := range outputList {
			p.Fd.Write(buffer[:n])
		}
	}
	os.Stdin.Close()

	for _, p := range outputList {
		p.Fd.Close()
		err := os.Remove(p.Fname)
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		err = os.Rename(p.TmpName, p.Fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}
