import json


def json_save(filename="data.json",datatype="us_media",data={}):
    jdata = {"type":datatype,"data":data}
    with open(filename,"w+",encoding='utf-8') as fp:
        json.dump(jdata,fp,ensure_ascii=False)

if __name__ == '__main__':
    json_save("data.json","dicom",data=[{"a":"1"},{"b":"2"}])

