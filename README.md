# Caesar Cipher

How to Run :

Set up Receiver and Sender 

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