package main

func main() {
	s := InitServer()
	s.Run("127.0.0.1:8088")
}
