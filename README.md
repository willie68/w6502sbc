# w6502sbc
willies 6502 single board computer
a 6502 sbc with the option to be used with a backplane as a module.
I'm a bit sorry, but there are so many 6502 projects out there that I've taken the liberty of documenting this project (for the time being) only in German.

Hier Dokumentiere ich den augenblicklichen Stand. Für die Historie gibt es den Blog. ([blog](/blog.md))

# Einleitung

Schon seit längerem beschäftigt mich ein altes Thema. Assembler und 6502. Mein C64 ist schon länger eingemottet, der VC20 verkauft. Trotzdem finde ich den 6502 von der Assemblersprache her sehr interessant. Was man da mit 1MHz und 64kb RAM alles machen konnte. Dieses Jahr (2022) war es dann soweit. Ich habe mich dazu durchgerungen einen 6502 SBC zu bauen. Im Internet gibt es viele, viele Beiträge und Vorschläge dazu. Leider gibt es keine Version, die man einfach nachbauen könnte und meine Erwartungen an einen 6502 Workbench Rechner erfüllt. Hier mal kurz meine Anforderungen:

- WDC65C02 CPU mit 4MHz (Takt per Programm umschaltbar)
- WDC65C22 VIA mit 4MHz, evtl. davon 2
- 6551 ASIC für die serielle Schnittstelle (R65C51) 
- 32Kb oder 64KB SRAM (beides war vorhanden) 
- EEPROM oder Flash als Kernal/Basic ROM
- NMI und IRQ belegbar
- variable Speicheraufteilung, evtl. auch von aussen modifizierbar
- Steckkartensystem, bzw. Systembus, am liebsten mit Corner Edge Connectors, wie im C64, VC20, eine passende Backplane wäre nett
- kein Video, Audio oder sowas
- SPI über 6522, damit man mal ein Display/Keyboard anschließen kann
- CPLD für irgendwas
- opt. LC-Display

Angefangen habe ich im Juli 22. Zunächst mit Bestellungen bei Mouser und Reichelt. 

Folgende Schritte zum Lernen:

- CLPD: dafür brauchte ich einen Programmieradapter (meine Wahl viel auf einen TL866II+)
- einfaches System auf Breadboard aufbauen mit 1MHz
- Wozmon implementieren
- ...

Am Ende des 1. Abschnittes soll dann eine CPU Steckkarte mit einem funktionierenden 6502 System stehen. Inkl. BASIC überTerminal. Als weitere Sprachen wären noch VTL-2 und Basl (meine eigene kleine Sprache) angedacht.

# Speichereinteilung

Der 6502 kann einen Adressbereich von 64Kb ansprechen. Dabei müssen folgende Bedingungen gegeben sein:

die beiden unteren Seiten (eine Seite sind immer 256 Bytes) sollten/müssen RAM sein. Die Zeropage (ZP) also der Bereich von 0x0000 -0x00FF hat eine besondere Bedeutung im 6502. Hier liegen "Register" die sehr schnell angesprochen werden können. Auch Bitmanipulation sind sehr einfach. 

Evtl. macht es somit Sinn ähnlich wie im 6510 (C64) die beiden unteren Bytes mit einem digitalen Ein/Ausgang zu besetzen. Damit könnte man sehr schnelle Anbindungen schreiben. z.B. SPI oder Taktumschaltung, oder auch ROM Selektion (dazu später mehr).

In der Page 1 liegt der "interne Stack" vom 6502. Auch dieser Bereich muss mit RAM versorgt sein, will man Stackbefehle oder auch nur Unterprogramme benutzen. 

Im obersten Bereich liegen beim 6502 die Start und Interrupt Vektoren. d.h. die Adressen 0xFFFA..0xFFFF sollten im ROM liegen. 

Somit ergibt sich zunächst folgende Einteilung

| 0x0000-0x00FF | 0x0100-0x01FF |      |      | ...  |      |      | 0xFF00-0xFFFF    |
| ------------- | ------------- | ---- | ---- | ---- | ---- | ---- | ---------------- |
| RAM ZP        | RAM Stack     |      |      |      |      |      | ROM für Vektoren |

Nun gibt es aber mittlerweile RAM und ROM im Überfluss, d.h. an kleine RAM oder ROM Bausteine zu kommen ist leider schwierig und Preislich auch nicht nötig. (JA, auch der Preis des Systemes ist immer wieder zu betrachten)

Da ich gerne mit BASIC oder auch mit anderen Sprachen herumhantieren möchte, ist es notwendig, dem Rechner eine ordentlich Portion RAM zu spendieren. Ich habe (ich weiß leider nicht mehr genau woher) noch ein paar RAM Bausteine da. Einen 62256 (32KBx8) und mehrere UM61512-15 (64Kbx8). Beides habe ich mit dem TL866 schon getestet und es funktioniert. Also kann ich die Bausteine weiter verwenden. Somit macht es durchaus Sinn das RAM einfach komplett unter den Adressbereich zu legen und dann die anderen Bereiche (IO und ROM) zu überlagern. Evtl. kann man sogar später mal ROM Bereiche ausblenden.

Auch das ROM gibt es in verschiedenen Varianten. Ich habe mir sowohl ein klassisches EEPROM mit 8Kx8 besorgt, wie auch ein NOR Flash mit 128KBx8. Das 2. ist dafür gedacht ROM Bereiche umschaltbar zu machen. D.h. ein Bereich Flash kann sowohl im Kernel ROM eingeblendet werden, ein andere dann als Interpreter ROM.

