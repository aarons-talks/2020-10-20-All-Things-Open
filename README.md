# 2020-10-20-All-Things-Open
My slides and code for All Things Open 2020


# Python v. Go Notes

- Black code formatting - closest thing to go fmt
- mypy type checking (there are others too) = closest thing to static type checks that compilers/static languages give you
- **consider showing go test generation**
  - maybe to show how much boilerplate tests can have?
  - maybe to show the tooling is great
  - CODE GENERATION IN GO IS A HUGE ADVANTAGE
- Python has property based testing (https://hypothesis.readthedocs.io/en/latest/)
- There's another way to connect Go / Python together, we won't talk about this, but it's very cool. A fellow Twitch Live Coder (Shoutout anthonywritescode) built this: https://github.com/asottile/setuptools-golang-examples
- We're gonna separate out the concept of a "frontend tier" (server that ties together lots of data sources, backend systems, etc...) and a "backend tier" (servers that do specific tasks) - kind of like the unix philosophy, where backends do one thing well and frontend just "pipes" them all together

TODO: try out pipenv - noobypirate said so lol
