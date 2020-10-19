# 2020-10-20-All-Things-Open

My slides and code for All Things Open 2020

There are two main apps in this repository:

- A frontend Python app that serves up a limited REST API and HTML templates
- A backend Go app that does the expensive tasks like downloading and compressing images

# Testing the Python App

The Python app has two "entrypoint" files - `app.py` and `app.backup.py`. `app.py` does all of the work itself - downloading the image, saving it, indexing it in the database, and so on. It does not do any compression because that is very computationally expensive to do

`app.backup.py` relies on the Go backend (which is in [backend](./backend)) to do the compression work.

## The All In One App

To run the "all in one" app, follow the below instructions

First, go into the directory with the app in it:

```shell
cd frontend
```

Then, set up the Python virtual environment:

```shell
source ./start-env.sh
```

Start the app with:

```shell
flask run
```

Now, you can upload an image:

```shell
curl -XPOST -H "Content-Type: application/json" -d '{"url": "https://www.vizteams.com/wp-content/uploads/2013/08/python-logo-master.png", "tags": ["not", "found"], "name":"python"}' localhost:5000/upload
```

And then test that it worked properly:

```shell
curl localhost:5000
```

You should see something like this:

```json
{
  "last_uploaded": "Thu, 15 Oct 2020 22:10:19 GMT",
  "num_images": 1
}
```


# Notes from Streamers

I built these slides, the code and this talk in general on my [Twitch stream](https://twitch.tv/arschles). The following are notes and ideas that chat had along the way. Not all of this made it into the talk, but I believe that everything here is relevant and interesting.

- Black code formatting - closest thing to go fmt
- mypy type checking (there are others too) = closest thing to static type checks that compilers/static languages give you
- **consider showing go test generation**
  - maybe to show how much boilerplate tests can have?
  - maybe to show the tooling is great
  - CODE GENERATION IN GO IS A HUGE ADVANTAGE
- Python has property based testing (https://hypothesis.readthedocs.io/en/latest/)
- There's another way to connect Go / Python together, we won't talk about this, but it's very cool. A fellow Twitch Live Coder (Shoutout anthonywritescode) built this: https://github.com/asottile/setuptools-golang-examples
- We're gonna separate out the concept of a "frontend tier" (server that ties together lots of data sources, backend systems, etc...) and a "backend tier" (servers that do specific tasks) - kind of like the unix philosophy, where backends do one thing well and frontend just "pipes" them all together
- Environment management & dependencies - try out pipenv - noobypirate said so lol
