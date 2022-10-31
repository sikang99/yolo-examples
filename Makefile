#
# Makefile
#
usage:
	@echo "usage: make [download|git]"

download d:
	@echo "> make (download) [yolov3|yolov5]"

download-yolov3 d3:
	mkdir -p ./data
	wget https://pjreddie.com/media/files/yolov3.weights -O ./data/yolov3.weights
	wget https://github.com/pjreddie/darknet/blob/master/cfg/yolov3.cfg?raw=true -O ./data/yolov3.cfg
	wget https://github.com/pjreddie/darknet/blob/master/data/coco.names?raw=true -O ./data/coco.names

download-yolov5 d5:
	mkdir -p ./data
	wget https://github.com/doleron/yolov5-opencv-cpp-python/raw/main/config_files/yolov5s.onnx -O ./data/yolov5s.onnx
	wget https://github.com/doleron/yolov5-opencv-cpp-python/raw/main/config_files/yolov5n.onnx -O ./data/yolov5n.onnx
	wget https://github.com/doleron/yolov5-opencv-cpp-python/raw/main/config_files/classes.txt -O ./data/yolov5.names


#--------------------------------------------------------------------------------
USER=stoney
BUILD=0.0.3.2
git g:
	@echo "make (git:g) [update|store]"
git-reset gr:
	git reset --soft HEAD~1
git-update gu:
	git add .
	git commit -a -m "$(BUILD),$(USER)"
	git push
git-store gs:
	git config credential.helper store
#--------------------------------------------------------------------------------
	
