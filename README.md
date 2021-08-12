# TURNS

`turns` is a simple turn server using pion/turn lib. Some codes are from pion/turn/example/turn-server/simple/main.go, but added config file and a http api to add user/password.

# Build

``` shell
// clone this repo and cd turns, then just run
go build
```

# Usage

Repo provides a systemd service file, you can change it by necessary and put it in /lib/systemd/system/, then you can use systemd to use `turns`.

Or just running `./turns` is ok if you have a config file named `turns.json` in current directory.

About this config file, all fields are as follows:

- `api` is http server listened addr
- `port` is turn server listened UDP addr
- `public` shoule be your real public IP running `turns`
- `realm` is a key to encrypt password
- `auth` is a user:password table, of course you can add a pair of user:password dynamicly by `http://api/user`

Repo also has a example about this config file.


# License

MIT License - Since I used some code in pion/turn

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
