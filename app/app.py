from flask import Flask, render_template, jsonify, request
import pickle
import requests
import uuid
import datetime
from io import BytesIO

PICKLE_DB_FILENAME = "pickled_db.db"

# TODO: maybe put this into an app factory:
# https://flask.palletsprojects.com/en/1.1.x/patterns/appfactories/

app = Flask(__name__)

@app.route("/")
def home():
    try:
        with open(PICKLE_DB_FILENAME, "rb") as dbFile:
            unpickled = pickle.load(dbFile)
            return jsonify({
                "num_images": unpickled["num_images"],
                "last_uploaded": unpickled["last_uploaded"],
            })
    except:
        return jsonify({
            "num_images": 0,
            "last_uploaded": "never"
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
    
@app.route("/upload", methods = ['POST'])
def upload():

    """
    TODO: call out to the Go backend to:
    compress
    index & tag
    return continuation token
    
    Also TODO: need to figure out the communication between this app and the Go app
    """

    json = request.get_json()
    if not json:
        return "No JSON found", 400
    
    print("In the upload function")

    # TODO: make sure these are the right data types
    image_url = json["url"] # needs to be a string
    image_tags = json["tags"] # needs to be a list of strings
    image_name = json["name"] # needs to be a string

    image_binary = requests.get(image_url).content
    image_size = len(image_binary)

    # TODO: we cannot do a compress operation concurrently because
    # Python doesn't support multicore operations
    #
    # We would have to use multiprocessing to do it, but it would be nicer
    # to not need to fork a new process each time we get a new request, 
    # because that's pretty heavyweight
    

    # FOR LIVE CODING:
    # rip out all of this low level DB etc... code
    # and replace it with a call to the Go backend server

    unpickled = None
    with open("pickled_db.db", "rb") as dbFile:
        filename = "{}-{}.image".format(image_name, str(uuid.uuid4()))
        with open(filename, "wb") as imageFile:
            imageFile.write(image_binary)

        # instantiate the unpickled data if the file
        # was empty
        try:
            unpickled = pickle.load(dbFile)
        except:
            unpickled = {
                "images": [],
                "num_images": 0,
                "total_size": 0,
            }
        
        unpickled["last_uploaded"] = datetime.datetime.now(datetime.timezone.utc)
        unpickled["num_images"] = unpickled["num_images"] + 1
        unpickled["total_size"] = unpickled["total_size"] + image_size
        unpickled["avg_size"] = unpickled["total_size"] / unpickled["num_images"]

        unpickled["images"].append({
            str(uuid.uuid4()): {
                "name": image_name,
                "tags": image_tags,
                "url": image_url,
                "file_location": "./{}".format(filename)
            }
        })
    
    with open(PICKLE_DB_FILENAME, "wb") as f:
        pickle.dump(unpickled, f)

    return jsonify({
        # "token": "abvcd"
        "status": "done"
    })
    
if __name__ == "__main__":
    # create the pickle DB file so that we can read it later,
    # without having to check if it already exists
    open(PICKLE_DB_FILENAME, 'w+')
    app.run(debug=True)
