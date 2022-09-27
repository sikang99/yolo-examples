import onnxruntime
import numpy as np

IN_PATH = ...      # path to your .onnx model (OUT_PATH if you used the previous example)
src_out = ... # output of your source framework, if this is a PyTorch tensor convert it to a numpy array using to_numpy
ort_session = onnxruntime.InferenceSession(IN_PATH)
def to_numpy(tensor):
    return tensor.detach().cpu().numpy() if tensor.requires_grad else tensor.cpu().numpy()
# compute ONNX Runtime output prediction
ort_inputs = {ort_session.get_inputs()[0].name: to_numpy(x)}
ort_outs = ort_session.run(None, ort_inputs)
# compare ONNX Runtime and PyTorch results
np.testing.assert_allclose(src_out, ort_outs[0], rtol=1e-03, atol=1e-05)