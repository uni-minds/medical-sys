// 腹部横切面
let PresetCommon = {
    "default":{"color": "khaki"},
    "UV":{"tool":"com","cid":"UV","name":"脐静脉","color":"#ffe92c"},
    "ST":{"tool":"com","cid":"ST","name":"胃泡","color":"#9aff25"},
    "DV":{"tool":"com","cid":"DV","name":"静脉导管","color":"#ff3e10"},
    "IVC":{"tool":"com","cid":"IVC","name":"下腔静脉","color":"#FD843F"},
    "DAO":{"tool":"com","cid":"DAO","name":"降主动脉","color":"#E3FD5B"},
    "JZ":{"tool":"com","cid":"JZ","name":"脊柱","color":"#709DFD"},
    "LIVER":{"tool":"com","cid":"LIVER","name":"肝脏","color":"#FB43FD"},
    "JJM":{"tool":"com","cid":"JJM","name":"奇静脉","color":"#ffe2c3"},
    "DN":{"tool":"com","cid":"DN","name":"胆囊","color":"#f0e533"},
    "Mark":{"tool":"text","cid":"Mark","name":"备注"}
};
let PresetQuility ={
    "FQ5":{"tool":"q","cid":"FQ5","name":"优秀"},
    "FQ4":{"tool":"q","cid":"FQ4","name":"良好"},
    "FQ3":{"tool":"q","cid":"FQ3","name":"一般"},
    "FQ2":{"tool":"q","cid":"FQ2","name":"差"},
    "FQ1":{"tool":"q","cid":"FQ1","name":"不可评估"}
};
let PresetTime = {
    "SSMQ":{"tool":"time","cid":"SSMQ","name":"收缩末期"},
    "SZMQ":{"tool":"time","cid":"SZMQ","name":"舒张末期"},
    "SPEC":{"tool":"time","cid":"SPEC","name":"特殊时间"}
};
let PresetTSC = {
    "CZJM":{"tool":"com","cid":"CZJM","name":"垂直静脉","color":"#fcb1a3"},
    "YCXG":{"tool":"com","cid":"YCXG","name":"异常血管","color":"#ffe1a7"},
    "ERR1":{"tool":"com","cid":"ERR1","name":"异常结构\r1","color":"#FFA"},
    "ERR2":{"tool":"com","cid":"ERR2","name":"异常结构\r2","color":"#FFC"},
    "ERR3":{"tool":"com","cid":"ERR3","name":"异常结构\r3","color":"#FFE"}
};


let CRFButton_Common=PresetCommon;
let CRFButton_Q = PresetQuility;
let CRFButton_T = PresetTime;
let CRFButton_Spec = PresetTSC;

CRFButton_Common['default']={};
CRFButton_Common['default'].color="khari";