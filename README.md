# upvideo server

## Installation for development
1. Clone repo
2. Create/Update **config.json** (or other like dev.config.json) in root directory (see config.json.example)
3. Create/Update **dbconf.yml** in **db/** directory (see db/dbconf.yml.example file)
4. Install **goose** migration tool if not exists (go get bitbucket.org/liamstask/goose/cmd/goose)
5. Run migrations 
   ```bash
   goose up
   ```
6. Run/Build server like:
   ```bash
   go run server.go dev.config.json
   # or
   # go build 
   # ./upvideo dev.config.json
   ```