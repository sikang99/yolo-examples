import sys
import torch
print(f"Python version: {sys.version}, {sys.version_info} ")
print("Pytorch version: ", torch.__version__)
print("MPS support: ", torch.backends.mps.is_available(), torch.backends.mps.is_built())
