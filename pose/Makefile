
usage:
	@echo "usage: make [run]"

model m:
	wget http://posefs1.perception.cs.cmu.edu/OpenPose/models/pose/coco/pose_iter_440000.caffemodel

run r:
	go run main.go

run-tiny rt:
	go run main.go -model=tiny

run-spp rs:
	go run main.go -model=spp

run-help rh:
	go run main.go -h

git-update gu:
	cd ..; make gu


