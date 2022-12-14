#!/usr/bin/env python
# -*- coding: utf-8 -*-
import copy
import argparse

import cv2 as cv
import numpy as np
import mediapipe as mp

from collections import deque

class CvFpsCalc(object):
    def __init__(self, buffer_len=1):
        self._start_tick = cv.getTickCount()
        self._freq = 1000.0 / cv.getTickFrequency()
        self._difftimes = deque(maxlen=buffer_len)

    def get(self):
        current_tick = cv.getTickCount()
        different_time = (current_tick - self._start_tick) * self._freq
        self._start_tick = current_tick

        self._difftimes.append(different_time)

        fps = 1000.0 / (sum(self._difftimes) / len(self._difftimes))
        fps_rounded = round(fps, 2)

        return fps_rounded

# from utils import CvFpsCalc


def get_args():
    parser = argparse.ArgumentParser()

    parser.add_argument("--device", type=int, default=0)
    parser.add_argument("--width", help='cap width', type=int, default=1280)
    parser.add_argument("--height", help='cap height', type=int, default=720)

    parser.add_argument("--f_min_detection_confidence",
                        help='face mesh min_detection_confidence',
                        type=float,
                        default=0.5)
    parser.add_argument("--f_min_tracking_confidence",
                        help='face mesh min_tracking_confidence',
                        type=int,
                        default=0.5)
    parser.add_argument("--h_min_detection_confidence",
                        help='hands min_detection_confidence',
                        type=float,
                        default=0.7)
    parser.add_argument("--h_min_tracking_confidence",
                        help='hands min_tracking_confidence',
                        type=int,
                        default=0.5)
    parser.add_argument("--p_min_detection_confidence",
                        help='pose min_detection_confidence',
                        type=float,
                        default=0.5)
    parser.add_argument("--p_min_tracking_confidence",
                        help='pose min_tracking_confidence',
                        type=int,
                        default=0.5)

    parser.add_argument('--use_brect', action='store_true')

    args = parser.parse_args()

    return args


def main():
    args = get_args()

    cap_device = args.device
    cap_width = args.width
    cap_height = args.height

    f_min_detection_confidence = args.f_min_detection_confidence
    f_min_tracking_confidence = args.f_min_tracking_confidence
    h_min_detection_confidence = args.h_min_detection_confidence
    h_min_tracking_confidence = args.h_min_tracking_confidence
    p_min_detection_confidence = args.p_min_detection_confidence
    p_min_tracking_confidence = args.p_min_tracking_confidence

    use_brect = args.use_brect

    # ??????????????? ###############################################################
    cap = cv.VideoCapture(cap_device)
    cap.set(cv.CAP_PROP_FRAME_WIDTH, cap_width)
    cap.set(cv.CAP_PROP_FRAME_HEIGHT, cap_height)

    # ?????????????????? #############################################################
    mp_face_mesh = mp.solutions.face_mesh
    face_mesh = mp_face_mesh.FaceMesh(
        min_detection_confidence=f_min_detection_confidence,
        min_tracking_confidence=f_min_tracking_confidence,
    )
    mp_hands = mp.solutions.hands
    hands = mp_hands.Hands(
        min_detection_confidence=h_min_detection_confidence,
        min_tracking_confidence=h_min_tracking_confidence,
    )
    mp_pose = mp.solutions.pose
    pose = mp_pose.Pose(
        min_detection_confidence=p_min_detection_confidence,
        min_tracking_confidence=p_min_tracking_confidence,
    )

    # FPS????????????????????? ########################################################
    cvFpsCalc = CvFpsCalc(buffer_len=10)

    while True:
        display_fps = cvFpsCalc.get()

        # ???????????????????????? #####################################################
        ret, image = cap.read()
        if not ret:
            break
        image = cv.flip(image, 1)  # ???????????????
        debug_image = copy.deepcopy(image)

        # ???????????? #############################################################
        image = cv.cvtColor(image, cv.COLOR_BGR2RGB)

        image.flags.writeable = False
        hands_results = hands.process(image)
        face_results = face_mesh.process(image)
        pose_results = pose.process(image)
        image.flags.writeable = True

        # Face Mesh ###########################################################
        if face_results.multi_face_landmarks is not None:
            for face_landmarks in face_results.multi_face_landmarks:
                # ?????????????????????
                brect = calc_bounding_rect(debug_image, face_landmarks)
                # ??????
                debug_image = draw_face_landmarks(debug_image, face_landmarks)
                debug_image = draw_bounding_rect(use_brect, debug_image, brect)

        # Pose ###############################################################
        if pose_results.pose_landmarks is not None:
            # ?????????????????????
            brect = calc_bounding_rect(debug_image,
                                       pose_results.pose_landmarks)
            # ??????
            debug_image = draw_pose_landmarks(debug_image,
                                              pose_results.pose_landmarks)
            debug_image = draw_bounding_rect(use_brect, debug_image, brect)

        # Hands ###############################################################
        if hands_results.multi_hand_landmarks is not None:
            for hand_landmarks, handedness in zip(
                    hands_results.multi_hand_landmarks,
                    hands_results.multi_handedness):
                # ?????????????????????
                cx, cy = calc_palm_moment(debug_image, hand_landmarks)
                # ?????????????????????
                brect = calc_bounding_rect(debug_image, hand_landmarks)
                # ??????
                debug_image = draw_hands_landmarks(debug_image, cx, cy, hand_landmarks, handedness)
                debug_image = draw_bounding_rect(use_brect, debug_image, brect)

        cv.putText(debug_image, "FPS:" + str(display_fps), (10, 30),
                   cv.FONT_HERSHEY_SIMPLEX, 1.0, (0, 255, 0), 2, cv.LINE_AA)
        cv.imshow('MediaPipe Demo', debug_image)

        key = cv.waitKey(1)
        if key == 27:  # ESC
            break

    cap.release()
    cv.destroyAllWindows()


