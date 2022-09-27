import cv2

cap = cv2.VideoCapture(0)

while True:
    ret, img = cap.read()

    cv2.imshow('Video', img)
    
    if cv2.waitKey(10) == ord('q'):
        break

cap.release()
