import cv2
import numpy as np
import mediapipe as mp

v_cap = cv2.VideoCapture(0) #カメラのIDを選ぶ。映らない場合は番号を変える。

with mp.solutions.face_mesh.FaceMesh( #mesh化の設定をしてその処理名をface_meshとする
    max_num_faces=1,refine_landmarks=True,min_detection_confidence=0.5,min_tracking_confidence=0.5) as face_mesh:
    # 顔の数, ランドマークのリファイン?, ランドマーク検出成功判定の閾値, ランドマークトラッキング成功判定の閾値
  while v_cap.isOpened():
    success, image = v_cap.read() #キャプチャが成功していたら画像データとしてimageに取り込む
    results = face_mesh.process(image)#メッシュ化計算の結果がresutsに入る
    # 表示画面の背景を黒で塗りつぶす
    image_blank = np.zeros((int(v_cap.get(cv2.CAP_PROP_FRAME_HEIGHT)), int(v_cap.get(cv2.CAP_PROP_FRAME_WIDTH)), 3))
    image = image_blank

    if results.multi_face_landmarks:
      for face_landmarks in results.multi_face_landmarks:
        mp.solutions.drawing_utils.draw_landmarks(
            image=image, landmark_list=face_landmarks, connections=mp.solutions.face_mesh.FACEMESH_TESSELATION,
            # 画像データ, ランドマークのリスト, メッシュの繋ぎ方?(テッセレーション)
            landmark_drawing_spec=None,connection_drawing_spec=mp.solutions.drawing_styles.get_default_face_mesh_tesselation_style())
            # わからぬ, わからぬ
    cv2.imshow('MediaPipe Face Mesh', cv2.flip(image, 1))#imageを鏡像で表示
    if cv2.waitKey(5) & 0xFF == 27: #ESCキーが押されたら終わる
      break

v_cap.release()
