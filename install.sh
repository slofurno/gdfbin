#!/bin/sh

basedir=$(pwd)
echo "#!/bin/sh" > gdf3
echo "" >> gdf3
echo "basedir=$basedir" >> gdf3
cat $basedir/gdf3.sh >> gdf3
chmod +x $basedir/gdf3
ln -s $basedir/gdf3 /usr/local/bin/gdf3

