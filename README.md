# Caesar Cipher

How to Run :

Set up Receiver and Sender 

Init package (if u havent) :

```
go mod init
```

### Receiver 

Listen on some port, example :

```
go run main.go -mode=receiver -addr=:9000 -shift=3
```

### Sender 

Fill all the flag, for example :

```
go run main.go -mode=sender -addr="127.0.0.1:9000" -shift=3
```

You will prompted to input the text, and send those text by press `enter`


### How To Run

Find the address of receiver computer/devices. This can be done with :

```
// from receiver device
hostname -I
```

to test :

```
ping <<IP_ADDRESS>>
```

in Ubuntu, (I use VM because it's much more simple)

after that, use above instructions on receiver and sender to send and receive the message