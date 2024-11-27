# small-chatroom

A small concept of a chatroom in Go

## Usage

```bash
go run main.go hub.go client.go handler.go 
```

In a separate terminal (TMUX would be really good for this) open a WebSocket client.

I was using the `wscat` node package.

```shell
wscat -c ws://localhost:8080/ws
```
Now that you are connected you can send messages!

