### API

Path|Get|Post|Put|Del|Example|Return
:---|:---:|:---:|:---:|:---:|:---|:---:
/api/v1/login |g|p
/api/v1/user  |g|p|p|d
/api/v1/media |g
/api/v1/label |g|p| |d
/api/v1/algo  |g|p
/api/v1/raw  |g|
**Blockchain**|
/api/v1/blockchain/nodelist|g|
/api/v1/blockchain/tps|g|
/api/v1/blockchain/height|g|
**DB**|
/api/v1/database/ct/${class}/rs|g
/api/v1/database/ct/${class}/wado|g
**AI-Analysis/Search/Result**|
/api/v1/ai/${modal}/${class}/${algo_mode} | |p| | |ai/ct/cta/analysis_deep1 [data]|data=aid
/api/v1/ai/${modal}/${class}/result/${aid}|g| | | |ai/ct/cta/result/AID

### UI

Path|Return
:---|:---
/ui/ai/result/aid/
