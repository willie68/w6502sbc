package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/willie68/w6502sbc/tree/main/software/emulator/internal/config"
	log "github.com/willie68/w6502sbc/tree/main/software/emulator/internal/logging"
	"github.com/willie68/w6502sbc/tree/main/software/emulator/pkg/emulator"
	"gopkg.in/yaml.v3"

	flag "github.com/spf13/pflag"
)

var binFile string
var c02 bool
var cfn string
var sbcConfig map[string]any

func init() {
	// variables for parameter override
	flag.StringVarP(&binFile, "bin", "b", "", "this is the path and filename to the ROM image file")
	flag.StringVarP(&cfn, "cfg", "c", "", "this is the path and filename to the ROM image file")
	flag.BoolVar(&c02, "c02", false, "emulate a 65C02")
}

func main() {
	log.Logger.Info("W6502SBC Emulator")
	flag.Parse()
	if c02 {
		log.Logger.Info("using CMOS 6502 Version")
	}
	if binFile == "" && cfn == "" {
		log.Logger.Error("no ROM or config given.")
		os.Exit(-1)
	}
	b := emulator.NewEmu6502()
	if cfn != "" {
		b = processConfiguration(cfn, b)
	} else {
		log.Logger.Info("ROM Image in : " + binFile)
		dat, err := os.ReadFile(binFile)
		if err != nil {
			log.Logger.Errorf("can't read ROM: %v", err)
			os.Exit(-1)
		}
		b = b.WithROM(0xE000, dat).WithRAM(0x000, 0x7fff)
		if c02 {
			b.With65C02()
		}
	}

	w6502sbc := b.Build()
	str := w6502sbc.Start()
	fmt.Printf("%s\r\n", str)
	fmt.Printf("Adr: $%.4X, SP: $%.2X, A: $%.2X, X: $%.2X, Y: $%.2X, S: %08b\r\n", w6502sbc.Adr(), w6502sbc.SP(), w6502sbc.A(), w6502sbc.X(), w6502sbc.Y(), w6502sbc.ST())
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
	end := false
	for !end {
		select {
		case c := <-ch:
			if c == 'x' {
				end = true
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
				res := w6502sbc.Step()
				fmt.Println(res)
				fmt.Printf("Adr: $%.4X, SP: $%.2X, A: $%.2X, X: $%.2X, Y: $%.2X, S: %08b\r\n", w6502sbc.Adr(), w6502sbc.SP(), w6502sbc.A(), w6502sbc.X(), w6502sbc.Y(), w6502sbc.ST())
				fmt.Print("x for exit\r\n>")
			}
		default:
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func processConfiguration(cfn string, b *emulator.Build6502) *emulator.Build6502 {
	if _, err := os.Stat(cfn); errors.Is(err, os.ErrNotExist) {
		log.Logger.Alertf("file: \"%s\" does not exists", cfn)
		os.Exit(-1)
	}
	log.Logger.Infof("reading configuration: %s", cfn)
	yamlFile, err := ioutil.ReadFile(cfn)
	if err != nil {
		log.Logger.Fatalf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &sbcConfig)
	if err != nil {
		log.Logger.Fatalf("Unmarshal: %v", err)
	}
	//fmt.Printf("%v", sbcConfig)
	for k, v := range sbcConfig {
		m := v.(map[string]any)
		t, ok := m["type"]
		if !ok {
			log.Logger.Errorf("config: missing type of key %s", k)
		}
		if strings.ToLower(k) == "cpu" {
			switch strings.ToLower(t.(string)) {
			case "65C02", "65c02":
				b.With65C02()
				log.Logger.Info("switing CPU: using 65C02\r\n")
			}
		}
		switch strings.ToLower(t.(string)) {
		case "ram":
			ram := config.RAM{}
			if err := mergo.Map(&ram, m); err != nil {
				log.Logger.Errorf("config: RAM config error %v", err)
			}
			b.WithRAM(uint16(ram.Start), uint16(ram.Start+ram.Length-1))
			log.Logger.Infof("adding RAM : start: $%.4x, end: $%.4x, length: $%.4x\r\n", ram.Start, ram.Start+ram.Length-1, ram.Length)
		case "rom":
			rom := config.ROM{}
			if err := mergo.Map(&rom, m); err != nil {
				log.Logger.Errorf("config: ROM config error: %v", err)
				os.Exit(-1)
			}
			binFile, _ := m["file"].(string)
			log.Logger.Info("config: ROM Image in : " + binFile)
			data, err := os.ReadFile(binFile)
			if err != nil {
				log.Logger.Errorf("can't read ROM: %v", err)
				os.Exit(-1)
			}

			b.WithROM(uint16(rom.Start), data)
			log.Logger.Infof("adding ROM : start: $%.4x, end: $%.4x, length: $%.4x\r\n", rom.Start, rom.Start+len(data)-1, len(data))
		case "6522":
			via := config.P6522{}
			if err := mergo.Map(&via, m); err != nil {
				log.Logger.Errorf("config: VIA config error: %v", err)
				os.Exit(-1)
			}
			log.Logger.Infof("adding VIA : start: $%.4x\r\n", via.Start)
		case "6551":
			acia := config.P6551{}
			if err := mergo.Map(&acia, m); err != nil {
				log.Logger.Errorf("config: ACIA config error: %v", err)
				os.Exit(-1)
			}
			log.Logger.Infof("adding ACIA: start: $%.4x\r\n", acia.Start)
		}
	}
	return b
}
