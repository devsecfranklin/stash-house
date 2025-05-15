#!/usr/bin/env python

from setuptools import setup

with open("README.md", encoding="ascii") as file:
    long_description = file.read()

setup(
    name="mcrcon",
    version="0.7.0",
    description="A client for handling Remote Commands (RCON) to a "
    "Minecraft server.",
    long_description=long_description,
    url="https://github.com/uncaught-exceptions/mcrcon",
    author="Adrian Turjak",
    author_email="minecraft@uncaught-exceptions.com",
    license="MIT",
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Developers",
        "Intended Audience :: System Administrators",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3.7",
    ],
    keywords="minecraft mcrcon rcon remote command mojang cli",
    py_modules=["mcrcon"],
    entry_points={"console_scripts": ["mcrcon=mcrcon:mcrcon_cli"]},
)
