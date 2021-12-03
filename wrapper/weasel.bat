@REM This script is a wrapper for weasel that downloads and calls the executable for you

@REM -------------------------------- License -----------------------------------

@REM Copyright © 2021 Christoph Swoboda

@REM Licensed under the Apache License, Version 2.0 (the "License");
@REM you may not use this file except in compliance with the License.
@REM You may obtain a copy of the License at

@REM     http://www.apache.org/licenses/LICENSE-2.0

@REM Unless required by applicable law or agreed to in writing, software
@REM distributed under the License is distributed on an "AS IS" BASIS,
@REM WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
@REM See the License for the specific language governing permissions and
@REM limitations under the License.

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
