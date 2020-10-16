from flask import Flask, render_template, jsonify, request
import requests
import uuid
import datetime
import os


# TODO: maybe put this into an app factory:
# https://flask.palletsprojects.com/en/1.1.x/patterns/appfactories/

app = Flask(__name__)
BACKEND_HOST="http://localhost:5001"

@app.route("/")
def home():
    # get the basic stats back from Go backend
    # do we convert this from json to a python dict??
    resp = requests.get("{}/basic_stats".format(BACKEND_HOST))
    if resp.status_code != 200:
        return "Error: {}".format(resp.text), resp.status_code
    ret_json = resp.json()
    return jsonify(ret_json)

@app.route("/image/<image_name>")
def view_image(image_name):
    img_data = requests.get("{}/image/{}".format(BACKEND_HOST, image_name)).json()
    return render_template(
        "/image.html",
        img_src="http://localhost:5001{}".format(img_data["src"]),
        img_alt=img_data["alt"]
    )

@app.route("/upload", methods = ['POST'])
def upload():
    json = request.get_json()
    # We technically don't need to do this because the Go backend
    # can check the json for us automatically. This is because it has
    # types, and the JSON library over there will automatically check
    # the types of the incoming JSON for us
    if ("url" not in json) or ("tags" not in json) or ("name" not in json):
        return "Invalid JSON", 400
    
    payload = {
        "url": json["url"],
        "tags": json["tags"],
        "name": json["name"],
    }

    # this is where we return the 201 CREATED response because the
    # Go backend begins processing the image but immediately returns 
    # to us
    res = requests.post("{}/process_image".format(BACKEND_HOST), data=payload)

    return jsonify({
        # "token": "abvcd"
        "status": "done"
    })
    
if __name__ == "__main__":
    if os.environ.get("APP_BACKEND_HOST") is not None:
        BACKEND_HOST = os.environ.get("APP_BACKEND_HOST")
    app.run(debug=True)
