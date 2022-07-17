@echo off
set LIBCUPL=c:\Wincupl\Shared\atmel.dl
c:\Wincupl\Shared\cupl.exe -jaxfsl %1.pld
rem c:\Wincupl\Shared\csim.exe -l g16v8 -u c:\Wincupl\Shared\Atmel.dl TEST