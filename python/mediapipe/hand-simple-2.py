import cv2
import mediapipe
 
drawingModule = mediapipe.solutions.drawing_utils
handsModule = mediapipe.solutions.hands
 
capture = cv2.VideoCapture(0)
 
frameWidth = capture.get(cv2.CAP_PROP_FRAME_WIDTH)
frameHeight = capture.get(cv2.CAP_PROP_FRAME_HEIGHT)
 
with handsModule.Hands(
    static_image_mode=False, 
    min_detection_confidence=0.7, 
    min_tracking_confidence=0.7, 
    max_num_hands=2) as hands:
 
    while (True):
        ret, frame = capture.read()
 
        results = hands.process(cv2.cvtColor(frame, cv2.COLOR_BGR2RGB))
 
        if results.multi_hand_landmarks != None:
            for handLandmarks in results.multi_hand_landmarks:
                for point in handsModule.HandLandmark:
 
                    normalizedLandmark = handLandmarks.landmark[point]
                    pixelCoordinatesLandmark = drawingModule._normalized_to_pixel_coordinates(normalizedLandmark.x, normalizedLandmark.y, frameWidth, frameHeight)
 
                    cv2.circle(frame, pixelCoordinatesLandmark, 5, (0, 255, 0), -1)
 
        cv2.imshow('Test hand #2', frame)
        if cv2.waitKey(1) == 27:
            break
 
cv2.destroyAllWindows()
capture.release()
