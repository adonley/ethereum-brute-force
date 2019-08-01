ethereum-brute-force
---

A small go program to check random privatekeys for a balance and record if a used ethereum address is found.

### Setup
Run the get addresses script to download all the ethereum addresses (this will take a while):
```
./get-addresses.sh
```

### Usage
With golang installed;
```
go run main.go
```

The docker-compose is a work in progress (it's loading addresses too slow at the moment);
```
docker-compose up
```

### About

This is a proof of concept - the actual probability of finding a private key in use is about 90867014 / 115792089237316195423570985008687907853269984665640564039457584007913129639936. Even at millions of privatekeys a second, it will be a very long time before you will likely find anything. Though, I'll admit that there is some fun in leaving it running.

Some things have changed since when I first wrote this. Google now has an open dataset containing the whole Ethereum block chain. This
is what I used to get the addresses dataset. You can find the query I used in the file `bigquery` included in this project. The size
of the addresses file is at least double. This will take 100% of the memory on my 16GB RAM linux machine. It's probably time to work toward
a distributed solution.

May number gods smile upon you.

# License

Copyright 2017 Andrew Donley

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
