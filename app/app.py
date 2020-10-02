from flask import Flask, render_template, jsonify, request
import pickle
import requests
import uuid

app = Flask(__name__)

@app.route("/")
def home():
    return jsonify({
        "images": [],
        "num_images": 0,
        "last_uploaded": "",
    })
@app.route("/stats")
def stats():
    return jsonify({
        "avg_size": 0,
        "avg_compressed_size": 0,
        "num_tags": 0,
    })

@app.route("/size_histogram")
def size_histogram():
    return jsonify({})

@app.route("/tags")
def tags():
    return jsonify({})
    
@app.route("/upload", methods=["POST"])
def upload():

    """
    TODO: call out to the Go backend to:
    download image from URL
    compress
    index & tag
    return continuation token
    
    Also TODO: need to figure out the communication between this app and the Go app
    """

    json = request.get_json()
    if not json:
        return "No JSON found", 400
        
    # TODO: make sure these are the right data types
    image_url = json["url"] # needs to be a string
    image_tags = json["tags"] # needs to be a list of strings
    image_name = json["name"] # needs to be a string

    image_binary = requests.get(image_url)

    with open("pickled_db.db", "w+") as dbFile:
        filename = "{}-{}".format(image_name, str(uuid.uuid4()))
        with open(filename, "w") as imageFile:
            imageFile.write(image_binary)

        unpickled = pickle.load(dbFile)
        if not isinstance(unpickled["images"], list):
            unpickled["images"] = []
        
        unpickled.push({
            str(uuid.uuid4()): {
                "name": image_name,
                "tags": image_tags,
                "url": image_url,
                "file_location": "./{}".format(filename)
            }
        })

    return jsonify({
        # "token": "abvcd"
        "status": "done"
    })
    
if __name__ == "__main__":
    app.run(debug=True)
