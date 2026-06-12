# Demo: discover and query metrics
Output assets/edx-metrics.gif
Set Shell "bash"
Set FontSize 15
Set Width 1200
Set Height 500
Set Padding 12
Set Theme "Catppuccin Mocha"
Set TypingSpeed 25ms

Type "edx metrics list --keyword cpu --output table --columns name,count"
Sleep 500ms
Enter
Sleep 3s
Type "edx metrics query --name k8s.node.cpu.seconds.rate --agg avg --group-by ed.cluster.name --lookback 15m --graph-type table"
Sleep 500ms
Enter
Sleep 4s
