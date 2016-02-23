## Usage

### Using curl 

```
cat index.js | curl --data-binary @- https://gdf3.com

or

curl --data-binary @index.js gdf3.com
```

### Using the cli

run install

```
./install.sh
```

register an account

```
gdf3 register <email> <password>
```

paste and bookmark a file

```
cat webpack.config.js | gdf3 mark webpack
```

list bookmarks

```
gdf3 ls
```

download a bookmarked paste

```
gdf3 get webpack > webpack.config.js
```


## Hosting 

### init the sqlite database

```
sqlite3 pastes.db ".read schema.sql"
```

install go, get deps, build, and run
