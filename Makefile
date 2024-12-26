build: 
	go build -o bin/main cmd/api/main.go

run:
	go run cmd/api/main.go

scrape:
	go run cmd/scraper/main.go

tune:
	go run tuning/parameter_tuning.go

migrate_flatfiles:
	go run cmd/migrate/main.go

format: 
	gofmt -s -w .
	cd web/src && npx prettier --write . --ignore-path .prettierignore 

tailwind:
	cd web/src && npx tailwindcss -i ./styles.css -o ../dist/output.css --watch
