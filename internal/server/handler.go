package server

import (
	"bufio"
	"io"
	"net"
	"os"
	"strconv"

	"stash/internal/protocol"
)

func HandleConn(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	defer w.Flush()

	h, err := protocol.ReadHeader(r)
	if err != nil {
		return
	}

	switch h.Command {
	case "BACKUP":
		handleBackup(h, r, w)
	case "RESTORE":
		handleRestore(h, w)
	case "LIST":
		handleList(w)
	default:
		protocol.WriteError(w, "unknown command")
	}
}

func handleBackup(h *protocol.Header, r *bufio.Reader, w *bufio.Writer) {
	name := h.Fields["name"]
	size, _ := strconv.ParseInt(h.Fields["size"], 10, 64)

	f, offset, err := OpenForWrite(name)
	if err != nil {
		protocol.WriteError(w, "open failed")
		return
	}
	defer f.Close()

	protocol.WriteOK(w, offset)
	w.Flush()

	remaining := size - offset
	if remaining > 0 {
		io.CopyN(f, r, remaining)
	}
}

func handleRestore(h *protocol.Header, w *bufio.Writer) {
	name := h.Fields["name"]
	reqOffset, _ := strconv.ParseInt(h.Fields["offset"], 10, 64)

	f, size, err := OpenForRead(name)
	if err != nil {
		protocol.WriteError(w, "not found")
		return
	}
	defer f.Close()

	if reqOffset > size {
		reqOffset = 0
	}

	protocol.WriteOK(w, reqOffset)
	w.Flush()

	f.Seek(reqOffset, io.SeekStart)
	io.Copy(w, f)
}

func handleList(w *bufio.Writer) {
	ensureDataDir()

	entries, err := os.ReadDir(dataDir)
	if err != nil {
		protocol.WriteError(w, "read failed")
		return
	}

	w.WriteString("OK\n\n")
	for _, e := range entries {
		if !e.IsDir() {
			w.WriteString(e.Name() + "\n")
		}
	}
}
