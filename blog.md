# Blog die erste (18.07.22)

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

Am Ende des 1. Abschnittes soll dann eine CPU Steckkarte mit einem funktionierenden 6502 System stehen. Inkl. BASIC über Terminal. Als weitere Sprachen wären noch VTL-2 und Basl (meine eigene kleine Sprache) angedacht.

# Start the engine (19.07.22)

Heute (19.07.22) kamen die letzten Teile. Jetzt kann ich loslegen. Also habe ich zunächst mit der Clocksection und dem Reset angefangen. Beides recht simple.

Für die Clock nehme ich im ersten Entwurf einfach einen entsprechenden DIP 14 Quarzoszillator. Ich habe 1, 2 und 4MHz zur Auswahl. Angefangen habe ich mit 1MHz. Auch der Reset ist schnell erledigt. Der DS1813 von Maxim erledigt den von ganz alleine. 

<img src="./images/clock_reset.jpg" alt="clock_reset" style="zoom:25%;" />

Schematic dazu

![sch_clock_reset](./images/sch_clock_reset.png)

## BOM

Quarzoszillator 1MHz (Oder auch mehr)

Maxim DS1813 EconoReset

PushSwitch

# Busgedanken (19.07.22)

Wie schon im 1. Teil erwähnt, soll das Projekt später auf einer Backplane aufsetzen. Die Backplane übernimmt dann die Stromversorgung der verschiedenen Karten. Gerne hätte ich Edge Card Connectoren, wie sie auch im C64 Verwendung gefunden haben.  So braucht man für die Karten nur etwas Platinenplatz und keinen eigenen Connector. 

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

# Adressen (20.07.22)

Heute hab ich mir mal ein paar Gedanken über die Adressierung gemacht. Eine erste Version der geplanten Speicheraufteilung steht ja schon in der [README.md](/readme.md) . Dazu muss aber auch ein PLD File für das CPLD geschrieben werden. Zunächst habe ich die verschiedenen Signal in einer Tabelle mal zusammengefasst.

| **Bereich**    | **A15** | **A14** | **A13** | **A12** |      | **/AloRAM** | **/AhiRAM** | **/AloROM** | **/AhiROM** |      | **/CSRAM** | **/CSHiROM** | **/CSLoROM** | **/CSIO** | **/LoRAM** | **/LoROM** | **/HiRAM** | **/HiROM** |
| -------------- | ------- | ------- | ------- | ------- | ---- | ----------- | ----------- | ----------- | ----------- | ---- | ---------- | ------------ | ------------ | --------- | ---------- | ---------- | ---------- | ---------- |
| **Base RAM**   | 0       | x       | x       | x       |      | x           | x           | x           | x           |      | 0          | 1            | 1            | 1         | 1          | 1          | 1          | 1          |
| **Lo RAM**     | 1       | 0       | 0       | x       |      | 1           | x           | x           | x           |      | 0          | 1            | 1            | 1         | 1          | 1          | 1          | 1          |
| **Lo RAM ext** | 1       | 0       | 0       | x       |      | 0           | x           | x           | x           |      | 1          | 1            | 1            | 1         | 0          | 1          | 1          | 1          |
| **Lo ROM**     | 1       | 0       | 1       | x       |      | x           | x           | 1           | x           |      | 1          | 1            | 0            | 1         | 1          | 1          | 1          | 1          |
| **Lo ROM ext** | 1       | 0       | 1       | x       |      | x           | x           | 0           | x           |      | 1          | 1            | 1            | 1         | 1          | 0          | 1          | 1          |
| **Hi RAM**     | 1       | 1       | 0       | 0       |      | x           | 1           | x           | x           |      | 0          | 1            | 1            | 1         | 1          | 1          | 1          | 1          |
| **Hi RAM ext** | 1       | 1       | 0       | 0       |      | x           | 0           | x           | x           |      | 1          | 1            | 1            | 1         | 1          | 1          | 0          | 1          |
| **IO**         | 1       | 1       | 0       | 1       |      | x           | x           | x           | x           |      | 1          | 1            | 1            | 0         | 1          | 1          | 1          | 1          |
| **Hi ROM**     | 1       | 1       | 1       | x       |      | x           | x           | x           | 1           |      | 1          | 0            | 1            | 1         | 1          | 1          | 1          | 1          |
| **Hi ROM ext** | 1       | 1       | 1       | x       |      | x           | x           | x           | 0           |      | 1          | 1            | 1            | 1         | 1          | 1          | 1          | 0          |

