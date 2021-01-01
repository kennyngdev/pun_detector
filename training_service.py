import sys
import base64
import requests


def get_wave(fname):
    with open(fname) as infile:
        return base64.b64encode(infile.read())


endpoint = "https://snowboy.kitt.ai/api/v1/train/"

token = "e8a1876f4159b76f163319fa6324aaddc3b04b76"
hotword_name = "???"
language = "en"
age_group = "20_29"
gender = "M"
microphone = "macbook microphone"
# python training_service.py ./user_recording/1_new.wav ./user_recording/2_new.wav ./user_recording/3_new.wav saved_model.pmdl

if __name__ == "__main__":
    try:
        [_, wav1, wav2, wav3, out] = sys.argv
    except ValueError:
        print "Usage: %s wave_file1 wave_file2 wave_file3 out_model_name" % sys.argv[0]
        sys.exit()

    data = {
        "name": hotword_name,
        "language": language,
        "age_group": age_group,
        "gender": gender,
        "microphone": microphone,
        "token": token,
        "voice_samples": [
            {"wave": get_wave(wav1)},
            {"wave": get_wave(wav2)},
            {"wave": get_wave(wav3)}
        ]
    }

    response = requests.post(endpoint, json=data)
    if response.ok:
        with open("./resources/models/"+out+".pmdl", "w") as outfile:
            outfile.write(response.content)
        print "Saved model to '%s'." % out
    else:
        print "Request failed."
        print response.text
