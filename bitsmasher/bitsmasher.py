# -*- coding: utf-8 -*-
"""
  bitsmasher.net
  :copyright: (c) 2019 by Franklin Diaz
  :license: MIT
"""
import os
import flask
from flask import Flask
from flask import render_template
from flask import send_from_directory

app = Flask(__name__)

def create_app(debug):
    app.debug = debug
    return app
    
@app.route('/')
def render_static():
    return render_template('index.html', title='bitsmasher.net')


#@app.route('/indybsides')
#@app.route('/bsidesindy')

@app.route('/favicon.ico')
def favicon():
    return send_from_directory(os.path.join(app.root_path, 'static'), 'favicon.ico',mimetype='image/vnd.microsoft.icon')

@app.errorhandler(404)
def page_not_found(e):
    return render_template('404.php'), 404
   
if __name__ == '__main__':
    app = create_app(debug=True)
    app.run(host="0.0.0.0")