# w6502sbc
willies 6502 single board computer

a 6502 sbc with the option to be used with a backplane as a module.

I'm a bit sorry, but there are so many 6502 projects out there that I've taken the liberty of documenting this project (for the time being) only in German.

Hier dokumentiere ich den aktuellen Stand. Für die Historie gibt es den Blog. ([blog](/blog.md))

# Einleitung

Schon seit längerem beschäftigt mich ein altes Thema. Assembler und 6502. Mein C64 ist schon länger eingemottet, der VC20 verkauft. Trotzdem finde ich den 6502 von der Assemblersprache her sehr interessant. Was man da mit 1MHz und 64kb RAM alles machen konnte. Dieses Jahr (2022) war es dann soweit. Ich habe mich dazu durchgerungen einen 6502 SBC zu bauen. Im Internet gibt es viele, viele Beiträge und Vorschläge dazu. Leider gibt es keine Version, die man einfach nachbauen könnte und meine Erwartungen an einen 6502 Workbench Rechner erfüllt. Hier mal kurz meine Anforderungen:

- WDC65C02 CPU mit 4MHz (Takt per Programm umschaltbar)
- WDC65C22 VIA mit 4MHz, evtl. davon 2
- 6551 ACIA für die serielle Schnittstelle (R65C51) 
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

In der Page 1 liegt der "interne Stack" vom 6502. Auch dieser Bereich muss mit RAM versorgt sein, will man Stackbefehle oder auch nur Unterprogramme benutzen. 

Im obersten Bereich liegen beim 6502 die Start und Interrupt Vektoren. d.h. die Adressen 0xFFFA..0xFFFF sollten im ROM liegen. 

Somit ergibt sich zunächst folgende Einteilung

| 0x0000-0x00FF | 0x0100-0x01FF |      |      | ...  |      |      | 0xFF00-0xFFFF    |
| ------------- | ------------- | ---- | ---- | ---- | ---- | ---- | ---------------- |
| RAM ZP        | RAM Stack     |      |      |      |      |      | ROM für Vektoren |

Nun gibt es aber mittlerweile RAM und ROM im Überfluss, d.h. an kleine RAM oder ROM Bausteine zu kommen ist leider schwierig und Preislich auch nicht nötig. (Ja, auch der Preis des Systemes ist immer wieder zu betrachten)

Da ich gerne mit BASIC oder auch mit anderen Sprachen herumhantieren möchte, ist es notwendig, dem Rechner eine ordentlich Portion RAM zu spendieren. Ich habe (ich weiß leider nicht mehr genau woher) noch ein paar RAM Bausteine da. Einen 62256 (32KBx8) und mehrere UM61512-15 (64Kbx8). Beides habe ich mit dem TL866 schon getestet und es funktioniert. Also kann ich die Bausteine weiter verwenden. Somit macht es durchaus Sinn das RAM einfach komplett unter den Adressbereich zu legen und dann die anderen Bereiche (IO und ROM) zu überlagern. Evtl. kann man sogar später mal ROM Bereiche ausblenden.

Auch das ROM gibt es in verschiedenen Varianten. Ich habe mir sowohl ein klassisches EEPROM mit 8Kx8 besorgt, wie auch ein NOR Flash mit 128KBx8. Das 2. ist dafür gedacht ROM Bereiche umschaltbar zu machen. D.h. ein Bereich Flash kann sowohl im Kernel ROM eingeblendet werden, ein andere dann als Interpreter ROM.

Zur Zeit stell ich mir die Aufteilung so vor:

## simple Variante

| Bereich                     | Hi Address Nibble                       | Beschreibung                                                 |
| --------------------------- | --------------------------------------- | ------------------------------------------------------------ |
| 0xFFFF<br />...<br />0xC000 | 11xx xxxx                               | 16KB Kernel ROM, HiROM                                       |
| 0xBFFF<br /><br />0xB000    | 1011 xxxx                               | IO Bereich aufgeteilt in 16 Bereiche für die Peripherie.<br />0xD300: CS3 #0011<br />0xD200: CS2 #0010<br />0xD100: ASIC 1 #0001<br />0xD000: VIA 1 #0000<br />CS2 und 3 gehen nur auf den "Bus"<br /> |
| 0xAFFF<br />...<br />0x8000 | 1010 xxxx<br />1001 xxxx<br />1000 xxxx | 12k ext ROM, BASIC oder anderes ROM                          |
| 0x7FFF<br />...<br />0x0200 | 0xxx xxxx                               | 31469 Bytes RAM (BaseRAM)                                    |
| 0x01FF<br />...<br />0x0100 | 0xxx xxxx                               | 256 Bytes Stack                                              |
| 0x00FF<br />...<br />0x0000 | 0xxx xxxx                               | 256 Bytes ZP RAM                                             |

Für den Adressdecoder will ich einen ATF16V8B, also ein CPLD verwenden. Dadurch habe ich die Möglichkeit die Adressdekodierung variabel gestalten zu können. Da noch genügend Pins im CPLD vorhanden sind kann ich die unteren 4 IO Leitungen auch direkt dekodieren. Falls mehr gewünscht sind, kann man per CSIO einen weitern Dekodierer kaskadieren. (74HC138 o.ä.)

Die Verknüpfung von PHI2 und RAM ist bereits enthalten. 
Hier der Entwurf des PLDs: 

```wpld
header:
Name     adr_simple ;
PartNo   01 ;
Date     24.07.2022 ;
Revision 03 ;
Designer wkla ;
Company  nn ;
Assembly None ;
Location  ;
Device   G16V8 ;

pld:
/* *************** INPUT PINS *********************/
PIN [1..8]   =  [A15..A8]; 
PIN 9   =  PHI2;

/* *************** OUTPUT PINS *********************/
PIN 12   =  CSRAM;
PIN 13   =  CSHIROM;
PIN 14   =  CSEXTROM;
PIN 15   =  CSIO;
PIN 16   =  CSIO3;
PIN 17   =  CSIO2;
PIN 18   =  CSIO1;
PIN 19   =  CSIO0;
/* *************** LOGIC *********************/

FIELD Addr = [A15..A8];
CSRAM_EQU = Addr:[0000..7FFF]; // 32KB
IOPORT_EQU = Addr:[B000..BFFF]; // 4KB
VIAPORT_EQU = Addr:[B000..B0FF];
ACIAPORT_EQU = Addr:[B100..B1FF];
CSIO2PORT_EQU = Addr:[B200..B2FF];
CSIO3PORT_EQU = Addr:[B300..B3FF];
CSEXTROM_EQU = Addr:[8000..AFFF]; // 12KB
CSROM_EQU = Addr:[C000..FFFF];  // 16KB

/* ZP */
CSEXTROM = !CSEXTROM_EQU;

/* RAM */
CSRAM = !CSRAM_EQU # !PHI2;

/* 8kb of ROM */
CSHIROM = !CSROM_EQU;

/* IO */
CSIO= !IOPORT_EQU;
CSIO0 = !VIAPORT_EQU;
CSIO1 = !ACIAPORT_EQU;
CSIO2 = !CSIO2PORT_EQU;
CSIO3 = !CSIO3PORT_EQU;

simulator:
ORDER: A15, A14, A13, A12, A11, A10, A9, A8, PHI2, CSEXTROM, CSRAM, CSHIROM, CSIO, CSIO0, CSIO1, CSIO2, CSIO3; 

VECTORS:
/* internal RAM */
0 X X X X X X X 0 H H H H H H H H 
0 X X X X X X X 1 H L H H H H H H 

/* 8000-AFFF external Rom */ 
1 0 0 0 X X X X X L H H H H H H H 
1 0 0 1 X X X X X L H H H H H H H 
1 0 1 0 X X X X X L H H H H H H H 

/* IO */ 
/* CSIO0 */
1 0 1 1 0 0 0 0 X H H H L L H H H 
/* CSIO1 */
1 0 1 1 0 0 0 1 X H H H L H L H H 
/* CSIO2 */
1 0 1 1 0 0 1 0 X H H H L H H L H 
/* CSIO3 */
1 0 1 1 0 0 1 1 X H H H L H H H L 
/* nicht direkt benutzt */
1 0 1 1 0 1 X X X H H H L H H H H 
1 0 1 1 1 X X X X H H H L H H H H 
/* ROM */
1 1 X X X X X X X H H L H H H H H 

```

Das Format ist wpld. Eine Erweiterung von pld von mir. Mit diesem kleinen Tool kann ich PLD und SI File in einer Datei bearbeiten. Das Tool generiert automatisch die erforderlichen Dateien (pld und si) aus dieser Quelle in einem eigenen Unterverzeichnis und startet dort dann CUPL.

## C64 like

Evtl. macht es somit Sinn ähnlich wie im 6510 (C64) die beiden unteren Bytes mit einem digitalen Ein/Ausgang zu besetzen. Damit könnte man sehr schnelle Anbindungen schreiben. z.B. SPI oder Taktumschaltung, oder auch ROM Selektion (dazu später mehr).

Eine C64 ähnliche Aufteilung wäre das hier:
| Bereich | Hi Adress | Beschreibung |
|-|-|-|
| 0xFFFF<br />...<br />0xE000 | 111x xxxx | 8KB Kernel ROM, HiROM |
| 0xDFFF<br /><br />0xD000    | 1101 xxxx | IO Bereich aufgeteilt in 16 Bereiche für die Peripherie.<br />0xD500: CS5<br />0xD400: CS4<br />0xD300: CS3<br />0xD200: CS2<br />0xD100: ASIC 1<br />0xD000: VIA 1<br /> |
| 0xCFFF<br />...<br />0xC000 | 1100 xxxx | 4k RAM (HiRAM, kennt man aus dem C64) |
| 0xBFFF<br />...<br />0xA000 | 101x xxxx | 8KB Interpreter ROM, LoROM<br />Dieses ROM ist per NoLoROM abschaltbar und wird durch RAM ersetzt. Somit erhält man 51KB durchgängigen RAM Speicher. |
| 0x9FFF<br />...<br />0x8000 | 100x xxxx | 8KB RAM (LoRAM) |
| 0x7FFF<br />...<br />0x0200 | 0xxx xxxx | 31469 Bytes RAM (BaseRAM) |
| 0x01FF<br />...<br />0x0100 | 0xxx xxxx | 256 Bytes Stack |
| 0x00FF<br />...<br />0x0002 | 0xxx xxxx | 254 Bytes ZP RAM |
| 0x0001, 0x0000 | 0xxx xxxx | 16 Bit digitale Ein/Ausgabe (optional) |

Die Adressdekodierung erfolgt nun auf Basis dieser Tabelle. 

| **Bereich**    | **A15** | **A14** | **A13** | **A12** |      | **/NOLOROM** | **/AloRAM** | **/AhiRAM** | **/AloROM** | **/AhiROM** |      | **/CSRAM** | **/CSHiROM** | **/CSLoROM** | **/CSIO** | **/LoRAM** | **/LoROM** | **/HiRAM** | **/HiROM** |                  |
| -------------- | ------- | ------- | ------- | ------- | ---- | ------------ | ----------- | ----------- | ----------- | ----------- | ---- | ---------- | ------------ | ------------ | --------- | ---------- | ---------- | ---------- | ---------- | ---------------- |
| **Base RAM**   | 0       | x       | x       | x       |      | x            | x           | x           | x           | x           |      | 0          | 1            | 1            | 1         | 1          | 1          | 1          | 1          | Lower 32K of RAM |
| **Lo RAM**     | 1       | 0       | 0       | x       |      | x            | 1           | x           | x           | x           |      | 0          | 1            | 1            | 1         | 1          | 1          | 1          | 1          | internal Lo RAM  |
| **Lo RAM ext** | 1       | 0       | 0       | x       |      | x            | 0           | x           | x           | x           |      | 1          | 1            | 1            | 1         | 0          | 1          | 1          | 1          | external Lo RAM  |
| **No Lo ROM**  | 1       | 0       | 1       | x       |      | 0            | x           | x           | 1           | x           |      | 0          | 1            | 1            | 1         | 1          | 1          | 1          | 1          | no Lo ROM        |
| **Lo ROM**     | 1       | 0       | 1       | x       |      | 1            | x           | x           | 1           | x           |      | 1          | 1            | 0            | 1         | 1          | 1          | 1          | 1          | internal Lo ROM  |
| **Lo ROM ext** | 1       | 0       | 1       | x       |      | 1            | x           | x           | 0           | x           |      | 1          | 1            | 1            | 1         | 1          | 0          | 1          | 1          | external Lo ROM  |
| **Hi RAM**     | 1       | 1       | 0       | 0       |      | x            | x           | 1           | x           | x           |      | 0          | 1            | 1            | 1         | 1          | 1          | 1          | 1          | internal Hi RAM  |
| **Hi RAM ext** | 1       | 1       | 0       | 0       |      | x            | x           | 0           | x           | x           |      | 1          | 1            | 1            | 1         | 1          | 1          | 0          | 1          | external Hi RAM  |
| **IO**         | 1       | 1       | 0       | 1       |      | x            | x           | x           | x           | x           |      | 1          | 1            | 1            | 0         | 1          | 1          | 1          | 1          | IO               |
| **Hi ROM**     | 1       | 1       | 1       | x       |      | x            | x           | x           | x           | 1           |      | 1          | 0            | 1            | 1         | 1          | 1          | 1          | 1          | internal Hi ROM  |
| **Hi ROM ext** | 1       | 1       | 1       | x       |      | x            | x           | x           | x           | 0           |      | 1          | 1            | 1            | 1         | 1          | 1          | 1          | 0          | external Hi ROM  |

