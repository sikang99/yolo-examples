import torch

# path to your .onnx model (OUT_PATH if you used the previous example)
IN_PATH = "../assets/models/yolov5x.onnx"
onnx_model = torch.onnx.load(IN_PATH)
torch.onnx.checker.check_model(onnx_model)
