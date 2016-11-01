# Learning Well - Code Challenge

## Översikt

Mjukvaran är implementerad med programspråket [Go][go], och
är terminalbaserad. Tanken med koden är att den ska vara så
idiomatisk som möjligt för valt programspråk.

## Exekvering

Detta projekt innehåller förkompilerade binärer för Mac OS X,
Linux samt Windows. Samtliga givet en 64-bitars arkitektur.

Exempelvis så heter binären för Mac OS X `votingrates-darwin`.

Givetvis så kan mjukvaran också byggas från källkoden direkt
om så önskas m.h.a. `go build` som producerar filen `votingrates`
(eller `votingrates.exe` för Windows). Detta förutsätter att
[Go][go] finns installerat på datorn.

En exekvering av mjukvaran ska producera ett liknande resultat
som exemplutskriften i PDF:en där uppgiften beskrivs, med skillnaden
att namnen på valdistrikten är sorterade i bokstavsordning.

## Testning

Testning sker genom `go test -v`

## Möjliga förändringar och förbättringar

I sin nuvarande form så hämtar programmet alltid data via SCB:s API.
Resultatet skulle kunna snabblagras per användare (p.g.a. skrivskydd)
så att data endast hämtas om då ett nytt val har genomförts eftersom
statistiken inte ändras under tiden.

[go]:http://golang.org/
