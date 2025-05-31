# Board outline

This describes how to make the PCB into the correct "shape" in KiCad. There is an
"edge layer" in the PCB design. We use tools to trace and scale the original artwork,
and then import the trace into the Edge Layer.

## Outline the Image

* Import the jpg/png to inkscape.
* Use the "Draw bezier curves and Straight lines" tool to outline the image completely.
* Click on "Object => Transform -> Scale" to resize the image proportionally.
  * Make sure the "scale proportionally" box is checked.
* Save the outline as a "plain SVG", and save it that way (rather than inkscape SVG)
* [Outline the image](https://www.youtube.com/watch?v=lcmY77oxd8w)

## Import the SVG to KiCad

* We can now bring in the outline to the edge cuts layer on the PCB.
* Make sure the "edge cuts" Layer is selected.
* Click on "file -> import -> graphics" and select the plain SVG file.
