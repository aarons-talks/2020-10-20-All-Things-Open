from app import app, PICKLE_DB_FILENAME
import pytest
import os

@pytest.fixture
def flask_app(request):
    app.config["TESTING"] = True
    # pre-create the pickle DB file
    open(PICKLE_DB_FILENAME, "w+")

    with app.test_client() as client:
        yield client
    
    os.remove(PICKLE_DB_FILENAME)
