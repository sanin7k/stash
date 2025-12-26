# stash

`stash` is a **TCP-based command-line backup and restore tool** written in Go.

It is a systems-focused learning project built on **raw TCP** with a **custom application-layer protocol**, emphasizing **correctness, concurrency, and explicit protocol invariants**.  
No HTTP, frameworks, or external libraries are used.

---

## Features (v1)

- File backup (upload)
- File restore (download)
- Server-side file listing
- Resumable transfers using explicit file offsets
- Concurrent clients with bounded resource usage

---

## Design Notes

- Uses a **request–response protocol** over persistent TCP connections
- Separates **control messages** from bulk data transfer
- Enforces protocol invariants to prevent partial writes and duplication
- Designed for correctness and failure awareness over feature breadth

---

## Usage

### Run the server

```bash
go run cmd/stashd/main.go
```
Starts the TCP server and listens for client connections.

### Run the client

The client is invoked via its entrypoint with an operation argument.

#### List files stored on the server
```bash
go run cmd/stash/main.go list
```

#### Backup (upload) a file
```bash
go run cmd/stash/main.go backup <local-file-path>
```
Uploads the specified file to the server.
If a partial upload exists, the transfer resumes from the last known offset.

#### Restore (download) a file
```bash
go run cmd/stash/main.go restore <remote-file-name>
```
Downloads the specified file from the server.
If a partial download exists locally, the transfer resumes from the last known offset.

## Status
v1 complete.
The core protocol and behavior are stable; future versions may extend functionality
without breaking v1 semantics.

## License
MIT
