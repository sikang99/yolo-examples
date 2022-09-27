#
# Makefile for model, config data
#
usage:
	@echo "usage: make [yolov3]"

model:
	open https://sparsezoo.neuralmagic.com/

yolov3:
	wget https://pjreddie.com/media/files/yolov3.weights
	wget https://pjreddie.com/media/files/yolov3-tiny.weights

yolov4:
	wget https://raw.githubusercontent.com/AlexeyAB/darknet/master/cfg/yolov4-tiny.cfg
	wget https://github.com/AlexeyAB/darknet/releases/download/darknet_yolo_v4_pre/yolov4-tiny.weights

yolov5:
	open https://github.com/ultralytics/yolov5/releases

yolov7:
	wget 
	

install i:
	pip install yolov5

convert:
	yolov5 export --weights yolov5x.pt --include torchscript,onnx,coreml,pb,tfjs

run-detect rd:
	python3 detect.py

git-update gu:
	cd ../..; make gu

