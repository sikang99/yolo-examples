#
# Makefile for cast
#
PROG=pang-text

usage:
	@echo "usage: make [edit|download|build|run]"

download d:
	wget https://github.com/opencv/opencv_3rdparty/raw/dnn_samples_face_detector_20170830/res10_300x300_ssd_iter_140000.caffemodel
	wget https://raw.githubusercontent.com/opencv/opencv/master/samples/dnn/face_detector/deploy.prototxt

edit e:
	vi $(PROG).go

build b:
	go build -o $(PROG) *.go

run r:
	@echo "> make (run) [Local|Remote|Media|Cctv]"

kill k:
	pkill $(PROG)

set-env se:
	source ../gocv/customenv

LOCAL=localhost
REMOTE=cobot.center
CHANNEL=c40hp6epjh65aeq6ne50

run-local rl:
	./$(PROG) -device=0 -scale=0.5 -url="wss://$(LOCAL):8277/pang/ws/pub?channel=$(CHANNEL)" &
run-remote rr:
	./$(PROG) -device=0 -scale=0.5 -url="wss://$(REMOTE):8277/pang/ws/pub?channel=$(CHANNEL)" &

run-local-media rlm:
	./$(PROG) -media=sample.mp4 &
run-local-cctv rlc:
	./$(PROG) -media="rtsp://admin:admin1234@221.155.2.140:555/cam/realmonitor?channel=1&subtype=1&unicast=true&proto=Onvif"
rlc2:
	./$(PROG) -scale=0.5 -media=rtsp://64.187.201.16/axis-media/media.amp

run-help rh:
	./$(PROG) -h

run-version rv:
	./$(PROG) -version
	pkg-config --cflags  -- opencv4
	brew info opencv


t1:
	go run sample.go -input images/bottle.jpg -model frozen_east_text_detection.pb -ocr CRNN_VGG_BiLSTM_CTC.onnx


open o:
	open http://localhost:8080/

media-info mi:
	packager input=sample.mp4 --dump_stream_info

clean:
	rm -f $(PROG)

build-run br:
	@make build
	@make run-local

git-update gu:
	cd ../..; make gu

