#!/bin/bash
set -e
mkdir -p data
cd data
wget -N https://download.geonames.org/export/dump/cities500.zip
unzip -o cities500.zip 