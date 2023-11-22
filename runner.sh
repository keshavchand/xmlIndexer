#! /bin/bash

go run *.go -method $1 -file SwissProt.xml -index index.gob
