# Demo: search logs with CQL
# Requires an authenticated edx (edx auth login). No secrets appear on screen.
Output assets/edx-logs.gif
Set Shell "bash"
Set FontSize 15
Set Width 1200
Set Height 300
Set Padding 12
Set Theme "Catppuccin Mocha"
Set TypingSpeed 25ms

Type `edx logs search -q 'service.name:"admin-api-deployment" AND severity_text:"ERROR"' --lookback 1h --limit 5 --output table --columns severity_text,resource.service.name,resource.ed.tag`
Sleep 600ms
Enter
Sleep 4.5s
