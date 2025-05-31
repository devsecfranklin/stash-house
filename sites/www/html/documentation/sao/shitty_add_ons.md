# AND!XOR SAO Reference Designs

- [SHITTY ADD-ON V1.69BIS STANDARD](https://hackaday.com/2019/03/20/introducing-the-shitty-add-on-v1-69bis-standard/)

## Hardware

- [badge header connectors](https://www.ebay.com/itm/173960769032)
- search Ebay for `2×3 Pin 6 Pin With Polarizing Key`

![sao1](https://github.com/dead10c5/badge-project-template/blob/main/docs/images/sao1.png)

![sao2](https://github.com/dead10c5/badge-project-template/blob/main/docs/images/sao2.png)

![sao_connector](https://github.com/dead10c5/badge-project-template/blob/main/docs/images/sao_connector.jpg)

## I2C Addresses

| Type  | IC | 7-bit address |
|:-: | :-: |:-: |
| GPIO Expander | MCP23017 | 0x20 |
| EEPROM  | AT24C32  | 0x50    |
| REDACTED   | REDACTED  | TBD |

## EEPROM

Simple I2C EEPROM SAO based on the AT24C32. All address pins are pulled to GND. Badges that detect the AT24C32 on the I2C bus should query for bytes at known addresses then perform actions.

AND!XOR has adopted the following format for data stored in the EEPROM, we hope that others from #badgelife do too.

| 0   | 1   | 2    | 3...n  |
|:-:  |:-:  |:-:   |:-:  |
| DC Year | Maker ID | SAO Type ID  | Data   |

- DC Year: Use 0x1B for DC27
- Maker ID: Unique identifier for SAO maker, AND!XOR uses 0x49 (Middle 8bits of our registered Bluetooth ID)
- SAO Type ID: Unique identifier assigned by the maker for the SAO
- Data: Arbitrary data parseable by anything recognizing DC, Maker, and SAO values

## GPIO Expander

Returning from DC26 we will be using MCP23017-based SAO(s). When the MCP23017 is detected on the bus, the badge will play special patterns.
