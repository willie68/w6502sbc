@echo off
E:\Sprachen\6502\retroassembler\retroassembler.exe -c -C=65SC02 -g %1.asm

rem ca65 --debug -l %1.lst --cpu 65C02 %1.asm
rem ld65 %1.o -o %1.bin -t none