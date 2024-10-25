tailwind:
	@echo "Generating tailwind css"
	./tailwindcss -i ./static/css/input.css -o ./static/css/styles.min.css --watch --minify

templ-watch:
	@echo "Watching templ files"
	templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."

dev:
	make -j2 tailwind templ-watch
		
		
