package main

import (
	"atomicgo.dev/cursor"
	"flag"
	"fmt"
	"time"
)

const refresh = time.Millisecond * 100

func main() {
	// parse flags to get username
	var name string
	flag.StringVar(&name, "n", "", "Specify name. Default is NEO")
	flag.Parse()
	// start the actual function
	do(name)
}

func do(name string) {
	// hide cursor and show it again before leaving
	cursor.Hide()
	defer cursor.Show()

	// c channel is for sending letter
	// d channel is for sending, that you are done
	c := make(chan rune)
	d := make(chan bool)

	// start the updater into the background
	go updater(c, d, name)

	// the output is what we print every time after the refresh time is over, it will grow with time from the updater
	out := ""
	// goup (meant as "go up") is how much space (in lines) is above and below the message line
	goup := 17
	// when we have a linebreak, we increase the lineset, to always go to appropriate line when printing
	lineset := 0
	// infinite loop where we keep receiving new characters from the updater
	for {
		// what did we receive from the updater if anything
		select {
		// we received a character
		case a := <-c:
			// what is this character about?
			switch a {
			case 8:
				// we received a backspace, we reset the output and reset the output and the lineset
				out = ""
				lineset = 0
			case '\n':
				// we received a linebreak, increase lineset, so we can accomodate, that the output will naturally be a line longer
				lineset++
				out = out + string(a)
			default:
				// nothing special? just append it to the output
				out = out + string(a)
			}
		// we received something on the d channel, so we are done
		case <-d:
			cursor.ClearLinesUp(goup * 2)
			return
		default:
			// the updater did not have an update for us, we don't change output and print the same output again
		}

		// Now we print out output again
		cursor.ClearLine()

		cursor.ClearLinesUp(goup * 2)
		cursor.Down(goup)

		cursor.StartOfLine()
		fmt.Print(out)
		for i := 0; i < goup-lineset; i++ {
			fmt.Println()
		}
		// wait before printing the next time
		time.Sleep(refresh)
	}
}

// the updater sends new characters for the output
func updater(c chan rune, d chan bool, name string) {
	// default name
	if name == "" {
		name = "Neo"
	}
	// message
	// \x08 is backspace to reset the message
	// \n is a linebreak
	out := "Wake Up, " + name + "...\x08The Matrix Has You\x08Follow the white rabbit.\x08Knock, Knock, " + name + "."

	// we continously send the characters of the message
	for i := 0; i < len(out); i++ {
		time.Sleep(getWait(out[i]) * time.Millisecond)
		c <- rune(out[i])
	}
	// after a final wait we are done
	time.Sleep(time.Duration(2000) * time.Millisecond)
	d <- true
}

// how much do we wait depending on what we print next
func getWait(input byte) time.Duration {
	switch input {
	case '\x08':
		return time.Duration(2000)
	case '\n':
		return time.Duration(1500)
	default:
		return time.Duration(300)
	}
}
