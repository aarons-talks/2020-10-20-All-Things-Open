import pytest
from app import app, PICKLE_DB_FILENAME
import os.path
import pickle

def test_home_no_images(flask_app):
    resp = flask_app.get("/")
    assert resp.status_code == 200
    json = resp.get_json()
    assert json["last_uploaded"] == "never"
    assert json["num_images"] == 0

def test_home_some_images(flask_app):
    # upload an image
    resp = flask_app.post("/upload", json=dict({
        "url": "https://http.cat/404",
        "tags": ["http", "cat", "404"],
        "name": "httpcat"
    }))
    assert os.path.isfile("./pickled_db.db") == True
    assert resp.status_code == 200
    json = resp.get_json()
    assert json["status"] == "done"

    with open(PICKLE_DB_FILENAME, "rb") as f:
        unpickled = pickle.load(f)
        assert unpickled["last_uploaded"] != "never"
        assert unpickled["num_images"] == 1


    resp = flask_app.get('/')
    json = resp.get_json()
    assert json["num_images"] == 1
    assert json["last_uploaded"] != "never"
