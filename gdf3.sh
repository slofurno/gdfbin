#!/bin/sh

url="https://gdf3.com"

usage (){
  echo "gdf3 register <email> <password>"
  echo "gdf3 login <email> <password>"
  echo "cat <file> | gdf3 mark <name>"
  echo "cat <file> | gdf3 paste"
  echo "gdf3 ls"
  echo "gdf3 get <name>"
}

token=`cat .gdf3.token 2> /dev/null`

case $1 in
  mark )
    name=`curl --data-binary @- $url | sed 's|https://gdf3.com/||'`
    curl -X POST "$url/bookmarks/$name/$2" -H "Auth: $token"
    ;;
  register )
    token=`curl -X POST $url/user -d "{\"email\":\"$2\", \"password\":\"$3\"}"` 
    echo "$token" > .gdf3.token 
    ;;
  ls )
    curl $url/bookmarks -H "Auth: $token"
    ;;
  get )
    curl "$url/bookmarks/$2" -H "Auth: $token"
    ;;
  * )
    usage 
    ;;
esac