Der IO Bereich wird mit einem 74138 noch einmal extra unterteilt und braucht somit nicht mit ins CPLD. Geplant ist auch nur ein kleines ATF16V8. Wie man hier sieht, habe ich 8 Eingänge und 8 Ausgänge. Für das RAM wird die notwendige Kombination mit dem PHI2 extern gemacht. Evtl. zieh ich das aber auch hier noch mit rein. Steht noch nicht fest.

Herausgekommen sind die beiden folgenden Dateien: PLD

```pld
Name     W6502SBC_ADR ;
PartNo   00 ;
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

/* *************** OUTPUT PINS *********************/
PIN 12   =  CSRAM;
PIN 13   =  CSHIROM;
PIN 14   =  CSLOROM;
PIN 15   =  CSIO;
PIN 16   =  LORAM;
PIN 17   =  LOROM;
PIN 18   =  HIRAM;
PIN 19   =  HIROM;

CSRAM = (A15 & !A14 & !A13 & !ALORAM) # (A15 & !A14 & A13) # (A15 & A14 & !A13 & !A12 & !AHIRAM) # (A15 & A14 & !A13 & A12) # (A15 & A14 & A13) ;
CSHIROM = !(A15 & A14 & A13 & AHIROM);
CSLOROM = !(A15 & !A14 & A13 & ALOROM);
CSIO= !(A15 & A14 & !A13 & A12);
LORAM= !(A15 & !A14 & !A13 & !ALORAM);
LOROM= !(A15 & !A14 & A13 & !ALOROM);
HIRAM= !(A15 & A14 & !A13 & !A12 & !AHIRAM);
HIROM= !(A15 & A14 & A13 & !AHIROM);
```

Und der Simulator dazu:

```
Name     W6502SBC_ADR ;
PartNo   00 ;
Date     20.07.2022 ;
Revision 01 ;
Designer wkla ;
Company  nn ;
Assembly None ;
Location  ;
Device   G16V8 ;

ORDER: A15, A14, A13, A12, ALORAM, AHIRAM, ALOROM, AHIROM, CSRAM, CSHIROM, CSLOROM, CSIO, LORAM, LOROM, HIRAM, HIROM; 

VECTORS:
0 X X X X X X X L H H H H H H H 
1 0 0 X 1 X X X L H H H H H H H 
1 0 0 X 0 X X X H H H H L H H H 
1 0 1 X X X 1 X H H L H H H H H 
1 0 1 X X X 0 X H H H H H L H H 
1 1 0 0 X 1 X X L H H H H H H H 
1 1 0 0 X 0 X X H H H H H H L H 
1 1 0 1 X X X X H H H L H H H H 
1 1 1 X X X X 1 H L H H H H H H 
1 1 1 X X X X 0 H H H H H H H L 
```

Wer einen Fehler findet mag ihn mir gerne mitteilen. 

# Der NOP Generator

Der erste Schritt mit der CPU für mich ist ein sog. NOP Generator. Wenn der 6502 startet, ließt er zunächst aus den Adressen 0xFFFC/0xFFFD die Adresse, wo er seine Startroutine finden soll. Dort springt der Prozessor dann hin und führt den Code aus. Wenn man nun den Datenbus auf das NOP (no operation)  Kommando $EA fest verdrahtet, passiert folgendes: Zunächst ließt die CPU nach dem Einschalten (oder auch nach einem Reset) den Reset Vektor $EAEA. Damit wird nun der Adresszeiger geladen. Als nächstes ließt die CPU den ersten Befehl, auch wieder ein $EA und führt diesen aus. dabei wird der Adresszeiger inkrementiert. Da nichts zu tun ist, NOP, liesst er von der neuen Adresse den nächsten, erhöht wieder den Adresszeiger und führt diesen aus. Und so weiter... Wie man sieht erhöht sich die Adresse bis zum Überlauf des Adresszeiger. Dort wird dann einfach bei $0000 weiter gemacht und das ganze wiederholt sich. Wenn man nun an den Adressbus (vor allen an den höheren Bits) LED anschließt, sollte man einen typischen Zähler sehen. Wenn das funktioniert, funktioniert sowohl der Takt wie auch die CPU selber.

Hier mal ein Plan dazu:

![nop_generator](./images/nop_generator.png)

Hier mal ein Blick auf das Steckbrett.

![nop_board](./images/nop_board.jpg)

Die Widerstände bilden das Wort $EA ab. D0 befindet sich links, D7 rechts.
Der kleine schwarze Kasten ist mein 8 Kanal Logik Analyser. Angeschlossen habe ich R/W, A0, A10-A15.
Hier mal ein Bild von dem "Programm"

![nop_logic](./images/nop_logic.png)