def calc_palm_moment(image, landmarks):
    image_width, image_height = image.shape[1], image.shape[0]

    palm_array = np.empty((0, 2), int)

    for index, landmark in enumerate(landmarks.landmark):
        landmark_x = min(int(landmark.x * image_width), image_width - 1)
        landmark_y = min(int(landmark.y * image_height), image_height - 1)

        landmark_point = [np.array((landmark_x, landmark_y))]

        if index == 0:  # ??????1
            palm_array = np.append(palm_array, landmark_point, axis=0)
        if index == 1:  # ??????2
            palm_array = np.append(palm_array, landmark_point, axis=0)
        if index == 5:  # ?????????????????????
            palm_array = np.append(palm_array, landmark_point, axis=0)
        if index == 9:  # ??????????????????
            palm_array = np.append(palm_array, landmark_point, axis=0)
        if index == 13:  # ??????????????????
            palm_array = np.append(palm_array, landmark_point, axis=0)
        if index == 17:  # ??????????????????
            palm_array = np.append(palm_array, landmark_point, axis=0)
    M = cv.moments(palm_array)
    cx, cy = 0, 0
    if M['m00'] != 0:
        cx = int(M['m10'] / M['m00'])
        cy = int(M['m01'] / M['m00'])

    return cx, cy


def calc_bounding_rect(image, landmarks):
    image_width, image_height = image.shape[1], image.shape[0]

    landmark_array = np.empty((0, 2), int)

    for _, landmark in enumerate(landmarks.landmark):
        landmark_x = min(int(landmark.x * image_width), image_width - 1)
        landmark_y = min(int(landmark.y * image_height), image_height - 1)

        landmark_point = [np.array((landmark_x, landmark_y))]

        landmark_array = np.append(landmark_array, landmark_point, axis=0)

    x, y, w, h = cv.boundingRect(landmark_array)

    return [x, y, x + w, y + h]


