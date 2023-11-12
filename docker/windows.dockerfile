FROM chocolatey/choco:latest-windows

RUN choco install cygwin --yes

RUN choco install cyg-get --yes

CMD PowerShell.exe