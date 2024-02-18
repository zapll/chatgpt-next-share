#!/bin/bash

dir=service

rm -rf $dir

mkdir -vp $dir/admin

run_admin=$(cat <<- EOF
#!/bin/sh
set -e
echo "admin start at http://127.0.0.1:3001"
exec bun run index.js 2>&1
EOF
)

echo "$run_admin" > $dir/admin/run

chmod +x $dir/admin/run

mkdir -vp $dir/share

run_admin=$(cat <<- EOF
#!/bin/sh
set -e
echo "share start at http://127.0.0.1:3000"
exec ./share 2>&1
EOF
)

echo "$run_admin" > $dir/share/run

chmod +x $dir/share/run

echo "build share"
cd share
CGO_ENABLED=1 go build -ldflags "-s -w" -o ../$dir/share .
cd -

echo "build admin"
cp -rf admin/ui $dir/admin
cd admin
bun install && bun run build
cd -

echo "build ok"