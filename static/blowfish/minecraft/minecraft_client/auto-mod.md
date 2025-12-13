# Minecraft Auto Mod Downloader

A Python script which automatically downloads mods from a given list.

## Installing Dependencies

`pip install pretty-downloader`

## Requirements

1. Upload your version of "mods.list" and "modpack.ver".
  * I suggest uploading them to a github repo because the script requires direct links to them.
2. The "mods.list" file has to only include the ***Direct*** links to the mods (basically it has to look like this):

![mods.list](https://github.com/Rayrsn/Minecraft-Auto-Mod-Downloader/raw/main/images/mods.list.png?raw=true)

3. The "modpack.ver" file also has to follow the same rule and it has to look like this:

![modpack.ver](https://github.com/Rayrsn/Minecraft-Auto-Mod-Downloader/raw/main/images/modpack.ver.png?raw=true)

* Edit "modlistlink" and "versionLink" to the correct format.

### Setting up and Running the script from source

1. Open the script with a text editor.
2. Edit lines 6 and 7 ("modlistlink", "versionLink") to the correct values. 
  (Read [Requirements](https://github.com/Rayrsn/Minecraft-Auto-Mod-Downloader/blob/main/README.md#requirements).)
3. Change line 11 to your preferred launcher. (You might have to change the "AppData" values if the launcher is ***not*** located in the "AppData/Roaming/.minecraft" directory.)
4. Execute the script by running:

```sh
python main.py
```

## #Making an executable (.exe file)

1. Install Pyinstaller `pip install pyinstaller`.
2. Open the script directory in your preferred terminal.
3. Run this line of code: `pyinstaller --onefile main.py`
4. Your file will be located in the "dist" directory.

* You can also give it an icon by adding `--icon=<Icon_Name>.ico` parameter to this line.
* You can replace your minecraft shortcut with this .exe file

## Questions ⁉️

If yall have any questions or just wanna talk, add me on
[Discord](https://rayr.ml/LinkInBio)
