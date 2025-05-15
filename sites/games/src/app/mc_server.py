"""Minecraft server class."""

from mcstatus import JavaServer


class McServer:
    """McServer."""


    def __init__(self, ip, port):
        # You can pass the same address you'd enter into the address field in minecraft into the 'lookup' function
        # If you know the host and port, you may skip this and use JavaServer("example.org", 1234)
        self.server = JavaServer.lookup(ip + ":" + port)
        self.latency = self.server.ping()
        self.status = self.server.status()
        self.query = status.players.online
        
    def mc_status(self):
        # 'status' is supported by all Minecraft servers that are version 1.7 or higher.
        # Don't expect the player list to always be complete, because many servers run
        # plugins that hide this information or limit the number of players returned or even
        # alter this list to contain fake players for purposes of having a custom message here.
        result = f"The server has {self.query} player(s) online and replied in {self.latency} ms"
        print(result)
        logger.debug(result)
        return result

    def mc_latency(self):
        # 'ping' is supported by all Minecraft servers that are version 1.7 or higher.
        # It is included in a 'status' call, but is also exposed separate if you do not require the additional info.
        result = f"The server replied in {self.__class__.latency} ms"
        print(result)
        logger.debug(result)
        return result

    def mc_query(self):
        # 'query' has to be enabled in a server's server.properties file!
        # It may give more information than a ping, such as a full player list or mod information.
        result = f"The server has the following players online: {', '.join(self.__class__.query.players.names)}"
        print(result)
        logger.debug(result)
        return result


# SPDX-FileCopyrightText: 2014-2025 <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT
