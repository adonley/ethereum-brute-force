#!/bin/bash

main () {
	# Make addresses directory if it doesn't exist
	if [ ! -d "addresses" ]; then
		mkdir addresses;
	fi

	# Delete the addresses if it exists
	if [ -f "addresses/addresses.csv" ]; then
	  rm addresses/addresses.csv;
	fi

	touch addresses/addresses.csv;

	# Download all the files that make up the addresses
	for i in {0..9}; do 
		# Only download if doesn't exist
		if [ ! -f "addresses/part$i.gz" ]; then
			echo "Downloading part$i.gz";
			wget https://storage.googleapis.com/ethereum-addresses/addresses00000000000$i -O addresses/part$i.gz;
		fi

		# Unzip if doesn't exist
		if [ ! -f "addresses/part$i.csv" ]; then
			echo "Decompressing part$i.csv";
			gunzip -c addresses/part$i.gz > addresses/part$i.csv;
    fi

    tail -n +2 "addresses/part$i.csv" >> addresses/addresses.csv;
    rm "addresses/part$i.csv" "addresses/part$i.gz";
	done

  local LINES;
  LINES=$(wc -l < addresses/addresses.csv)
	printf "Done grabbing and parsing address files, %s addresses loaded.\n" "$LINES";
}

main;
