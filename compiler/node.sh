#!/bin/sh
echo Your container args are: "$@"
echo "$@" >> index.js
cat index.js
node index.js
