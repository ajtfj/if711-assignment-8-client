package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Event string

const (
	ENTER Event = "enter"
	LEAVE Event = "leave"
)

func enter(client *rpc.Client, c *Client) error {
	args := SushiBarEnterArgs{
		Client: c,
	}
	reply := SushiBarEnterReply{}
	if err := client.Call("SushiBar.Enter", args, &reply); err != nil {
		return err
	}

	return nil
}

func leave(client *rpc.Client, ticket int) error {
	args := SushiBarLeaveArgs{
		Ticket: ticket,
	}
	reply := SushiBarLeaveReply{}
	if err := client.Call("SushiBar.Leave", args, &reply); err != nil {
		return err
	}

	fmt.Println(reply.Farewell)

	return nil
}

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatal("undefined HOST")
	}

	client, err := rpc.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	var event Event
	for {
		fmt.Printf("Do you want to enter or leave the sushibar? (enter/leave) ")
		fmt.Scanf("%s\n", &event)
		switch event {
		case ENTER:
			fmt.Printf("What is your name? ")
			c := Client{}
			fmt.Scanf("%s\n", &c.Name)

			fmt.Printf("Requesting entrace of client %s\n", c.Name)
			go func() {
				if err := enter(client, &c); err != nil {
					fmt.Println(err)
				}
			}()

		case LEAVE:
			fmt.Printf("What is your ticket? ")
			var ticket int
			fmt.Scanf("%d\n", &ticket)

			if err := leave(client, ticket); err != nil {
				fmt.Println(err)
			}

		default:
			fmt.Println("Undefined event")
		}
	}
}

type Client struct {
	Name string
}

type SushiBarEnterArgs struct {
	Client *Client
}

type SushiBarEnterReply struct {
	Ticket int
}

type SushiBarLeaveArgs struct {
	Ticket int
}

type SushiBarLeaveReply struct {
	Farewell string
}
