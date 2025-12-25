# stash

`stash` is a **TCP-based command-line backup and restore tool** written in Go.

It is a systems-focused learning project that uses **raw TCP** and a **custom application-layer protocol**, with emphasis on correctness, concurrency, and resource control. No HTTP, frameworks, or external libraries are used.

## Features (v1)
- Backup (upload files)
- Restore (download files)
- List stored files
- Resumable transfers
- Concurrent clients with bounded resources

## Status
Work in progress.

## License
MIT
