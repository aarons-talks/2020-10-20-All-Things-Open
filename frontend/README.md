# Test Python App

You need Python 3 to run this.

Get the environment set up:

```shell
source ./start-env.sh
```

Start the server:

```shell
flask
```

Upload an image:

```shell
curl -X POST -H "Content-Type: application/json" -d '{"url": "https://http.cat/404", "name": "cat404", ["cat", "notfound"]}' http://localhost:5000
```
