package main

import "os"

func main() {
	os.Exit(New().Run())
}

type ExitStatus = int

const (
	ExitSuccess ExitStatus = 0
	ExitFail    ExitStatus = 1
)

func (s *Server) Run() ExitStatus {
	if err := s.handler(); err != nil {
		return ExitFail
	}
	return ExitSuccess
}