Für den Adressdecoder will ich einen ATF16V8B, also ein CPLD verwenden. Dadurch habe ich die Möglichkeit die Adressdekodierung variabel gestalten zu können. Für den IO Bereich setze ich einen zus. 74HC138 ein, der die unteren Pages aufteilt.

Die Verknüpfung von PHI2 und RAM ist bereits enthalten. 
Hier der Entwurf des PLDs: 

```pld
Name     W6502SBC_ADR ;
PartNo   01 ;
Date     20.07.2022 ;
Revision 01 ;
Designer wkla ;
Company  nn ;
Assembly None ;
Location  ;
Device   G16V8 ;

/* *************** INPUT PINS *********************/
PIN 1   =  A12; 
PIN 2   =  A13;
PIN 3   =  A14;
PIN 4   =  A15;
PIN 5 	=  ALORAM;
PIN 6 	=  AHIRAM;
PIN 7 	=  ALOROM;
PIN 8 	=  AHIROM;
PIN 9   =  PHI2;
PIN 11  =  NOLOROM;

/* *************** OUTPUT PINS *********************/
PIN 12   =  CSRAM;
PIN 13   =  CSHIROM;
PIN 14   =  CSLOROM;
PIN 15   =  CSIO;
PIN 16   =  LORAM;
PIN 17   =  LOROM;
PIN 18   =  HIRAM;
PIN 19   =  HIROM;

 
CSRAM = (A15 & !A14 & !A13 & !ALORAM) # (A15 & !A14 & A13 & NOLOROM) # (A15 & A14 & !A13 & !A12 & !AHIRAM) # (A15 & A14 & !A13 & A12) # (A15 & A14 & A13) # !PHI2;
CSHIROM = !(A15 & A14 & A13 & AHIROM);
CSLOROM = !(A15 & !A14 & A13 & ALOROM & NOLOROM);
CSIO= !(A15 & A14 & !A13 & A12);
LORAM= !(A15 & !A14 & !A13 & !ALORAM);
LOROM= !(A15 & !A14 & A13 & !ALOROM & NOLOROM);
HIRAM= !(A15 & A14 & !A13 & !A12 & !AHIRAM);
HIROM= !(A15 & A14 & A13 & !AHIROM);
```

Und der Simulator dazu:

```
Name     W6502SBC_ADR ;
PartNo   01 ;
Date     20.07.2022 ;
Revision 01 ;
Designer wkla ;
Company  nn ;
Assembly None ;
Location  ;
Device   G16V8 ;

ORDER: A15, A14, A13, A12, ALORAM, AHIRAM, ALOROM, AHIROM, NOLOROM, PHI2, CSRAM, CSHIROM, CSLOROM, CSIO, LORAM, LOROM, HIRAM, HIROM; 

VECTORS:
0 X X X X X X X X 0 H H H H H H H H 
0 X X X X X X X X 1 L H H H H H H H 
1 0 0 X 1 X X X X 0 H H H H H H H H 
1 0 0 X 1 X X X X 1 L H H H H H H H 
1 0 0 X 0 X X X X X H H H H L H H H 
1 0 1 X X X 1 X 1 X H H L H H H H H 
1 0 1 X X X 0 X 1 X H H H H H L H H 
1 0 1 X X X X X 0 0 H H H H H H H H 
1 0 1 X X X X X 0 1 L H H H H H H H 
1 1 0 0 X 1 X X X 0 H H H H H H H H 
1 1 0 0 X 1 X X X 1 L H H H H H H H 
1 1 0 0 X 0 X X X X H H H H H H L H 
1 1 0 1 X X X X X X H H H L H H H H 
1 1 1 X X X X 1 X X H L H H H H H H 
1 1 1 X X X X 0 X X H H H H H H H L 

```



# Der Bus

## einfache Version

Für die 1. Version muss es zunächst auch ein einfacher Bus tun. Ich habe mich für eine einfache 40-pol Stiftleiste entschieden.

Darauf gibt es dann folgende Signale:

| Pin    | Bez     | Pin    | Bez       |
| ------ | ------- | ------ | --------- |
| 1      | +5V     | 2      | +5V       |
| 3      | /IRQ    | 4      | /RES      |
| 5      | /NMI    | 6      | CLK       |
| 7      | A0      | 8      | BE        |
| 9      | A1      | 10     | /CSHiROM  |
| 11     | A2      | 12     | /CSExtROM |
| 13     | A3      | 14     | /CSIO     |
| 15     | A4      | 16     | /CSIO2    |
| 17     | A5      | 18     | /CSIO3    |
| 19     | A6      | 20     | RDY       |
| 21     | A7      | 22     | R/W       |
| 23..37 | A8..A15 | 24..38 | D0..D7    |
| 39     | GND     | 40     | GND       |

Mit dem Bus kann man schnell mal das HiROM als Karte (z.B. mit ZIF Sockel oder mit ISP Möglichkeit) aufsetzen. Man kann später dann auch mal schnell ein ExtROM zufügen, oder noch eine VIA oder ACIA. Ein echter BUS ist das aber nicht.

## spätere Version

Der SBC soll später auf einer Backplane aufgesetzt werden können. Die Backplane übernimmt dann die Stromversorgung der verschiedenen Karten. Gerne hätte ich Edge Card Connectoren, wie sie auch im C64 Verwendung gefunden haben.  So braucht man für die Karten nur etwas Platinenplatz und keinen eigenen Connector. 

![sch_clock_reset](/images/edge_card_connector.png)

### Aber welche Signale müssen auf den Bus?

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
| 1      | GND      | 2      | GND      |
| 3      | +5V      | 4      | +5V      |
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
| 43     | CS2B     | 44     | CS3B     |
| 45     | CS4B     | 46     | CS5B     |
| 47     | n.n.     | 48     | n.n.     |
| 49     | GND      | 50     | GND      |

### Aber wofür?

#### CS2..CS5

Hier kann man per Karte eigene zus. Peripherie anbinden. So ist es damit schnell möglich eine zus. ASIC oder eine VIA ohne viel Aufwand dazu zu stecken.

#### HiROM, LoROM, HiRAM, LoRAM

Naja das ist ein moving Target. Ich bin mir noch nicht sicher, ob ich die Signale brauche.

Mit HiROM kann man ein eigenes Kernel einblenden. 

Mit LoROM eine eigene Shell mit eigener Sprache. So kann man Karten mit BASIC, VTL-2, Basl oder ähnliches verwenden.

LoRAM und HiRAM könnten für Grafikkarten Verwendung finden.
Alles zusammen gibt einem 24KB ROM (8KB (0xD000...0xDFFF sind den IO geschuldet)

Eine Speicherkarte kann nun aber per NoLoROM Signal auch das LoROM ausschalten und erhält dadurch zus. RAM. Es ergibt sich damit ein 51KB grosser RAM Bereich (0x0000 .. 0xCFFF). 
Warum liegt das Signal nicht auf dem Bus? 
Auf dem Bus macht das Signal keinen Sinn, es wird durch einen Port geschaltet. Dadurch kann man programatisch das ROM aus und einblenden. Für eine BusKarte macht das nur Sinn, wenn es das HiROM besetzt. Somit kann man dort auch die ROM Abschaltung implementieren. 

# Bauteile

Die verschiedenen Bauteile zu bekommen war eine kleine Herausforderung.

1. Ukrainekrieg, dadurch musste ich einiges an Schriftverkehr leisten, um überhaupt an CPU und Peripherie zu kommen
2. kleine 8KB EEPROMs sind sehr schwer zu bekommen. Deswegen habe ich auch ein 1MBit großes NOR Flash für die CPU Karte vorgesehen. Ich habe dabei folgende Adresszuordnungen vorgenommen.
   A16..A14 gehen auf ein Jumperfeld mit 3 Jumpern, A13 -> A14, und A12..A0 gehen auf die entsprechenden Adressleitungen. Auf der CPU liegen LoROM und HiROM 8kb auseinander. durch diese Zuordnung ist es mir nun möglich, im Flash HiROM und LoROM jeweils direkt hintereinander zu legen und dann in einem Flash 8 verschiedene Versionen (selektierbar durch Jumper) zu programmieren.

## BOM

WDC65C02 CPU 

WDC65C22 VIA

R65C51 ACIA

UM62256 32KB SRAM

SST39F010-70 1MBit NOR Flash (128x8)

ATF16V8 CPLD

74HC138 3-8 Binärdecoder

1MHZ Quarzoszilator

DS1813 EconoReset

div Widerständer, Switches, Jumper...