Man sieht hier sehr schön wie die Adressleitungen hoch gezählt werden. Wir können auch mal hineinzoomen um die Frequenz zu messen. Am besten der A0 Leitung. Jeder Wechsel ist eine neue Instruktion. (Kanalbelegung geändert R/W, Clk, A0, A11-A15)

![nop_logic](./images/nop_logic_detail.png)



Die Adressleitung wird mit 250KHz umgeschaltet, d.h. die NOP Kommandos werden mit 500KHz verarbeitet, das bedeutet ein NOP braucht 2 Takte.

![nop_logic](./images/nop_logic_detail_2.png)

Also funktioniert mein Aufbau erst einmal.

# EEPROM und mehr?

Nun geht's richtig los. Jetzt wird ein EEPROM (AT28C64 8Kx8) angeschlossen. Und da das RAM ja eh auch bald dran ist, hab ich das gleich schon mal mit verkabelt. Zunächst aber nur Adress- und Datenleitung. Die Steuerleitungen vom RAM sind alle auf +5V, sodass das eigentlich nicht aktiv sein sollte. Beim EEPROM sind die Steuerleitung so gesetzt, dass es immer mit einem READ reagiert. Allerdings so richtig tun kann der SBC ja noch nix. Keine Ausgabe, keine Eingabe, da muss ich mir noch was überlegen. 

![nop_logic](./images/rom_ram.jpg)

Vorher:

![](./images/hookup_start.jpg)

Nachher:

![nop_logic](./images/hookup_rom_ram.jpg)

Hmm, derzeit wird also alles aus dem ROM geladen. Egal wo der Prozessor hin greift, es antwortet immer das ROM. Die 8KB werden im kompletten Speicher dupliziert. Also 8x. Wenn ich nun das ROM mit $EA (NOP) fülle, sollte sich eigentlich der gleiche Effekt wie beim NOP Generator ergeben. Also hab ich das WWPROM in den Prommer gepackt, mit $EA gefüllt und wieder eingebaut.

Klar das hat nicht auf Anhieb funktioniert. Wer findet den Fehler? Der Adresszähler machte komische Dinge. Erst als ich das RAM heraus genommen hatte, gings. Also das RAM mit dem Progger getestet. Der sagt: "alles OK". Noch mal auf die Schaltung geschaut, da sah ich es... Ich hab den GND auf die falsche Rail (+5V) gesteckt. Also RAM wieder eingebaut, GND richtig verbunden und schon gings. 

![video_nop_ram_rom](./images/nop_ram_rom.mp4)

# Die Logik

