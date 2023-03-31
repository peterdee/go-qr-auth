## go-qr-auth

This is a server for [QR Auth](https://github.com/peterdee/qr-auth) mobile application

Stack: [Gorilla Websocket](https://github.com/gorilla/websocket)

This application also features a custom ping mechanism that automatically disconnects unresponsive clients

### Deploy

```shell script
git clone https://github.com/peterdee/go-qr-auth
cd ./go-qr-auth
gvm use go1.19
go mod download
```

### Environment variables

The `.env` file is required if launching locally, see [.env.example](./.env.example) for details

### Launch

```shell script
go run ./
```

Can be used with [AIR](https://github.com/cosmtrek/air)

### License

[MIT](./LICENSE.md)
