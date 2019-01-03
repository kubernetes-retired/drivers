# CSI Image driver

## Usage:

### Build imageplugin
```
$ make image
```

### Start Image driver
```
$ sudo ./_output/imageplugin --endpoint tcp://127.0.0.1:10000 --nodeid CSINode -v=5
```

### Test using csc
Get ```csc``` tool from https://github.com/rexray/gocsi/tree/master/csc

### Mount the image
$ csc -e tcp://127.0.0.1:10000 node publish abcdefg --attrib image=kfox1111/misc:test --cap MULTI_NODE_MULTI_WRITER,block --target-path /tmp/csi

### Unmount the image
$ csc -e tcp://127.0.0.1:10000 node unpublish abcdefg --target-path /tmp/csi
