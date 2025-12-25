package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"stash/internal/protocol"
)

func Backup(addr, path string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	fi, _ := os.Stat(path)

	fmt.Fprintf(w,
		"BACKUP\nname:%s\nsize:%d\n\n",
		filepath.Base(path), fi.Size(),
	)
	w.Flush()

	offset := protocol.ReadOffset(r)

	f, _ := os.Open(path)
	defer f.Close()
	f.Seek(offset, io.SeekStart)

	io.Copy(w, f)
}

func Restore(addr, name string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var offset int64 = 0
	if st, err := os.Stat(name); err == nil {
		offset = st.Size()
	}

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	fmt.Fprintf(w,
		"RESTORE\nname:%s\noffset:%d\n\n",
		name, offset,
	)
	w.Flush()

	confirmedOffset := protocol.ReadOffset(r)

	f, _ := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	f.Seek(confirmedOffset, 0)
	io.Copy(f, r)
}

func List(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	w.WriteString("LIST\n\n")
	w.Flush()

	h, _ := protocol.ReadHeader(r)
	if h.Command != "OK" {
		return
	}

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		fmt.Print(line)
	}
}
