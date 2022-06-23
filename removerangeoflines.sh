#!/bin/bash

# Used for removing redundant lines fom big CSV files created from .xlsx files

export START=50011

export STOP=1015052

export FILENAME=blabla3.csv

export OUT=out.csv

sed "${START},${STOP}d" $FILENAME > $OUT 
