package server

import (
	"bufio"
	"io"
	"net"
	"os"

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
	size := protocol.MustInt(h, "size")

	f, offset, err := OpenForWrite(name)
	if err != nil {
		protocol.WriteError(w, "open failed")
		return
	}
	defer f.Close()

	if offset > size {
		offset = 0
	}

	remaining := size - offset
	protocol.WriteOK(w, map[string]int64{
		"offset": offset,
		"size":   size,
	})
	w.Flush()

	n, err := io.CopyN(f, r, remaining)
	if err != nil || n != remaining {
		return
	}
}

func handleRestore(h *protocol.Header, w *bufio.Writer) {
	name := h.Fields["name"]
	reqOffset := protocol.MustInt(h, "offset")

	f, size, err := OpenForRead(name)
	if err != nil {
		protocol.WriteError(w, "not found")
		return
	}
	defer f.Close()

	if reqOffset > size {
		reqOffset = 0
	}

	remaining := size - reqOffset
	protocol.WriteOK(w, map[string]int64{
		"offset": reqOffset,
		"size":   size,
	})
	w.Flush()

	f.Seek(reqOffset, io.SeekStart)

	n, err := io.CopyN(w, f, remaining)
	if err != nil || n != remaining {
		return
	}
}

func handleList(w *bufio.Writer) {
	ensureDataDir()

	entries, err := os.ReadDir(dataDir)
	if err != nil {
		protocol.WriteError(w, "read failed")
		return
	}

	count := int64(0)
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}

	protocol.WriteOK(w, map[string]int64{
		"count": count,
	})
	w.Flush()

	for _, e := range entries {
		if !e.IsDir() {
			w.WriteString(e.Name() + "\n")
		}
	}
}
