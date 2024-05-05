package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/warthog618/modem/at"
	"github.com/warthog618/modem/serial"
)

type Pin struct {
	pin   int
	err   error
	value bool
}

func (p *Pin) open() {
	cmd := exec.Command("gpio", "mode", fmt.Sprintf("%d", p.pin), "output")
	p.err = cmd.Run()
}

func (p *Pin) high() {
	p.value = true
	cmd := exec.Command("gpio", "write", fmt.Sprintf("%d", p.pin), "1")
	p.err = cmd.Run()
}

func (p *Pin) low() {
	p.value = false
	cmd := exec.Command("gpio", "write", fmt.Sprintf("%d", p.pin), "0")
	p.err = cmd.Run()
}
func (p *Pin) state() bool {
	return p.value
}

var (
	// Use mcu pin 10, corresponds to physical pin 19 on the pi
	pinRed   = Pin{pin: 23}
	pinGreen = Pin{pin: 24}
	pinGsm   = Pin{pin: 21}
)

// var ttys = []string{"/dev/ttyS1", "/dev/ttyS2", "/dev/ttyS3", "/dev/ttyS4", "/dev/ttyS5", "/dev/ttyS6", "/dev/ttyS7"}
var ttys = []string{"/dev/ttyS4"}

func gsmtest() {
	pinGsm.open()
	pinGsm.high()
	for {
		for _, v := range ttys {
			fmt.Printf("test pin gsm %v device %s\n", pinGsm.state(), v)
			if err := gsmwork(v, 115200); err != nil {
				fmt.Printf("%s ...is not sucsessed! \n", err.Error())
			} else {
				fmt.Println("is sucsess...")
				os.Exit(0)
			}
		}
		if pinGsm.state() {
			pinGsm.low()
		} else {
			pinGsm.high()
		}
	}

}

func gsmwork(dev string, baud int) error {
	timeout := 1 * time.Second
	m, err := serial.New(serial.WithPort(dev), serial.WithBaud(baud))
	if err != nil {
		return err
	}
	var mio io.ReadWriter = m
	a := at.New(mio, at.WithTimeout(timeout))
	err = a.Init()
	if err != nil {
		return err
	}
	cmds := []string{
		"I",
		"+GCAP",
		"+CMEE=2",
		"+CGMI",
		"+CGMM",
		"+CGMR",
		"+CGSN",
		"+CSQ",
		"+CIMI",
		"+CREG?",
		"+CNUM",
		"+CPIN?",
		"+CEER",
		"+CSCA?",
		"+CSMS?",
		"+CSMS=?",
		"+CPMS=?",
		"+CCID?",
		"+CCID=?",
		"^ICCID?",
		"+CNMI?",
		"+CNMI=?",
		"+CNMA=?",
		"+CMGF?",
		"+CMGF=?",
		"+CUSD?",
		"+CUSD=?",
		"^USSDMODE?",
		"^USSDMODE=?",
	}
	for _, cmd := range cmds {
		info, err := a.Command(cmd)
		fmt.Println("AT" + cmd)
		if err != nil {
			fmt.Printf(" %s\n", err)
			continue
		}
		for _, l := range info {
			fmt.Printf(" %s\n", l)
		}
	}
	return nil
}
func main() {
	pinRed.open()
	pinGreen.open()
	go gsmtest()
	w := true
	for {
		if w {
			pinGreen.high()
			pinRed.low()
		} else {
			pinGreen.low()
			pinRed.high()
		}
		w = !w
		time.Sleep(1000 * time.Millisecond)
	}
}
