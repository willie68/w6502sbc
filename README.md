# w6502sbc
willies 6502 single board computer

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

die beiden unteren Seiten (eine Seite sind immer 256 Bytes) sollten/müssen RAM sein. Die Zeropage (ZP) also der BEreich von 0x0000 -0x00FF hat eine besondere Bedeutung im 6502. Hier liegen "Register" die sehr schnell angesprochen werden können. Auch Bitmanipulation sind sehr einfach. 

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
| 0xFFFF<br />...<br />0xE000 | 8KB Kernel ROM |
| 0xDFFF<br /><br />0xD000    | IO Bereich aufgeteilt in 16 Bereiche für die Peripherie |
| 0xCFFF<br />...<br />0xC000 | 4k RAM ( kennt man aus dem C64) |
| 0xBFFF<br />...<br />0xA000 | 8KB Interpreter ROM |
| 0x9FFF<br />...<br />0x8000 | 4KB RAM |
| 0x7FFF<br />...<br />0x0200 | 31469 Bytes RAM |
| 0x01FF<br />...<br />0x0100 | 256 Bytes Stack |
| 0x00FF<br />...<br />0x0002 | 254 Bytes ZP RAM |
| 0x0001, 0x0000 | 16 Bit digitale Ein/Ausgabe |
