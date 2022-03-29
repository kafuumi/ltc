echo off
@REM forceposix 表示在windows上参数也为linux风格，即以“-”开头
go build -tags="forceposix" -ldflags "-s -w" -o ltc.exe .
