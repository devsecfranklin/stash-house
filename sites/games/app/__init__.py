import logging
import logging.config

logging.config.fileConfig(
    "logging.conf",
    defaults={"logfilename": "website.log"},
    disable_existing_loggers=False,
)

logger = logging.getLogger("hanson")


class McUtils:
    """Minecraft utilities."""

    def __init__(self):
        self.ip = "10.10.8.11"
        self.minecraft_port = "25565"
        self.flask_port = "9001"


# SPDX-FileCopyrightText: 2014-2025 <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT
