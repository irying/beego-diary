#!/bin/bash
bee pack -be GOOS=linux -p src/api -o build 
rm -rf app/*
tar -C app/ -xzf build/api.tar.gz
rm -rf build/*
