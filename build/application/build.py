# working with these paths are so difficult, so i just bruteforced them

import os, os.path
import shutil

def directory_exists(name : str) -> bool:
    return os.path.isdir(name)

def make_directory(name : str) -> None:
    if not directory_exists(name):
        os.mkdir(name)

PLATFORMS = {
    "windows": ["exe", [
        "amd64",
    ]],
    "linux": ["", [
        "amd64",
        "arm64",
    ]],
}

os.chdir("../../dist")

for platform in PLATFORMS:
    make_directory(platform)
    os.chdir(platform)

    for architecture in PLATFORMS[platform][1]:
        _directory = f"dist/{platform}/{architecture}"
        directory = f"../{_directory}"

        make_directory(architecture)
        os.chdir(architecture)

        os.chdir("../../..")

        os.system(f"""wails build -clean -u -v 2 -platform "{platform}/{architecture}" -webview2 "embed" -trimpath "true" -o "{_directory}/flufbird.exe" """) # TODO: more flags

        os.chdir(_directory)

        if not directory_exists("resources"):
            shutil.copytree("../../../resources", "resources")

        os.chdir("..")

    os.chdir("..")