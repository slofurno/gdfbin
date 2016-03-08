## Example

```
gdfbin git:(master)✗ gdf3 ls
resume             BZ7J3RVB	2 days ago
gdf3               NJ7JKCR1	5 days ago
package-component  77JHFYBJ	6 days ago
webpack-component  6YMZEZTZ	6 days ago
index.html         R91TTENC	7 days ago
brb                DPBHHQH5	8 days ago
makefile           7Y287QN1	8 days ago
gitignore          1BN7FKJG	8 days ago
webpack            KR0B13Z5	8 days ago
.vimrc             HETQ6PMM	8 days ago
bookmarks.go       2NFHQ95Y	8 days ago
```

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
