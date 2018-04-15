ethereum-brute-force
---

A small goscript to check random privatekeys for a balance and record if a used ethereum address is found.

### Usage

```
go run main.go [--update=<boolean>] [--provider=<node location>]
```

Make sure to go to the releases section and download the list of addresses, otherwise you'll need to build it from scratch.

use the optional `-provider` in conjunction with the `-update=true` flag to overwrite the default node location `http://localhost:8545` to update the list of addresses. This doesn't work against infura as I didn't implement any retry mechanism in the case of failure.

When running with the provided `blances.csv` (located in the releases section), it'll take a little bit to get started up as there are ~27 million addresses in that list.

### About

This is more of a proof of concept - the actual probability of finding a private key in use is about 30000000 / 115792089237316195423570985008687907853269984665640564039457584007913129639936. Even at millions of privatekeys a second, it will be a very long time before you will likely find anything. Though, I'll admit that there is some fun in leaving it running.

I didn't scan uncle blocks for addresses or look into receipts for token transactions.

May number gods smile upon you.

### License

Copyright 2017 Andrew Donley

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
