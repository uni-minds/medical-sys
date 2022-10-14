import getopt
import os
from ffmpy3 import FFmpeg
from tqdm import tqdm
import sys
from json_io import json_save


def proc_folder(src_root, desc_root, view="", group="", keyword=""):
    print("index from {} to {}".format(src_root, desc_root))
    keywords = []
    li = []
    result = []

    for key in keyword.split(","):
        if key == "":
            continue
        keywords.append(key)

    descJson = os.path.join(desc_root, "data.json")
    desc_media_folder = os.path.join(desc_root, "media")
    try:
        os.makedirs(desc_media_folder)
    except:
        print("Target folder already existed")

    for root, dirs, files in os.walk(src_root):
        for name in files:
            if ".DS_Store" in name:
                continue
            if ".db" in name:
                continue
            li.append(os.path.join(root, name))

    for src_media_file in tqdm(li):
        data = media_process(src_media_file, desc_media_folder, desc_root, view, group, keywords)
        result.append(data)

    json_save(descJson, "us_media", {"media_info":result,"group":group,"keywords":keywords})


def media_process(src_media, dest_media_folder, dest_root, view, group, keywords, fcode="", patientid="", machineid=""):
    descript = ""
    srcFolder = os.path.dirname(src_media)
    srcDescript = os.path.join(srcFolder, "descript.txt")
    if os.path.isfile(srcDescript):
        with open(srcDescript, 'r', encoding='gb2312') as fp:
            descript = fp.read()

    detail = {}
    output = media_convert(src_media, dest_media_folder)
    detail["source"] = src_media
    detail["target"] = output[len(dest_root) + 1:]
    detail["descript"] = descript
    detail["groupname"] = group
    detail["view"] = view
    detail["keywords"] = keywords
    detail["fcode"] = fcode
    detail["patientid"] = patientid
    detail["machineid"] = machineid
    return detail


def media_convert(srcFile, descFolder):
    filename = srcFile.split('/')[-1]
    filebase = filename[:-4]
    fileext = filename.split(".")[-1].lower()
    print(fileext)
    outFile = filebase.replace(' ', '_').replace('(', '').replace(')', '')

    output = ""
    if fileext in ("mp4", "avi", "ogv", "jpg", "bmp", "png"):
        output = os.path.join(descFolder, "{}.ogv".format(outFile))
        count = 0
        while os.path.isfile(output):
            print('target {} existed, change it'.format(output))
            output = os.path.join(descFolder, "{}_{}.ogv".format(outFile, count))
            count += 1
        ff = FFmpeg(inputs={srcFile: None}, outputs={output: '-b:v 50000k -auto-alt-ref 0'})
        ff.run()

    else:
        print("file type is not recognize:{} / {}".format(srcFile, filebase))

    return output


if __name__ == '__main__':
    output = ''
    src = ''
    view = ''
    group = ''
    keywords = ''

    usage = '{} -i <src> -o <output> -v <view> -g <group> -k <keyword1,keyword2>'.format(sys.argv[0])

    try:
        opts, args = getopt.getopt(sys.argv[1:], "hi:o:v:g:k:", ["src=", "output=", "view=", "group=", "keyword="])
    except getopt.GetoptError:
        print(usage)
        sys.exit(2)

    for opt_name, opt_value in opts:
        if opt_name in ('-h', '--help'):
            print(usage)
            sys.exit(1)
        elif opt_name in ("-i", "--src"):
            src = opt_value
        elif opt_name in ("-o", "--output"):
            output = opt_value
        elif opt_name in ("-v", "--view"):
            view = opt_value
        elif opt_name in ("-g", "--group"):
            group = opt_value
        elif opt_name in ("-k", "--keyword"):
            keywords = opt_value

    proc_folder(src, output, view, group, keywords)
