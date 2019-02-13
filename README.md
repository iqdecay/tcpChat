# A TCP chat server

## Usage :

- First install dependencies :
     `go get github.com/jroimartin/gocui`
- Then launch the server :
    `go run server.go`                                                                         
- And for each client : 
    `go run client.go`


## TODO :
### Client :
- [x] Print messages in history
- [x] gocui package to have a nicer UI
- [x] Modify receiveMessage
- [x] Autorefresh messages
- [x] Scrollable history
- [x] Merge client and UI
- [x] Don't display sent messages twice
- [ ] Fix infinite invalid pseudo loop
- [ ] Display the users
- [ ] Get pseudo before connecting
- [ ] Add commands like leave, pseudo

### Server : 
- [x] Check for duplicate pseudo
- [x] Use regexp package to check the pseudo
- [ ] Send user list
- [ ] Attach pseudo to each message
- [ ] Add admin system

