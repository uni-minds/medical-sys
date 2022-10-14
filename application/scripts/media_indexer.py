import getopt
import os
from ffmpy3 import FFmpeg
import json
from concurrent import futures
from tqdm import tqdm
import sys
import shutil

def start(src_root, desc_root, view="", group="", keyword=""):
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
        # try:
        #     data = media_process(src_media_file, desc_media_folder, desc_root, view, group, keywords)
        # except:
        #     print("Error convert {}".format(src_media_file))
        #     continue
        result.append(data)

    with open(descJson, 'w', encoding='utf-8') as fp:
        json.dump(result, fp, ensure_ascii=False)


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
    detail["target"] = output[len(dest_root)+1:]
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
    if fileext == "mp4" or fileext =="avi":
        output = os.path.join(descFolder, "{}.ogv".format(outFile))
        count = 0
        while os.path.isfile(output):
            print('target {} existed, change it'.format(output))
            output = os.path.join(descFolder, "{}_{}.ogv".format(outFile, count))
            count += 1
        ff = FFmpeg(inputs={srcFile: None}, outputs={output: '-b:v 50000k -auto-alt-ref 0'})
        ff.run()

    elif fileext == "jpg" or fileext == "bmp" or fileext == "png":
        output = os.path.join(descFolder, "{}{}".format(outFile,fileext))
        shutil.copyfile(srcFile,output)

    else:
        print("file type is not recognize:{} / {}".format(srcFile,filebase))

    return output


def main(argv):
    src = ''
    desc = ''
    view = ''
    group = ''
    keyword = ''
    try:
        opts, args = getopt.getopt(argv, "hi:o:v:g:k:", ["src=", "dest=", "view=", "group=", "keyword="])
    except getopt.GetoptError:
        print('.py -i <src> -o <dest> -v <view> -g <group> -k <keyword1,keyword2>')
        sys.exit(2)

    for opt, arg in opts:
        if opt == "-h":
            print('.py -i <src> -o <dest> -v <view> -g <group> -k <keyword1,keyword2>')
            sys.exit(1)
        elif opt in ("-i", "--src"):
            src = arg
        elif opt in ("-o", "--dest"):
            desc = arg
        elif opt in ("-v", "--view"):
            view = arg
        elif opt in ("-g", "--group"):
            group = arg
        elif opt in ("-k", "--keyword"):
            keyword = arg

    start(src, desc, view, group, keyword)


if __name__ == '__main__':
    main(sys.argv[1:])