Zur Zeit stell ich mir die Aufteilung so vor:
| | |
|-|-|
| 0xFFFF<br />...<br />0xE000 | 8KB Kernel ROM, HiROM |
| 0xDFFF<br /><br />0xD000    | IO Bereich aufgeteilt in 16 Bereiche für die Peripherie.<br />0xD500: CS4<br />0xD400: CS3<br />0xD300: CS2<br />0xD200: CS1<br />0xD100: ASIC 1<br />0xD000: VIA 1<br /> |
| 0xCFFF<br />...<br />0xC000 | 4k RAM ( kennt man aus dem C64) |
| 0xBFFF<br />...<br />0xA000 | 8KB Interpreter ROM, LoROM |
| 0x9FFF<br />...<br />0x8000 | 4KB RAM |
| 0x7FFF<br />...<br />0x0200 | 31469 Bytes RAM (BaseRAM) |
| 0x01FF<br />...<br />0x0100 | 256 Bytes Stack |
| 0x00FF<br />...<br />0x0002 | 254 Bytes ZP RAM |
| 0x0001, 0x0000 | 16 Bit digitale Ein/Ausgabe |

# Der Bus

Der SBC soll später auf einer Backplane aufgesetzt werden können. Die Backplane übernimmt dann die Stromversorgung der verschiedenen Karten. Gerne hätte ich Edge Card Connectoren, wie sie auch im C64 Verwendung gefunden haben.  So braucht man für die Karten nur etwas Paltinenplatz und keinen eigenen Connector. 

![sch_clock_reset](/images/edge_card_connector.png)

## Aber welche Signale müssen auf den Bus?

Fangen wir mal mit dem Naheliegendsten an.

**+5V** und **GND** müssen da auf jeden Fall drauf. Am besten mehrfach, damit die Strombelastung pro Kontakt nicht zu hoch wird.
Dann der Adress- und der Datenbus. **A0..A15** und **D0..D7**.

Nun kommen wir zu dem schwierigsten Teil. Den Steuerleitungen. Klar sind natürlich

**IRQ, NMI, RESET, RW, PHI2, RDY**

Weiterhin kommen von der Adressdekodierlogik:

**CSIO2..6**, 4 zusätzliche ChipSelect Leitungen, die direkt aus dem Adressdekoder stammen, um weitere Peripherie anzuschließen. Jeder Adressbereich umfasst eine Page also 256 Adressen. (Wie auf der CPU Karte selber auch)

Dann gibt es jeweils Tupel von Anforderungs- und ChipSelectleitungen

**AHiROM, HiROM**: Mit der Anforderungsleitung AHiROM signalisiert die Karte, daß sie für den oberen Adressbereich (Die obersten 8K) ein eigenes ROM zur Verfügung stellt. (0xE000..0xFFFF)

**ALoROM, LoROM**: das geliche gilt für das LoROM, also den Bereich der Shell, terminal, BASIC Interpreters.(0xA000..0xBFFF)

**AHiRAM, HiRAM, ALoRAM, LoRAM**: hier nun für die HiRAM und LoRAM Bereiche. (HiRAM: 0xC000..0xCFFF, LoRAM: 0x8000..0x9FFF)

Werden alle 3 Bereiche (LoROM, HiRAM, LoRAM) genutzt, können 16KB zusammenhängender Adressbereich benutzt werden. Mit dem HiROM stehen also 24KB zur Verfügung.

| Pin    | Belegung | Pin    | Belegung |
| ------ | -------- | ------ | -------- |
| 1      | +5V      | 2      | +5V      |
| 3      | GND      | 4      | GND      |
| 5..11  | D0..D6   | 6..12  | D1..D7   |
| 13     | A0       | 14     | A1       |
| 15..25 | ...      | 16..26 | ...      |
| 27     | A14      | 28     | A15      |
| 29     | IRQB     | 30     | NMIB     |
| 31     | RESETB   | 32     | PHI2     |
| 33     | RW       | 34     | RDY      |
| 35     | AHiROMB  | 36     | HiROMB   |
| 37     | ALoROMB  | 38     | LoROMB   |
| 39     | AHiRAMB  | 40     | HiRAMB   |
| 41     | ALoRAMB  | 42     | LoRAMB   |
| 43     | CSB2     | 44     | CSB3     |
| 45     | CSB4     | 46     | CSB5     |
| 47     | n.n.     | 48     | n.n.     |
| 49     | GND      | 50     | GND      |

## Aber wofür?

### CS2..CS5

Hier kann man per Karte eigene zus. Peripherie anbinden. So ist es damit schnell möglich eine zus. ASIC oder eine VIA ohne viel Aufwand dazu zu stecken.

### HiROM, LoROM, HiRAM, LoRAM

Naja das ist ein moving Target. Ich bin mir noch nicht sicher, ob ich die Signale brauche.

Mit HiROM kann man ein eigenes Kernel einblenden. 

Mit LoROM eine eigene Shell mit eigener Sprache. So kann man Karten mit BASIC, VTL-2, Basl oder ähnliches verwenden.

LoRAM und HiRAM könnten für Grafikkarten Verwendung finden.
Alles zusammen gibt einem 24KB ROM (8KB (0xD000...0xDFFF sind den IO geschuldet)

