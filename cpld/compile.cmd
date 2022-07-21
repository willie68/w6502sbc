@echo off
set LIBCUPL=e:\Wincupl\Shared\atmel.dl
e:\Wincupl\Shared\cupl.exe -jaxfsl %1.pld
rem c:\Wincupl\Shared\csim.exe -l g16v8 -u c:\Wincupl\Shared\Atmel.dl TEST
if errorlevel 1 (
%1.lst
%1.so
)