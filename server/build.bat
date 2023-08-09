@echo off
setlocal

for /f %%i in ('git rev-parse HEAD') do (
    set "ver=%%i"
)

echo %ver%

docker build --build-arg CURRENTVERSION=%ver% -t gva-server:1.0 .

endlocal