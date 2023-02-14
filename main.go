package main

func main() {
	server := NewApiServer("localhost:5000")
	server.Run()
}
