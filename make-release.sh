#!/bin/sh -e

programs="g4ping g4setclock g4update"
go_ldflags="-s -w"
target_archs="arm 386"
go_tags_arm=""
go_tags_386="nofilter"
dest=binaries
package=$(basename $(pwd))

for arch in $target_archs; do
    mkdir -p $dest/$arch
    echo "Building binaries for $arch architecture"
    varname=go_tags_${arch}
    go_tags="$(eval echo \$${varname})"
    for prog in $programs; do
	(cd cmd/$prog && \
	 GOARCH=$arch go build -tags "$go_tags" -ldflags "$go_ldflags" && \
	 mv -v $prog ../../$dest/$arch/)
    done
    tarball=${package}-${arch}.tar.xz
    echo Building $tarball
    tar --create --file $dest/$tarball --xz --verbose --directory $dest/$arch .
done

ls -l $dest/*.xz
