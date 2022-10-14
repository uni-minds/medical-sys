### API

| Path                                       | Get | Post | Put | Del | Example                         |  Return  |
|:-------------------------------------------|:---:|:----:|:---:|:---:|:--------------------------------|:--------:|
| /api/v1/login                              |  g  |  p   |     |     |                                 |          |
| /api/v1/user                               |  g  |  p   |  p  |  d  |                                 |          |
| /api/v1/media                              |  g  |      |     |     |                                 |          |
| /api/v1/label                              |  g  |  p   |     |  d  |                                 |          |
| /api/v1/algo                               |  g  |  p   |     |     |                                 |          |
| /api/v1/raw                                |  g  |      |     |     |                                 |          |
| **Blockchain**                             |     |      |     |     |                                 |          |
| /api/v1/blockchain/nodelist                |  g  |      |     |     |                                 |          |
| /api/v1/blockchain/tps                     |  g  |      |     |     |                                 |          |
| /api/v1/blockchain/height                  |  g  |      |     |     |                                 |          |
| **DB**                                     |     |      |     |     |                                 |          |
| /api/v1/database/ct/${class}/rs            |  g  |      |     |     |                                 |          |
| /api/v1/database/ct/${class}/wado          |  g  |      |     |     |                                 |          |
| **AI-Analysis/Search/Result**              |     |      |     |     |                                 |          |
| /api/v1/ai/${modal}/${class}/${algo_mode}  |     |  p   |     |     | ai/ct/cta/analysis_deep1 [data] | data=aid |
| /api/v1/ai/${modal}/${class}/result/${aid} |  g  |      |     |     | ai/ct/cta/result/AID            |          |
| **Sync**                                   |     |      |     |     |                                 |          |
| /api/v1/screen?action=sync                 |     |      |     |     |                                 |          |

### UI

| Path | Return |
|:-----|:-------|
/ui/ai/result/aid/

### import media (e.g. jpg/mp4)
python3 media_index.py -i tapvc-negative-4ap -o output1 -v 4ap -g tapvc-negative -k negative,tapvc
/ui/import?path=/root/output1

### Docker build
docker build -t medisys:v1 -f ./Dockerfile-ms .

### Docker run 
docker run --rm -p 80:80 \
    -v /data/medisys/database/:/data/database/ \
    -v /data/web:/usr/local/uni-ledger/medical-sys/application/web \
    medisys:v1