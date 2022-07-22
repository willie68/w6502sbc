@echo off
rem E:\Sprachen\6502\Ophis\ophis.exe -l %1.lst -c  %1.asm -o %1.o

ca65 --debug -l %1.lst --cpu 65C02 %1.asm
ld65 %1.o -o %1.bin -t none