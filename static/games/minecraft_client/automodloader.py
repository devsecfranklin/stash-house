from urllib.request import urlopen
import time
from pretty_downloader import pretty_downloader
import os

modlistlink = (
    "https://github.com/Rayrsn/Minecraft-Auto-Mod-Downloader/raw/main/mods.list"
)
versionLink = (
    "https://github.com/Rayrsn/Minecraft-Auto-Mod-Downloader/raw/main/modpack.ver"
)
AppData = os.environ["APPDATA"]
LauncherPath = "/.minecraft/TLauncher.exe"  # windows, but could be linux. OS detection?


def Launch():
    path = AppData + LauncherPath
    os.system(path)


def WriteModList():
    if os.path.exists(AppData + "/Python" + "/mods.list"):
        os.remove(AppData + "/Python" + "/mods.list")
    pretty_downloader.download(modlistlink, file_path=AppData + "/Python")


if not os.path.exists(AppData + "/Python"):
    os.makedirs(AppData + "/Python")


try:
    with open(AppData + "/Python/" + "modpack.ver", "x") as b:
        b.write("1")
        b.close
except:
    pass

with open(AppData + "/Python/" + "modpack.ver", "r") as x:
    local_ver = x.read()
    web_ver = str(urlopen(versionLink).read()).replace("b", "").replace("'", "")
    if local_ver != web_ver:
        os.system("cls")
        print("Local version doesn't match the fetched version")
        print("Updating...\n")
        WriteModList()
        with open(AppData + "/Python" + "/mods.list", "r") as f:
            lines = f.readlines()
            lines2 = [s.replace("\n", "") for s in lines]
            for line in lines2:
                if line.find("/"):
                    filename = line.rsplit("/", 1)[1]
                print("\nDownloading Mod:")
                print(line + "\n")
                pretty_downloader.download(line, file_path=AppData + "/.minecraft/mods")
            f.close()
            with open(AppData + "/Python/" + "modpack.ver", "w") as z:
                z.write(web_ver)
            print("Download Complete.")
            print("Launching game...")

    else:
        print("Mods are up to date!")
        print("Launching game...")
Launch()
time.sleep(5)
