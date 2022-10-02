## Setting Up (Apple M1/M2 Macbook)

### Terminology
- MPS: Metal Performance Shaders


### Articles
- 2022/09/13 [Locally sourced: you can now run Stable Diffusion AI art generation on your M1 Mac](https://cdm.link/2022/09/run-stable-diffusion-ai-art-on-your-m1-mac/)
- 2022/06/06 [Deploying Transformers on the Apple Neural Engine](https://machinelearning.apple.com/research/neural-engine-transformers)
- 2022/05/18 [Running PyTorch on the M1 GPU](https://sebastianraschka.com/blog/2022/pytorch-m1-gpu.html)
- 2022/05/18 [Introducing Accelerated PyTorch Training on Mac](https://pytorch.org/blog/introducing-accelerated-pytorch-training-on-mac/)
- 2022/05/06 [New Release: Anaconda Distribution Now Supporting M1](https://www.anaconda.com/blog/new-release-anaconda-distribution-now-supporting-m1)
- 2022/04/22 [Deep Learning on the M1 Pro with Apple Silicon](https://wandb.ai/tcapelle/apple_m1_pro/reports/Deep-Learning-on-the-M1-Pro-with-Apple-Silicon---VmlldzoxMjQ0NjY3)



### Information
- [Anaconda](https://www.anaconda.com/)
	- [conda](https://docs.conda.io/projects/conda/en/latest/#)
	- [miniconda](https://docs.conda.io/en/latest/miniconda.html)
- Apple [Metal](https://developer.apple.com/metal/)
- PyTorch: [MPS Backend](https://pytorch.org/docs/master/notes/mps.html)
	- [Accelerated PyTorch Training on Mac](https://huggingface.co/docs/accelerate/usage_guides/mps)


### Open Source
- [pyenv/pyenv](https://github.com/pyenv/pyenv) - Simple Python version management
- [pytorch/pytorch](https://github.com/pytorch/pytorch) - Tensors and Dynamic neural networks in Python with strong GPU acceleration


### Install for Run Environments (PyTorch, Tensorflow, YOLOv7, OpenVINO, ONNX)
```sh
$ conda create --name=base python==3.10 
$ conda env list
$ coda activate
(base) $ pip install -U pip pandas numpy scipy h5py scikit-learn matplotlib tqdm jupyterlab
```

```sh
$ conda create -n torch python==3.9 --yes
$ coda activate torch
(base) $ pip install -U pip torch torchvision torchaudio
```

```sh
$ conda create -n tflow python==3.9 --yes
$ coda activate tflow
(base) $ pip install -U pip tensorflow-macos tensorflow-metal
```

```sh
$ conda create -n yolov7 python==3.9 --yes
$ coda activate yolov7
(base) $ pip install -U pip yolov5 yolov7
```

```sh
$ conda create -n ovino python==3.9 --yes
$ coda activate ovino
(base) $ pip install -U pip opencv-python openvino-arm
```

```sh
$ conda create -n onnx python==3.9 --yes
$ coda activate onnx
(base) $ pip install -U pip onnxruntime
```