Um RAM und ROM und später auch die Interfacebausteine unter den richtigen Adressen ansprechen zu können, ist etwas Logik nötig. Die großen Speicherblöcke für RAM und ROM will ich über den CPLD ansprechen. Das ist schnell und man kann die Logik später modifizieren ohne die Schaltung zu ändern. Den für den IO Bereich hab ich den Bereich von $D000 bis $DFFF also 4Kb vorgesehen. Diese will ich in 16 Bereiche aufteilen. Jeder Bereich hat dann Platz für 256 Register (A0-A7). Auf dem SBC verwende ich dazu einen 74HC138. Dieser wird einerseits mit dem CSIO Signal aus dem CPLD selektiert. Als zweite Selektion kommt A11 zum Einsatz. Dekodiert werden dann A10..A8. Somit habe ich die untere Hälfte der 16 Bereiche bereits hier dekodiert. Der CPLD bekommt zur Dekodierung zunächst die Leitungen A12..A15. Damit kann man die jeweiligen 4K Blöcke aufteilen. Für das RAM muss man zusätzlich auch noch die CLK (PHI2) Leitung verwenden, damit die Schreibvorgänge richtig synchronisiert werden. (https://wilsonminesco.com/6502primer/addr_decoding.html) Im Projekt habe ich ja bereits die große Adress-CLPD Logik geschrieben. Für das Steckbrett verwende ich eine kleinere Variante. Hier wird erst einmal nur der der CPLD zur Adressdekodierung verwendet. Da sich die eigentliche Adressaufteilung nicht ändern soll muss ich auf den CPLD die Adressen A15..A8 legen. Für das RAM brauche ich zusätzlich noch das CLK 8PHI2) Signal. Hier mal die PLD Datei im meinem eigenen WPLD Format. Dazu später mehr.

adr_simple.wpld

```WCPLD
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
PIN 15   =  CSIO;
PIN 16   =  CSIO3;
PIN 17   =  CSIO2;
PIN 18   =  CSIO1;
PIN 19   =  CSIO0;
/* *************** LOGIC *********************/

FIELD Addr = [A15..A8];
IOPORT_EQU = Addr:[D000..DFFF];
VIAPORT_EQU = Addr:[D000..D0FF];
ACIAPORT_EQU = Addr:[D100..D1FF];
CSIO2PORT_EQU = Addr:[D200..D2FF];
CSIO3PORT_EQU = Addr:[D300..D3FF];
CSRAM_EQU = Addr:[0000..7FFF];
CSROM_EQU = Addr:[E000..FFFF];

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
ORDER: A15, A14, A13, A12, A11, A10, A9, A8, PHI2, CSRAM, CSHIROM, CSIO, CSIO0, CSIO1, CSIO2, CSIO3; 

VECTORS:
/* RAM */
0 X X X X X X X 0 H H H H H H H 
0 X X X X X X X 1 L H H H H H H 
/* 8000-CFFF nothing */ 
1 0 X X X X X X X H H H H H H H 
1 1 0 0 X X X X X H H H H H H H 
1 1 0 0 X X X X X H H H H H H H 
/* IO */ 
/* CSIO0 */
1 1 0 1 0 0 0 0 X H H L L H H H 
/* CSIO1 */
1 1 0 1 0 0 0 1 X H H L H L H H 
/* CSIO2 */
1 1 0 1 0 0 1 0 X H H L H H L H 
/* CSIO3 */
1 1 0 1 0 0 1 1 X H H L H H H L 
/* nicht direkt benutzt */
1 1 0 1 0 1 X X X H H L H H H H 
1 1 0 1 1 X X X X H H L H H H H 
/* ROM */
1 1 1 X X X X X X H L H H H H H 


```

Und hier mal gleich der Schaltplan dazu:

![](./images/adr_logic_simple.png)

## WPLD Format und das Tool WCUPL

Da WinCUPL auf keinem meiner Rechner mehr läuft, habe ich mir eine Batch geschrieben, mit denen ich sowohl das JED File erzeuge, wie auch gleich einen Test durchführen kann. Ich persönlich mag es allerdings, wenn ich die Quellen an einer Stelle habe. Die Dopplung der Header von der PLD und SI-Datei  und damit die manuelle syncronisation zwischen den beiden finde ich persönlich unschön. Deswegen hab ich mir das einfache WPLD Format ausgedacht und dazu das WCUPL Tool geschrieben. Das Tool WCUPL gibt's in meinem Repo (https://github.com/willie68/w6502sbc/releases). Dadurch entfällt das leidige Header abgleichen und Logic und Test liegen in einer Datei. Das Tool macht nicht viel. Evtl. Argumente werden direkt an CUPL durchgereicht. Aus der wpld Dateien werden `header:` und `pld:` zu einer #.pld Datei gemerged und `header`: und `simulator:` Teil zu der #.si Datei. Dann wird cupl aus dem gleichen Verzeichniss wo auch wcupl liegt gestartet. 

## Test

Bevor ich das ganze zusammen baue, teste ich mal die PLD Logik. Man glaubt es kaum, der 1. Test, kompilieren, brennen mit TL866, und auf dem Steckbrett aufbauen, war direkt erfolgreich. Sie tut was sie soll. Ich hab mal ein paar Kombinationen durchgespielt.

<img src="./images/adr_logic_simple_hookup.jpg" style="zoom:33%;" />

# Kommunikation ist alles

Nun kommt zunächst der 6522 mit dran. Damit ich endlich auch mal ein kleines Programm in den Rechner laden kann und dann auch sehe was passiert. Die Register des 6522 sollen an den VIA Port zu finden sein, also schließe ich den CS2B and den CSIO0B an. CS1 liegt direkt auf +5V. 

![](./images/via_schematic.png)

![hookup_6522_via](./images/hookup_6522_via.jpg)





# Terminal

Da mein SBC ja weder über Tastatur noch Monitor verfügt muss eine andere Lösung her. Dazu verwende ich eine serielle Schnittstelle. Als Baustein verwende ich einen R65C51. Der entsprechende W65C51 ist leider derzeit nicht zu bekommen. Weiterhin hört man von größeren Problemen mit dem Chip. Um damit ein Terminal auf dem PC anzusprechen muss zusätzlich auch noch ein Seriell/USB Adapter her. Die gibt es zum Glück mannigfach. 

# Software

Auf dem SBC soll verschiedene Software zum Einsatz kommen. Als einfachstes will ich zunächst Steve Wozniak WOZMON integrieren. Danach eine VTL-2 Variante. Und dann hätte ich gerne ein Basic. und und und...

- Conways "game of live" in einer Terminalversion...
