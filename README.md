1.build server.exe
go build -o server.exe main.go user.go server.go

2.build client.exe
go build -o client.exe client.go

3.run server.exe and client.exe
.\server.exe
.\client.exe