# 用于支持流媒体的本地化转换，正常情况下使用RTSP上传。
# 转换完成后需要手动在media表中添加媒体对应的目录信息
import glob
import os

src_folder = "/media/nas/心脏手术视频-安贞武玉多/datasets_1/"
desc_folder = "./rtsp/datasets_1"
dryrun = True
show_ignore = True


def convert_to_stream(video, output):
    if os.path.isdir(output):
        a = glob.glob(os.path.join(output, "*.ts"))
        if len(a) > 0:
            if show_ignore:
                print()
                print("error: output folder has already exists. ignore {}".format(output))
            return

    else:
        if dryrun:
            print("mkdir {}".format(output))

        else:
            os.mkdir(output)

    cmd = "ffmpeg -i \"{0}\" -codec copy -vbsf h264_mp4toannexb -map 0 -f segment -segment_list \"{1}/video.m3u8\" -segment_time 10 \"{1}/%05d.ts\"".format(
        video, output)

    if dryrun:
        print(cmd)
    else:
        os.system(cmd)


def main(src,desc="."):
    for root, _, files in os.walk(src):
        for file in files:
            filename, fileext = os.path.splitext(file)
            if fileext == ".mp4":
                output = os.path.join(desc, filename)
                video = os.path.join(root, file)
                convert_to_stream(video, output)

            elif fileext == ".ts":
                filepara = filename.split("_")
                if filepara[-1] == '0':
                    tmp = []
                    i = 0
                    while (1):
                        fn = "{}_{}{}".format('_'.join(filepara[:-1]), i, fileext)
                        target = os.path.join(root, fn)
                        if os.path.isfile(target):
                            tmp.append("{}".format(target))
                            i += 1
                        else:
                            break

                    str = tmp[0]
                    if len(tmp) > 1:
                        for i, f in enumerate(tmp):
                            if i != 0:
                                str = "{}|{}".format(str, f)
                            else:
                                str = "concat:\"{}".format(f)
                            #
                        str = str + '"'

                    output = os.path.join(desc, '_'.join(filepara[:-1]))
                    convert_to_stream(str, output)

                # video = os.path.join(root,file)
                # output = os.path.join(root,filename)
                # convert_to_stream(video,output)

        # print(root,dirs,files)


main(src_folder,desc_folder)
