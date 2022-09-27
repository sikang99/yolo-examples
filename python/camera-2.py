import cv2

cap = cv2.VideoCapture(0)

cap.set(cv2.CAP_PROP_FRAME_WIDTH, 1280)
cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 720)

width = cap.get(cv2.CAP_PROP_FRAME_WIDTH) 
height = cap.get(cv2.CAP_PROP_FRAME_HEIGHT) 

print("Size:", int(width), "x", int(height))

while True:
    ret, img = cap.read()

    cv2.imshow('Camera', img)
    if cv2.waitKey(1) == ord('q'):
        break
        
cap.release() 
