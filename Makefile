# cannot use relative path in GOROOT, otherwise 6g not found. For example,
#   export GOROOT=../go  (=> 6g not found)
# it is also not allowed to use relative path in GOPATH
export GOROOT=$(realpath ../go)
export GOPATH=$(realpath .)
export PATH := $(GOROOT)/bin:$(GOPATH)/bin:$(PATH)


default: kill
	gopherjs serve -w --http ":8000" &
	gopherjs build src/snake.go -w -o src/snake.js &

kill:
	ps aux | grep gopherjs
	-killall gopherjs
	ps aux | grep gopherjs

install:
	go get -u github.com/gopherjs/gopherjs
	go get -u honnef.co/go/js/dom

clean:
	-rm src/snake.js
	-rm src/snake.js.map
	-rm -rf bin pkg
	-rm -rf src/github.com
	-rm -rf src/golang.org
	-rm -rf src/gopkg.in
	-rm -rf src/honnef.co