def draw_hands_landmarks(image, cx, cy, landmarks, handedness):
    image_width, image_height = image.shape[1], image.shape[0]

    landmark_point = []

    # ??????????????????
    for index, landmark in enumerate(landmarks.landmark):
        if landmark.visibility < 0 or landmark.presence < 0:
            continue

        landmark_x = min(int(landmark.x * image_width), image_width - 1)
        landmark_y = min(int(landmark.y * image_height), image_height - 1)
        # landmark_z = landmark.z

        landmark_point.append((landmark_x, landmark_y))

        if index == 0:  # ??????1
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 1:  # ??????2
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 2:  # ??????????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 3:  # ????????????1??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 4:  # ???????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
            cv.circle(image, (landmark_x, landmark_y), 12, (0, 255, 0), 2)
        if index == 5:  # ?????????????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 6:  # ???????????????2??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 7:  # ???????????????1??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 8:  # ??????????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
            cv.circle(image, (landmark_x, landmark_y), 12, (0, 255, 0), 2)
        if index == 9:  # ??????????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 10:  # ????????????2??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 11:  # ????????????1??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 12:  # ???????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
            cv.circle(image, (landmark_x, landmark_y), 12, (0, 255, 0), 2)
        if index == 13:  # ??????????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 14:  # ????????????2??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 15:  # ????????????1??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 16:  # ???????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
            cv.circle(image, (landmark_x, landmark_y), 12, (0, 255, 0), 2)
        if index == 17:  # ??????????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 18:  # ????????????2??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 19:  # ????????????1??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 20:  # ???????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
            cv.circle(image, (landmark_x, landmark_y), 12, (0, 255, 0), 2)

    # ?????????
    if len(landmark_point) > 0:
        # ??????
        cv.line(image, landmark_point[2], landmark_point[3], (0, 255, 0), 2)
        cv.line(image, landmark_point[3], landmark_point[4], (0, 255, 0), 2)

        # ?????????
        cv.line(image, landmark_point[5], landmark_point[6], (0, 255, 0), 2)
        cv.line(image, landmark_point[6], landmark_point[7], (0, 255, 0), 2)
        cv.line(image, landmark_point[7], landmark_point[8], (0, 255, 0), 2)

        # ??????
        cv.line(image, landmark_point[9], landmark_point[10], (0, 255, 0), 2)
        cv.line(image, landmark_point[10], landmark_point[11], (0, 255, 0), 2)
        cv.line(image, landmark_point[11], landmark_point[12], (0, 255, 0), 2)

        # ??????
        cv.line(image, landmark_point[13], landmark_point[14], (0, 255, 0), 2)
        cv.line(image, landmark_point[14], landmark_point[15], (0, 255, 0), 2)
        cv.line(image, landmark_point[15], landmark_point[16], (0, 255, 0), 2)

        # ??????
        cv.line(image, landmark_point[17], landmark_point[18], (0, 255, 0), 2)
        cv.line(image, landmark_point[18], landmark_point[19], (0, 255, 0), 2)
        cv.line(image, landmark_point[19], landmark_point[20], (0, 255, 0), 2)

        # ?????????
        cv.line(image, landmark_point[0], landmark_point[1], (0, 255, 0), 2)
        cv.line(image, landmark_point[1], landmark_point[2], (0, 255, 0), 2)
        cv.line(image, landmark_point[2], landmark_point[5], (0, 255, 0), 2)
        cv.line(image, landmark_point[5], landmark_point[9], (0, 255, 0), 2)
        cv.line(image, landmark_point[9], landmark_point[13], (0, 255, 0), 2)
        cv.line(image, landmark_point[13], landmark_point[17], (0, 255, 0), 2)
        cv.line(image, landmark_point[17], landmark_point[0], (0, 255, 0), 2)

    # ?????? + ??????
    if len(landmark_point) > 0:
        # handedness.classification[0].index
        # handedness.classification[0].score

        cv.circle(image, (cx, cy), 12, (0, 255, 0), 2)
        cv.putText(image, handedness.classification[0].label[0],
                   (cx - 6, cy + 6), cv.FONT_HERSHEY_SIMPLEX, 0.6, (0, 255, 0),
                   2, cv.LINE_AA)  # label[0]:??????????????????

    return image


