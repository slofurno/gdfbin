#!/bin/sh
basedir=$(pwd)
url=${1-https://gdf3.com}

cat << EOM > ./gdf3
#!/bin/sh

basedir="$basedir"
url="$url"
EOM

cat << 'EOM' >> ./gdf3

usage (){
  echo "gdf3 register <email> <password>"
  echo "gdf3 login <email> <password>"
  echo "cat <file> | gdf3 mark <name>"
  echo "cat <file> | gdf3 paste"
  echo "gdf3 ls"
  echo "gdf3 get <name>"
  echo "gdf3 rm <name>"
  echo 'gdf3 cp <name>'
  echo 'gdf3 log <name>'
}

token=`cat $basedir/.gdf3.token 2> /dev/null`

case $1 in
  mark )
    name=`curl -s --data-binary @- $url | sed "s|$url/||"`
    curl -X POST "$url/bookmarks/$name/$2" -H "Auth: $token"
    ;;
  paste )
    curl -s --data-binary @- $url
    ;;
  register )
    token=`curl -s -X POST $url/user -d "{\"email\":\"$2\", \"password\":\"$3\"}"`
    echo "$token"
    echo "$token" > $basedir/.gdf3.token
    ;;
  login )
    echo "$url/login"
    token=`curl -s -X POST $url/login -d "{\"email\":\"$2\", \"password\":\"$3\"}"`
    echo "$token" > $basedir/.gdf3.token
    ;;
  ls )
    curl -s $url/bookmarks -H "Auth: $token"
    ;;
  lsa )
    echo "not implemented )"
    ;;
  rm )
    curl -X DELETE "$url/bookmarks/$2" -H "Auth: $token"
    ;;
  get )
    curl -s "$url/bookmarks/$2" -H "Auth: $token"
    ;;
  cp )
    [ -e "$2" ] && echo "file already exists" || curl -s "$url/bookmarks/$2" -H "Auth: $token" >> $2
    ;;
  log )
    curl -s "$url/bookmarks/$2/history" -H "Auth: $token"
    ;;
  * )
    usage
    ;;
esac

EOM

chmod +x $basedir/gdf3
ln -sf $basedir/gdf3 /usr/local/bin/gdf3
