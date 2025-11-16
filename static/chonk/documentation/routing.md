# Routing PCB Traces

- Download [Layout Editor](https://layouteditor.com/download)
- Use the [Free Routing tool](https://freerouting.org/freerouting/using-with-kicad) to route KiCad boards.

```sh
locate freeRouting.jar
java -jar /opt/layout/bin/freeRouting.jar
```

- From KiCad, Create the DSN file by clicking on File, Export, “Specctra DSN…”
- Run the tool: `java -jar /opt/layout/bin/freeRouting.jar`
- Open the DSN file with the tool.
- Click on the `Routing` drop down and select `Autorouting`
- To export from FreeRouting, click on File, Export, Export Specctra Session File.
