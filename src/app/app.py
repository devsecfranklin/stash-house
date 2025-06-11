"""Website."""

import logging

from flask import Flask, render_template
from mc_server import McServer

from __init__ import McUtils

logger = logging.getLogger("website")

my_utils = McUtils()
my_server = McServer(ip=my_utils.ip, port=my_utils.minecraft_port)

# Create an instance of the Flask class that is the WSGI application.
# The first argument is the name of the application module or package,
# typically __name__ when using a single module.
app = Flask(__name__)

# Flask route decorators map / and /hello to the hello function.
# To add other resources, create functions that generate the page contents
# and add decorators to define the appropriate resource locators for them.


@app.route("/")
@app.route("/minecraft")
def hello():
    return render_template("index.html")


@app.route("/history")
def my_history():
    """Generate a server history page"""
    logger.debug("Generate a server history page.")
    return render_template("history.html")


@app.route("/status")
def my_status():
    """Generate a server status page"""
    logger.debug("Generate a server status page.")
    return render_template("status.html", status=my_server.status, latency=my_server.latency, query=my_server.query)


if __name__ == "__main__":
    logger.debug("Booting up....")
    app.run(
        host=my_utils.ip, port=my_utils.flask_port, debug=True
    )  # Run the app server on localhost:9001


# SPDX-FileCopyrightText: 2014-2025 <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT
