.PHONY: slides
slides:
	caddy run

.PHONY: backend
backend:
	go build -o backend-server ./backend

.PHONY: frontend
frontend:
	cd frontend && flask run
