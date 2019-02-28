@echo off
title Crawler

setlocal
set uac=~uac_permission_tmp_%random%
md "%SystemRoot%\system32\%uac%" 2>nul
if %errorlevel%==0 ( rd "%SystemRoot%\system32\%uac%" >nul 2>nul ) else (
    echo set uac = CreateObject^("Shell.Application"^)>"%temp%\%uac%.vbs"
    echo uac.ShellExecute "%~s0","","","runas",1 >>"%temp%\%uac%.vbs"
    echo WScript.Quit >>"%temp%\%uac%.vbs"
    "%temp%\%uac%.vbs" /f
    del /f /q "%temp%\%uac%.vbs" & exit )
endlocal  

:BG
cls
echo ©°©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©´
echo ©¦                        Crawler                              ©¦
echo ©À©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©È
echo ©¦                                                             ©¦
echo ©¦INS USE                                                      ©¦
echo ©¦        build      run go get and build                      ©¦
echo ©¦        run        run crawler                               ©¦
echo ©¦        install    install crawler as service (use nssm)     ©¦
echo ©¦        uninstall  uninstall crawler service                 ©¦
echo ©¦        start      start crawler service (after install)     ©¦
echo ©¦        stop       stop crawler service                      ©¦
echo ©¦        restart    stop and start crawler                    ©¦
echo ©¦        version    show crawler version                      ©¦
echo ©¦                                                             ©¦
echo ©¸©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¤©¼

%~d0
cd %~dp0
::SET select=
SET /P select="Please Enter Instructions:"
IF "%select%"=="build" (
    go get -v ./...
    go build -o %~dp0\bin\crawler.exe  %~dp0\src\main.go
    echo Build Finish.. 
    pause
    GOTO BG
) ELSE (
    IF "%select%"=="run" (
        %~dp0/bin/crawler.exe 
    ) ELSE ( 
        IF "%select%"=="install" (
            %~dp0\\bin\\nssm.exe install crawler %~dp0\\bin\\crawler.exe 
            pause
            GOTO BG
        ) ELSE ( 
            IF "%select%"=="start" (
                net start crawler 
                pause
                GOTO BG
            ) ELSE (
                IF "%select%"=="stop" (
                    net stop crawler 
                    pause
                    GOTO BG
                ) ELSE (
                    IF "%select%"=="restart" (
                        net stop crawler 
                        net start crawler 
                        pause
                        GOTO BG
                    ) ELSE (
                        IF "%select%"=="uninstall" (
                            sc delete crawler 
                            pause
                            GOTO BG
                        ) ELSE (
                             IF "%select%"=="version" (
                                %~dp0\bin\crawler.exe -v 
                                pause
                                GOTO BG
                            ) ELSE (
                                 echo Param Error Try Again!
                                 pause
                                 GOTO BG
                            )
                        ) 
                    ) 
                ) 
            ) 
        ) 
    )
)

pause

exit