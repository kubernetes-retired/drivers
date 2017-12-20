# CSI glusterblock driver

## Usage:

### Start glusterblock driver
```
$ sudo ../_output/glusterblockdriver --endpoint tcp://127.0.0.1:10000 --nodeid CSINode
```

### Test using csc
Get ```csc``` tool from https://github.com/chakri-nelluri/gocsi/tree/master/csc

#### Get plugin info
```
$ csc identity plugin-info --endpoint tcp://127.0.0.1:10000
"glusterblock"	"0.1.0"
```

### Get supported versions
```
$ csc identity supported-versions --endpoint tcp://127.0.0.1:10000
0.1.0
```

#### NodePublish a volume
```
$ export glusterblock_TARGET="glusterblock Target Server IP (Ex: 10.10.10.10)"
$ export IQN="Target IQN"
$ csc node publish --endpoint tcp://127.0.0.1:10000 --target-path /mnt/glusterblock --attrib targetPortal=$glusterblock_TARGET --attrib iqn=$IQN --attrib lun=<lun-id> glusterblocktestvol
glusterblocktestvol
```

#### NodeUnpublish a volume
```
$ csc node unpublish --endpoint tcp://127.0.0.1:10000 --target-path /mnt/glusterblock glusterblocktestvol
glusterblocktestvol
```

#### Get NodeID
```
$ csc node get-id --endpoint tcp://127.0.0.1:10000
CSINode
```