def draw_face_landmarks(image, landmarks):
    image_width, image_height = image.shape[1], image.shape[0]

    landmark_point = []

    for index, landmark in enumerate(landmarks.landmark):
        if landmark.visibility < 0 or landmark.presence < 0:
            continue

        landmark_x = min(int(landmark.x * image_width), image_width - 1)
        landmark_y = min(int(landmark.y * image_height), image_height - 1)

        landmark_point.append((landmark_x, landmark_y))

        cv.circle(image, (landmark_x, landmark_y), 1, (0, 255, 0), 1)

    if len(landmark_point) > 0:
        # ?????????https://github.com/tensorflow/tfjs-models/blob/master/facemesh/mesh_map.jpg

        # ?????????(55????????????46?????????)
        cv.line(image, landmark_point[55], landmark_point[65], (0, 255, 0), 2)
        cv.line(image, landmark_point[65], landmark_point[52], (0, 255, 0), 2)
        cv.line(image, landmark_point[52], landmark_point[53], (0, 255, 0), 2)
        cv.line(image, landmark_point[53], landmark_point[46], (0, 255, 0), 2)

        # ?????????(285????????????276?????????)
        cv.line(image, landmark_point[285], landmark_point[295], (0, 255, 0), 2)
        cv.line(image, landmark_point[295], landmark_point[282], (0, 255, 0), 2)
        cv.line(image, landmark_point[282], landmark_point[283], (0, 255, 0), 2)
        cv.line(image, landmark_point[283], landmark_point[276], (0, 255, 0), 2)

        # ?????? (133????????????246?????????)
        cv.line(image, landmark_point[133], landmark_point[173], (0, 255, 0), 2)
        cv.line(image, landmark_point[173], landmark_point[157], (0, 255, 0), 2)
        cv.line(image, landmark_point[157], landmark_point[158], (0, 255, 0), 2)
        cv.line(image, landmark_point[158], landmark_point[159], (0, 255, 0), 2)
        cv.line(image, landmark_point[159], landmark_point[160], (0, 255, 0), 2)
        cv.line(image, landmark_point[160], landmark_point[161], (0, 255, 0), 2)
        cv.line(image, landmark_point[161], landmark_point[246], (0, 255, 0), 2)

        cv.line(image, landmark_point[246], landmark_point[163], (0, 255, 0), 2)
        cv.line(image, landmark_point[163], landmark_point[144], (0, 255, 0), 2)
        cv.line(image, landmark_point[144], landmark_point[145], (0, 255, 0), 2)
        cv.line(image, landmark_point[145], landmark_point[153], (0, 255, 0), 2)
        cv.line(image, landmark_point[153], landmark_point[154], (0, 255, 0), 2)
        cv.line(image, landmark_point[154], landmark_point[155], (0, 255, 0), 2)
        cv.line(image, landmark_point[155], landmark_point[133], (0, 255, 0), 2)

        # ?????? (362????????????466?????????)
        cv.line(image, landmark_point[362], landmark_point[398], (0, 255, 0), 2)
        cv.line(image, landmark_point[398], landmark_point[384], (0, 255, 0), 2)
        cv.line(image, landmark_point[384], landmark_point[385], (0, 255, 0), 2)
        cv.line(image, landmark_point[385], landmark_point[386], (0, 255, 0), 2)
        cv.line(image, landmark_point[386], landmark_point[387], (0, 255, 0), 2)
        cv.line(image, landmark_point[387], landmark_point[388], (0, 255, 0), 2)
        cv.line(image, landmark_point[388], landmark_point[466], (0, 255, 0), 2)

        cv.line(image, landmark_point[466], landmark_point[390], (0, 255, 0), 2)
        cv.line(image, landmark_point[390], landmark_point[373], (0, 255, 0), 2)
        cv.line(image, landmark_point[373], landmark_point[374], (0, 255, 0), 2)
        cv.line(image, landmark_point[374], landmark_point[380], (0, 255, 0), 2)
        cv.line(image, landmark_point[380], landmark_point[381], (0, 255, 0), 2)
        cv.line(image, landmark_point[381], landmark_point[382], (0, 255, 0), 2)
        cv.line(image, landmark_point[382], landmark_point[362], (0, 255, 0), 2)

        # ??? (308????????????78?????????)
        cv.line(image, landmark_point[308], landmark_point[415], (0, 255, 0), 2)
        cv.line(image, landmark_point[415], landmark_point[310], (0, 255, 0), 2)
        cv.line(image, landmark_point[310], landmark_point[311], (0, 255, 0), 2)
        cv.line(image, landmark_point[311], landmark_point[312], (0, 255, 0), 2)
        cv.line(image, landmark_point[312], landmark_point[13], (0, 255, 0), 2)
        cv.line(image, landmark_point[13], landmark_point[82], (0, 255, 0), 2)
        cv.line(image, landmark_point[82], landmark_point[81], (0, 255, 0), 2)
        cv.line(image, landmark_point[81], landmark_point[80], (0, 255, 0), 2)
        cv.line(image, landmark_point[80], landmark_point[191], (0, 255, 0), 2)
        cv.line(image, landmark_point[191], landmark_point[78], (0, 255, 0), 2)

        cv.line(image, landmark_point[78], landmark_point[95], (0, 255, 0), 2)
        cv.line(image, landmark_point[95], landmark_point[88], (0, 255, 0), 2)
        cv.line(image, landmark_point[88], landmark_point[178], (0, 255, 0), 2)
        cv.line(image, landmark_point[178], landmark_point[87], (0, 255, 0), 2)
        cv.line(image, landmark_point[87], landmark_point[14], (0, 255, 0), 2)
        cv.line(image, landmark_point[14], landmark_point[317], (0, 255, 0), 2)
        cv.line(image, landmark_point[317], landmark_point[402], (0, 255, 0), 2)
        cv.line(image, landmark_point[402], landmark_point[318], (0, 255, 0), 2)
        cv.line(image, landmark_point[318], landmark_point[324], (0, 255, 0), 2)
        cv.line(image, landmark_point[324], landmark_point[308], (0, 255, 0), 2)

    return image


