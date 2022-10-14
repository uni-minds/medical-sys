import sys
import pydicom as pd
import os
from tqdm import tqdm
import getopt
from json_io import json_save


def decode(val):
    return val.decode('utf-8').replace('\x00', '')


def get_dicom_ids(file):
    d = pd.dcmread(file)

    study_id = decode(d.get_item('StudyInstanceUID').value)
    series_id = decode(d.get_item('SeriesInstanceUID').value)
    instance_id = decode(d.get_item('SOPInstanceUID').value)
    instance_number = decode(d.get_item('InstanceNumber').value)
    series_number = decode(d.get_item('SeriesNumber').value)

    return study_id, series_id, instance_id, int(series_number), int(instance_number)


def proc_folder(folder,group,keyword):
    dicom_list = []
    result_data = {}
    keywords = []

    for key in keyword.split(","):
        if key == "":
            continue
        keywords.append(key)

    print("processing: {}".format(folder))
    for root, dirs, files in os.walk(folder, topdown=False):
        for name in files:
            if ".dcm" in name:
                dicom_list.insert(0, os.path.join(root, name))

    count = 0
    pbar = tqdm(dicom_list)
    for file in pbar:
        pbar.set_description("proc: {}".format(file))

        study_id, series_id, instance_id, s_number, i_number = get_dicom_ids(file)

        study_data = result_data.get(study_id)
        if study_data is not None:
            series_data = study_data.get(series_id)
            if instance_id in series_data:
                print("instance already included: {}".format(file))
            else:
                series_data.append(instance_id)
        else:
            study_data={series_id:[instance_id]}

        count += 1
        result_data[study_id]=study_data
        result= {"dicom_tree":result_data,"group":group,"keywords":keywords}

    return count, result


if __name__ == '__main__':
    output = "data.json"
    src = "dicoms"
    group = ''
    keywords = ''

    usage = '{} -i <src> -o <dest> -g <group> -k <keyword1,keyword2>'.format(sys.argv[0])

    try:
        opts, args = getopt.getopt(sys.argv[1:], "hi:o:g:k:",
                                   ["help", "src=", "output=", "group=", "keywords="])
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
        elif opt_name in ("-g", "--group"):
            group = opt_value
        elif opt_name in ("-k", "--keywords"):
            keywords = opt_value

    cnt, ret = proc_folder(src, group, keywords)
    if cnt > 0:
        json_save(output, "us_dicom", ret)
        print("{} records saved.".format(cnt))
