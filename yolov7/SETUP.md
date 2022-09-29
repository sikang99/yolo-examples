## Setting Up



### Articles
- 2022/05/18 [Running PyTorch on the M1 GPU](https://sebastianraschka.com/blog/2022/pytorch-m1-gpu.html)
- 2022/05/18 [Introducing Accelerated PyTorch Training on Mac](https://pytorch.org/blog/introducing-accelerated-pytorch-training-on-mac/)
- 2022/05/06 [New Release: Anaconda Distribution Now Supporting M1](https://www.anaconda.com/blog/new-release-anaconda-distribution-now-supporting-m1)



### Information
- [Anaconda](https://www.anaconda.com/)
	- [conda](https://docs.conda.io/projects/conda/en/latest/#)
	- [miniconda](https://docs.conda.io/en/latest/miniconda.html)
- Apple [Metal](https://developer.apple.com/metal/)
- [Accelerated PyTorch Training on Mac](https://huggingface.co/docs/accelerate/usage_guides/mps)



### Open Source
- [pyenv/pyenv](https://github.com/pyenv/pyenv) - Simple Python version management
- [pytorch/pytorch](https://github.com/pytorch/pytorch) - Tensors and Dynamic neural networks in Python with strong GPU acceleration


### Scripts
```sh
$ conda create -n tf python=3.9 --yes
$ conda activate tf
(tf) $ pip install tensorflow-macos
(tf) $ pip install tensotflow-metal
(tf) $ python
>>> import tensorflow as tf
>>> print(tf.__version__)
2.10.0
(tf) $ conda deactivate
```
