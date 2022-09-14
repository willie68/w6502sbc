package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/willie68/w6502sbc/tree/main/software/emulator/internal/logging"
	"github.com/willie68/w6502sbc/tree/main/software/emulator/pkg/emulator"

	flag "github.com/spf13/pflag"
)

var binFile string

func init() {
	// variables for parameter override
	flag.StringVarP(&binFile, "bin", "b", "", "this is the path and filename to the ROM image file")
}

func main() {
	log.Logger.Info("W6502SBC Emulator")
	flag.Parse()
	if binFile == "" {
		log.Logger.Error("no ROM given.")
		os.Exit(-1)
	}
	log.Logger.Info("ROM Image in : " + binFile)
	dat, err := os.ReadFile(binFile)
	if err != nil {
		log.Logger.Errorf("can't read ROM: %v", err)
		os.Exit(-1)
	}
	w6502sbc := emulator.NewEmu6502().WithROM(0xE000, dat).WithRAM(0x000, 0x7fff).Build()
	w6502sbc.Start()
	fmt.Printf("Adr: $%.4X, SP: $%.2X, A: $%.2X, X: $%.2X, Y: $%.2X\r\n", w6502sbc.Adr(), w6502sbc.SP(), w6502sbc.A(), w6502sbc.X(), w6502sbc.Y())
	fmt.Print("x for exit\r\n>")
	ch := make(chan byte)
	go func(ch chan byte) {
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- b[0]
		}
	}(ch)

	//var b []byte = make([]byte, 1)
	for {
		select {
		case c, _ := <-ch:
			//			os.Stdin.Read(b)
			if c == 'x' {
				break
			}
			switch c {
			case 'i':
				w6502sbc.IRQ()
			case 'n':
				w6502sbc.NMI()
			case '\r':
				break
			default:
				w6502sbc.Step()
				fmt.Printf("Adr: $%.4X, SP: $%.2X, A: $%.2X, X: $%.2X, Y: $%.2X\r\n", w6502sbc.Adr(), w6502sbc.SP(), w6502sbc.A(), w6502sbc.X(), w6502sbc.Y())
				fmt.Print("x for exit\r\n>")
			}
		default:
		}
		time.Sleep(time.Millisecond * 100)
	}
}
