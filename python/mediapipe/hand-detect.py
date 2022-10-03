import cv2
import mediapipe as mp
 
mp_drawing = mp.solutions.drawing_utils
mp_hands = mp.solutions.hands
 
cap = cv2.VideoCapture(0)
 
with mp_hands.Hands(
    max_num_hands=2,
    min_detection_confidence=0.5,
    min_tracking_confidence=0.5) as hands:
 
    while cap.isOpened():
        success, image = cap.read()
        if not success:
            continue

        image = cv2.cvtColor(cv2.flip(image, 1), cv2.COLOR_BGR2RGB)
        results = hands.process(image)
        image = cv2.cvtColor(image, cv2.COLOR_RGB2BGR)
 
        if results.multi_hand_landmarks:
            for hand_landmarks in results.multi_hand_landmarks:
                finger1 = int(hand_landmarks.landmark[4].x * 100 )
                finger2 = int(hand_landmarks.landmark[8].x * 100 )
                dist = abs(finger1 - finger2)

                cv2.putText(
                    image, text='f1=%d f2=%d dist=%d ' % (finger1,finger2,dist), org=(10, 30),
                    fontFace=cv2.FONT_HERSHEY_SIMPLEX, fontScale=1,
                    color=255, thickness=3)
 
                mp_drawing.draw_landmarks(
                    image, hand_landmarks, mp_hands.HAND_CONNECTIONS)
 
        cv2.imshow('MediaPipe: Hands', image)
        if cv2.waitKey(1) == ord('q'):
            break
 
cap.release()
