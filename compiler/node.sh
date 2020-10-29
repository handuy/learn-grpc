#!/bin/sh
echo "$@"
echo "$@" >> index.js
cat index.js
node index.js
