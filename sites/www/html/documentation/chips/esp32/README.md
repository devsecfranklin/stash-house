# ESP32

## Install and set up the ESP-IDF on your machine

The IDF is a powerful and efficient framework used to program the
ESP32 family of microcontrollers
using C or C++. The ESP-IDF currently powers millions of devices in the field and enables building
a variety of network-connected products, ranging from simple light bulbs and toys to big appliances
and industrial devices.

Add the following line at the end of the ~/.config/fish/config.fish file:

```sh
direnv hook fish | source
```

```sh
sudo apt purge -y python2.7-minimal
sudo apt-get install git wget flex bison gperf python3 python3-pip python3-venv cmake ninja-build ccache libffi-dev libssl-dev dfu-util libusb-1.0-0 python3-full
mkdir -p ~/esp
cd ~/esp
git clone -b v5.2 --recursive https://github.com/espressif/esp-idf.git
cd ~/esp/esp-idf
./install.fish all
idf.py set-target esp32 # from src directory
```

## The main ESP32 chip types, their differences and how to select the one you need

## Use logging

## Basic GPIO’s with blinkey

## Keyboard input

## Using ESP-IDF examples

## Use C concepts with a primer on pointers, structures, and callbacks

## Create your own libraries and use other libraries

## Use FreeRTOS to create, manage and synchronise multiple tasks

## Debug with break points

## Use and create configuration

## Manage and analyse memory

## Persist data in Flash or SD card

## Use GPIO’s; including inputs, outputs, Interrupts, Touch ADC etc

## Use communications protocols with UART, I2C and SPI

## Use sleep and deep sleep to preserve power

## Connect to Wi-Fi or set up your own Wi-Fi access point

## Interact with Web API’s on the internet

## Create and host your own webserver to serve web pages written in modern web frameworks such as React

## Use MQTT to push and pull messages over the internet

## Use ESP-NOW for chip-to-chip communications

## Update your firmware with OTA (Over The Air) updates

## Use Bluetooth Low Energy