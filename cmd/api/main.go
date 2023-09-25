package main

const port = ":9920"

func main() {
	srv := NewServer()
	srv.Run(port)
}
