import torch 
SHAPE = ...                                       # your batch shape
OUT_PATH = ...                                    # output path
x = torch.randn(SHAPE)
with torch.no_grad():
    if isinstance(model, torch.nn.DataParallel):  # extract the module from dataparallel models
        model = model.module
    model.cpu()
    model.eval()                                  # the converter works best on models stored on the CPU
    torch.onnx.export(model,                      # model being run
                      x,                          # model input (or a tuple for multiple inputs)
                      OUT_PATH,                   # where to save the model (can be a file or file-like object)
                      export_params=True,         # store the trained parameter weights inside the model
                      opset_version=11)           # it's best to specify the opset version. At time of writing 11 was the latest release