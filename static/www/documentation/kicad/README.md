# kicad

## Linux Mint

Install the latest version of KiCAD on Linux Mint like so:

```sh
sudo add-apt-repository ppa:kicad/kicad-8.0-releases
sudo apt update
sudo apt install kicad
```

## Footprints

These are the `*.mod` and `*.kicad_mod` files used on the PCB. They
should show up in your project if you add a library using the footprint
editor.

* [FOOTPRINT LIBRARY FILE FORMAT](https://dev-docs.kicad.org/en/file-formats/sexpr-footprint/)
* Footprint library files use the .kicad_mod extension.
* Footprint library files can only define a single footprint.
* Footprint libraries are defined a folder containing one or more footprint library files.
