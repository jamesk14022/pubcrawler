## 🎡 Citycrawler 🍺

Web app that uses graph theory to generate custom “city crawl” routes in EU and US cities. This was built mostly as as a Golang and graph theory learning exercise. These walking routes encourage you to explore the finest drinking establishments in a city. Some features:

- A scraped "pub directory" containing 50,000+ relevant pubs from major cities
- Customizable Routes: Users can generate routes with adjustable lengths and select the number of bars included. They can also choose to include cultural areas of interest in their route
- Fully-featured mobile site 
- High performance DFS for fast route finding

## Dependencies

- **Docker:** [Install](https://docs.docker.com/get-docker/)
- **Golang (v1.22.3):** [Install](https://go.dev/dl/)
- **Golang Libraries:** Install with `go mod tidy`.

- **Frontend:** Install with `cd web && npm install`

## Environment Variables 

- WEB_DIR=/usr/local/web/static
- LOCATION_DATA_DIR=/usr/local/data/location_data/
- CACHE_DATA_PATH=/usr/local/data/location_data/cache.json
- MONGO_URI=

## Run

```
# build executable
make build

# run webserver to serve static files and API
make run

# tailwind monitor for changes and rebuild
make tailwind
```

## Developerment Commands

```
# format code 
make format

# scrape for new pubcrawl locations
make scrape
```

![alt text](https://github.com/jamesk14022/barcrawler/blob/main/web/src/assets/screenshot.png)