def draw_pose_landmarks(image, landmarks):
    image_width, image_height = image.shape[1], image.shape[0]

    landmark_point = []

    for index, landmark in enumerate(landmarks.landmark):
        landmark_x = min(int(landmark.x * image_width), image_width - 1)
        landmark_y = min(int(landmark.y * image_height), image_height - 1)
        landmark_point.append([landmark.visibility, (landmark_x, landmark_y)])

        if landmark.visibility < 0 or landmark.presence < 0:
            continue

        if index == 0:  # ???
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 1:  # ???????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 2:  # ????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 3:  # ???????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 4:  # ???????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 5:  # ????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 6:  # ???????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 7:  # ??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 8:  # ??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 9:  # ????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 10:  # ????????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 11:  # ??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 12:  # ??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 13:  # ??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 14:  # ??????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 15:  # ?????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 16:  # ?????????
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 17:  # ??????1(?????????)
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 18:  # ??????1(?????????)
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 19:  # ??????2(??????)
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 20:  # ??????2(??????)
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 21:  # ??????3(?????????)
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 22:  # ??????3(?????????)
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 23:  # ???(??????)
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)
        if index == 24:  # ???(??????)
            cv.circle(image, (landmark_x, landmark_y), 5, (0, 255, 0), 2)

    if len(landmark_point) > 0:
        # ??????
        if landmark_point[1][0] > 0 and landmark_point[2][0] > 0:
            cv.line(image, landmark_point[1][1], landmark_point[2][1], (0, 255, 0), 2)
        if landmark_point[2][0] > 0 and landmark_point[3][0] > 0:
            cv.line(image, landmark_point[2][1], landmark_point[3][1], (0, 255, 0), 2)

        # ??????
        if landmark_point[4][0] > 0 and landmark_point[5][0] > 0:
            cv.line(image, landmark_point[4][1], landmark_point[5][1], (0, 255, 0), 2)
        if landmark_point[5][0] > 0 and landmark_point[6][0] > 0:
            cv.line(image, landmark_point[5][1], landmark_point[6][1], (0, 255, 0), 2)

        # ???
        if landmark_point[9][0] > 0 and landmark_point[10][0] > 0:
            cv.line(image, landmark_point[9][1], landmark_point[10][1], (0, 255, 0), 2)

        # ???
        if landmark_point[11][0] > 0 and landmark_point[12][0] > 0:
            cv.line(image, landmark_point[11][1], landmark_point[12][1], (0, 255, 0), 2)

        # ??????
        if landmark_point[11][0] > 0 and landmark_point[13][0] > 0:
            cv.line(image, landmark_point[11][1], landmark_point[13][1], (0, 255, 0), 2)
        if landmark_point[13][0] > 0 and landmark_point[15][0] > 0:
            cv.line(image, landmark_point[13][1], landmark_point[15][1], (0, 255, 0), 2)

        # ??????
        if landmark_point[12][0] > 0 and landmark_point[14][0] > 0:
            cv.line(image, landmark_point[12][1], landmark_point[14][1], (0, 255, 0), 2)
        if landmark_point[14][0] > 0 and landmark_point[16][0] > 0:
            cv.line(image, landmark_point[14][1], landmark_point[16][1], (0, 255, 0), 2)

        # ??????
        if landmark_point[15][0] > 0 and landmark_point[17][0] > 0:
            cv.line(image, landmark_point[15][1], landmark_point[17][1], (0, 255, 0), 2)
        if landmark_point[17][0] > 0 and landmark_point[19][0] > 0:
            cv.line(image, landmark_point[17][1], landmark_point[19][1], (0, 255, 0), 2)
        if landmark_point[19][0] > 0 and landmark_point[21][0] > 0:
            cv.line(image, landmark_point[19][1], landmark_point[21][1], (0, 255, 0), 2)
        if landmark_point[21][0] > 0 and landmark_point[15][0] > 0:
            cv.line(image, landmark_point[21][1], landmark_point[15][1], (0, 255, 0), 2)

        # ??????
        if landmark_point[16][0] > 0 and landmark_point[18][0] > 0:
            cv.line(image, landmark_point[16][1], landmark_point[18][1], (0, 255, 0), 2)
        if landmark_point[18][0] > 0 and landmark_point[20][0] > 0:
            cv.line(image, landmark_point[18][1], landmark_point[20][1], (0, 255, 0), 2)
        if landmark_point[20][0] > 0 and landmark_point[22][0] > 0:
            cv.line(image, landmark_point[20][1], landmark_point[22][1], (0, 255, 0), 2)
        if landmark_point[22][0] > 0 and landmark_point[16][0] > 0:
            cv.line(image, landmark_point[22][1], landmark_point[16][1], (0, 255, 0), 2)

        # ??????
        if landmark_point[11][0] > 0 and landmark_point[23][0] > 0:
            cv.line(image, landmark_point[11][1], landmark_point[23][1], (0, 255, 0), 2)
        if landmark_point[12][0] > 0 and landmark_point[24][0] > 0:
            cv.line(image, landmark_point[12][1], landmark_point[24][1], (0, 255, 0), 2)
        if landmark_point[23][0] > 0 and landmark_point[24][0] > 0:
            cv.line(image, landmark_point[23][1], landmark_point[24][1], (0, 255, 0), 2)
    return image


def draw_bounding_rect(use_brect, image, brect):
    if use_brect:
        # ????????????
        cv.rectangle(image, (brect[0], brect[1]), (brect[2], brect[3]),
                     (0, 255, 0), 2)

    return image


if __name__ == '__main__':
    main()
