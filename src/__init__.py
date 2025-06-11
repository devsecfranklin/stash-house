# SPDX-FileCopyrightText: © 2022-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: GPL-3.0-or-later

import logging.config

logging.config.fileConfig(
    "logging.conf",
    defaults={"logfilename": "website.log"},
    disable_existing_loggers=False,
)
