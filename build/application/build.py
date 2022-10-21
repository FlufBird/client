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
    "darwin": ["", [
        "amd64",
        "arm64"
    ]],
}

DEVELOPMENT_MODE = os.path.isfile()

os.chdir("../../dist")

for platform in PLATFORMS:
    make_directory(platform)
    os.chdir(platform)

    for architecture in PLATFORMS[platform][1]:
        directory = f"../dist/{platform}/{architecture}"

        make_directory(architecture)
        os.chdir(architecture)

        # build with `wails build` here

        if not directory_exists("resources"):
            shutil.copytree("../../../resources", "resources")

        os.chdir("..")

    os.chdir("..")