# gfu - gophermap format utility
`gfu` manipulates gophermaps (gopher menus). It is intended to be used as part of an automation chain for managing a gopherhole using maps as the main doctype, easily allowing any document to contain links.

`gfu` can:
- Convert all lines in a file that are not valid gopher links into gopher info text (item type 'i') lines
- Deconstructing all gopher info text lines in a file back to plain text for easy editing

There are also plans to include additional features, such as:
- Adding the contents of a header file into the gophermap
- Adding the contents of a footer file into the gophermap
*Please note - Many servers already support includes, so the above may not be needed. If you are interested in this feature, you may want to check your server documentation first.*

## Getting Started
These instructions will get a copy of the project up and running on your local machine.

### Prerequisites
If building from source, you will need to have [Go](https://golang.org/) version 1.10 installed.

### Installing
Assuming you have `go` installed, run the following:
```
git clone https://tildegit.org/sloum/gfu.git
cd gfu
make install
```
Assuming `go install` is set up to install to a place on your path, you should be able to execute `gfu` from the terminal:
```
gfu
```

#### Troubleshooting
If you run `gfu` and get `gfu: command not found`, try running `make` from within the cloned repo. Then try: `./gfu`. If that works it means that Go does not install to your path. `make` added an executable file to the repo directory. Move that file to somewhere on your path. I suggest `/usr/local/bin` on most systems, but that may be a matter of personal preference.

### Downloading
If you would prefer to download a binary for your system, rather than build from source, please visit the [gfu downloads](https://rawtext.club/~sloum/gfu.html#downloads) page.

### Usage
#### Syntax
`gfu [flags...] [filepath]` 

#### Examples
`gfu ~/gopher/phlog/gophermap`
Convert plain text lines to gophermap info text (item type 'i') lines

`gfu -d ~/gopher/phlog/gophermap`
Deconstruct a gophermap's info text (item type 'i') lines back to plain text

`gfu -stdout ~/gopher/phlog/gophermap`
Don't write changes to file, only print them to the console

### Documentation
Please see the [gfu homepage](https://rawtext.club/~sloum/gfu.html#docs) for more information about `gfu`.

## Contributing
If you would like to get involved, please submit an issue. At present the developers use the tildegit issues system to discuss new features, track bugs, and communicate with users about hopes and/or issues for/with the software.

## License
This project is not currently licensed.
