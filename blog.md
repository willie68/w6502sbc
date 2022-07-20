# Blog die erste

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

Am Ende des 1. Abschnittes soll dann eine CPU Steckkarte mit einem funktionierenden 6502 System stehen. Inkl. BASIC über Terminal. Als weitere Sprachen wären noch VTL-2 und Basl (meine eigene kleine Sprache) angedacht.

# Start the engine

Heute (19.07.22) kamen die letzten Teile. Jetzt kann ich loslegen. Also habe ich zunächst mit der Clocksection und dem Reset angefangen. Beides recht simple.

Für die Clock nehme ich im ersten Entwurf einfach einen entsprechenden DIP 14 Quarzoszillator. Ich habe 1, 2 und 4MHz zur Auswahl. Angefangen habe ich mit 1MHz. Auch der Reset ist schnell erledigt. Der DS1813 von Maxim erledigt den von ganz alleine. 

<img src="./images/clock_reset.jpg" alt="clock_reset" style="zoom:25%;" />

Schematic dazu

![sch_clock_reset](./images/sch_clock_reset.png)



# Busgedanken

Wie schon im 1. Teil erwähnt, soll das Projekt später auf einer Backplane aufsetzen. Die Backplane übernimmt dann die Stromversorgung der verschiedenen Karten. Gerne hätte ich Edge Card Connectoren, wie sie auch im C64 Verwendung gefunden haben.  So braucht man für die Karten nur etwas Paltinenplatz und keinen eigenen Connector. 

![sch_clock_reset](./images/edge_card_connector.png)

## Aber welche Signale müssen auf den Bus?

Fangen wir mal mit dem Naheliegendsten an.

**+5V** und **GND** müssen da auf jeden Fall drauf. Am besten mehrfach, damit die Strombelastung pro Kontakt nicht zu hoch wird.
Dann der Adress- und der Datenbus. **A0..A15** und **D0..D7**. Gerüchteküche: Ich würde gerne auch 24 Bit Adressen unterstützen, damit ich später mal eine W65C816 CPU Karte bauen kann... also A0..A23

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



Mal schauen wie sich das entwickelt. Bis Pin 34 steht die Aufteilung fest. Somit kann die Backplane schon mal entworfen werden, was für Signale dann da drauf kommen, ist erst mal egal.
