package main

type Command struct {
	Id			  int		`json:"id"`
	CommandString string 	`json:"commandstring"`
}

type Commands []Command