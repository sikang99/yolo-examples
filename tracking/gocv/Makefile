#
# Makefile for gocv tracking
#
PROG=tracker
usage:
	@echo "usage: make [edit|build|run]"

edit e:
	vi main.go

build b:
	go build -o $(PROG) *.go

run r:
	@echo "> make (run) [kcf|csrt]"

run-help rh:
	./$(PROG) -h
run-kcf rk:
	./$(PROG) -tracker=kcf
run-csrt rc:
	./$(PROG) -tracker=csrt

kill:
	pkill $(PROG)

clean:
	rm -f $(PROG)

git-update gu:
	cd ..; make gu


