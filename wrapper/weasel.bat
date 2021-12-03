@REM This script is a wrapper for weasel that downloads and calls the executable for you

@REM -------------------------------- License -----------------------------------

@REM MIT License

@REM Copyright (c) 2021 Christoph Swoboda

@REM Permission is hereby granted, free of charge, to any person obtaining a copy
@REM of this software and associated documentation files (the "Software"), to deal
@REM in the Software without restriction, including without limitation the rights
@REM to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
@REM copies of the Software, and to permit persons to whom the Software is
@REM furnished to do so, subject to the following conditions:

@REM The above copyright notice and this permission notice shall be included in all
@REM copies or substantial portions of the Software.

@REM THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
@REM IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
@REM FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
@REM AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
@REM LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
@REM OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
@REM SOFTWARE.

@echo off

REM -------------------------------- Options -----------------------------------

REM Set the following option to enable build & run mode; if left empty just weasel is called
set BUILD_ARGS=

set WEASEL_VERSION=v0.1.0
set CACHE_PATH=%~dp0/.weasel


REM ------------------------------- Advanced -----------------------------------

set EXEC_NAME=weasel-windows-amd64.exe
set RELEASE_URL=https://github.com/codetent/weasel/releases/download/%VERSION%/%EXEC_NAME%


REM ----------------------------- DO NOT TOUCH ---------------------------------

set VERSION_FILE=%CACHE_PATH%/VERSION
set EXEC_FILE=%CACHE_PATH%/%EXEC_NAME%

if not exist "%CACHE_PATH%" ( 
    mkdir "%CACHE_PATH%"
)

:check-version
if exist "%VERSION_FILE%" (
    for /F "tokens=* delims=" %%x in (%VERSION_FILE%) do (
        if %%x == %WEASEL_VERSION% (
            goto :execute
        )
    )
)

:download
curl --silent --show-error --fail --output "%EXEC_FILE%" "%RELEASE_URL%"
if not %ERRORLEVEL% == 0 (
    echo Error: The weasel executable could not be fetched.
    exit /B 1
)

echo %WEASEL_VERSION%> "%VERSION_FILE%"

:execute
if not exist "%EXEC_FILE%" ( 
    echo Error: The weasel executable could not be found.
    exit /B 1
)

if "%BUILD_ARGS%" == "" (
    "%EXEC_FILE%" %*
) else (
    "%EXEC_FILE%" build --errors-only --tag test %BUILD_ARGS%
    "%EXEC_FILE%" run test %*
)
