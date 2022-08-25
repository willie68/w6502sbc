@echo off
rem E:\Sprachen\6502\retroassembler\retroassembler.exe -c -C=65SC02 -g %1.asm

ca65 --debug -l %1.lst -I .\include\ %1.asm
ld65 %1.o -o %1.bin -C w6502sbc.cfg