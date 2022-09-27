#
# Makefile for gocv tracking
#
PROG=circle
usage:
	@echo "usage: make [edit|build|run]"

edit e:
	vi main.go

build b:
	go build -o $(PROG) *.go

run r:
	go run main.go

kill:
	pkill $(PROG)

clean:
	rm -f $(PROG)

git-update gu:
	cd ..; make gu


