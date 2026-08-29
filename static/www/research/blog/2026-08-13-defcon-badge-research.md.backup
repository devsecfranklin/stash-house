# DEFCON 34 Badge Firmware Analysis - Lab Notes

- **Date:** 2026-08-13
- **Tags:** security, reverse-engineering, firmware, defcon
- **Analyst:** robot 🤖

## Architecture Overview
- **MCU:** RP2040 / Custom ESP32-S3 dual-core variant.
- **Interfaces:** USB-C, Exposed 4-pin UART (115200 8N1), SWD debug pins.

## Findings Summary
1. Boot1 loader lacks signature verification on secondary partition payloads.
2. UART diagnostic prompt allows arbitrary memory reads using opcodes `0x3A` and `0x3B`.
3. Challenge-response verification logic successfully decompiled via Ghidra.
