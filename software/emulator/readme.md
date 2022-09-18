# W6502SBC Emulator

## Grundsätzliches

Dies ist der W6502SBC Emulator. Er dient dazu, meinen 6502 SBC zu emulieren, um Problemen mit dem Bios und der Peripherie auf den Grund zu gehen. Dieser Emulator ist speziell für meinen SBC zugeschnitten, weswegen diverse Dinge anders funktionieren als bei anderen Emulatoren. 

- Speziell auf das Adressmodell meines SBC zugeschnitten
- nicht alle Mnemoniks bzw. Operationmodes sind implementiert. (Beispielsweise fehlt derzeit der BCD Modus komplett)
- einfaches Sprungtabellen design, statt Nachbau der Gruppenstruktur zur Evaluierung des Adressmodes eines Befehls. Der Emulator hat eine einfache Sprungtabelle, wo zu jedem Commandbyte die entsprechende Funktion hinterlegt ist. D.h. es gibt keine Gruppenkommandos. Beispiel: LDA gibt es in 8 verschiedenen Adressmodi. Gemeinsam ist das ein Wert in den Akkumulator gelegt wird. Woher dieser Wert stammt, wird über die Adressmodi bestimmt. In vielen Emulatoren gibt es somit eine Funktion für den LDA Befehl, die dann zunächst den Adressmode bestimmt und diesen dann auswertet und die entsprechenden Operationen durchführt. In meinem Emulator gibt es für jeden LDA Befehl eine eigenen Funktion. 
  Hintergrund: Somit können auch andere Prozessoren wie der 65C02 oder auch ein 65C816 einfacher emuliert werden, da für diese Prozessoren die neuen Anweisungen teilweise nicht an entsprechenden Stellen im Commandoraum gelegt wurden. 
  Auch ist es für den Programmierer einfacher, den Code des Befehls nachzuvollziehen.
- Bausteine, wie z.B. ein 6522 oder ein 6551 können im Adressbereich angelegt werden und werden dann mit emuliert. (noch nicht implementiert)
- Emulation von IRQ und NMI. (noch nicht implementiert) 

## Start des Emulators

Der Emulator wird mit dem Befehl 

```
emu6502.exe -b rom.bin  
```

gestartet. Der Emulator startet automatisch im Einzelschrittmodus. Genau wie die 6502 CPU wird zunächst der Reset Vector ausgelesen und dann zu dieser Adresse gesprungen. In der aktuellen Version wird das ROM automatisch an die Adresse `0xe000` gelegt. (Wie bei der Hardware) Zusätzlich wird ein RAM Bereich von `0x0000-0x7fff` verwendet. Derzeit werden noch keine weiteren Baustein emuliert. 