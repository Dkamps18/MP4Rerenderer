@echo off
setlocal enabledelayedexpansion

IF exist build/ (rmdir /s/q build)
mkdir build

for /f "delims=" %%x in (version) do set version=%%x

SET GOOS=linux
SET GOARCH=amd64
go build -o build/mp4rerenderer
tar -C build -a -cf build/mp4rerenderer-%version%-Linux-%GOARCH%.tar.gz mp4rerenderer
set idx=0
for /f %%F in ('certutil -hashfile build/mp4rerenderer-%version%-Linux-%GOARCH%.tar.gz SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo mp4rerenderer-%version%-Linux-%GOARCH%.tar.gz  %out1% >> build/checksum.txt
SET GOARCH=arm64
go build -o build/mp4rerenderer
tar -C build -a -cf build/mp4rerenderer-%version%-Linux-%GOARCH%.tar.gz mp4rerenderer
set idx=0
for /f %%F in ('certutil -hashfile build/mp4rerenderer-%version%-Linux-%GOARCH%.tar.gz SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo mp4rerenderer-%version%-Linux-%GOARCH%.tar.gz  %out1% >> build/checksum.txt

SET GOOS=windows
SET GOARCH=amd64
go build -o build/mp4rerenderer.exe
tar -C build -a -cf build/mp4rerenderer-%version%-Windows-%GOARCH%.zip mp4rerenderer.exe
set idx=0
for /f %%F in ('certutil -hashfile build/mp4rerenderer-%version%-Windows-%GOARCH%.zip SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo mp4rerenderer-%version%-Windows-%GOARCH%.zip   %out1% >> build/checksum.txt
SET GOARCH=arm64
go build -o build/mp4rerenderer.exe
tar -C build -a -cf build/mp4rerenderer-%version%-Windows-%GOARCH%.zip mp4rerenderer.exe
set idx=0
for /f %%F in ('certutil -hashfile build/mp4rerenderer-%version%-Windows-%GOARCH%.zip SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo mp4rerenderer-%version%-Windows-%GOARCH%.zip   %out1% >> build/checksum.txt

SET GOOS=darwin
SET GOARCH=amd64
go build -o build/mp4rerenderer
tar -C build -a -cf build/mp4rerenderer-%version%-macOS-%GOARCH%.tar.gz mp4rerenderer
set idx=0
for /f %%F in ('certutil -hashfile build/mp4rerenderer-%version%-macOS-%GOARCH%.tar.gz SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo mp4rerenderer-%version%-macOS-%GOARCH%.tar.gz  %out1% >> build/checksum.txt
SET GOARCH=arm64
go build -o build/mp4rerenderer
tar -C build -a -cf build/mp4rerenderer-%version%-macOS-%GOARCH%.tar.gz mp4rerenderer
set idx=0
for /f %%F in ('certutil -hashfile build/mp4rerenderer-%version%-macOS-%GOARCH%.tar.gz SHA256') do (
    set "out!idx!=%%F"
    set /a idx += 1
)
echo mp4rerenderer-%version%-macOS-%GOARCH%.tar.gz  %out1% >> build/checksum.txt

del build\mp4rerenderer
del build\mp4rerenderer.exe
