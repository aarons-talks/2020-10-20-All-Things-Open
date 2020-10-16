.PHONY: slides
slides:
	caddy run

.PHONY: backend
backend:
	go build -o backend-server ./backend

.PHONY: frontend
frontend:
	cd frontend && flask run

.PHONY: clean
clean:
	rm -f backend/imagefiles/*.image
	rm -f backend/images.db
	rm -f frontend/pickled.db
