//ROVT 右室流出道
let PresetCommon = {
    "default":{"color": "khaki"},
    "JZ":{"tool":"com","cid":"JZ","name":"脊柱","color":"#709DFD"},
    "XG":{"tool":"com","cid":"XG","name":"胸骨","color":"#ffe92c"},
    "LG":{"tool":"com","cid":"LG","name":"肋骨1","color":"#FD843F"},
    "LG2":{"tool":"com","cid":"LG2","name":"肋骨2","color":"#FD84EF"},
    "DAO":{"tool":"com","cid":"DAO","name":"降主动脉\rDAO","color":"#E3FD5B"},
    "LA":{"tool":"com","cid":"LA","name":"左房\rLA","color":"#76ff3e"},
    "LV":{"tool":"com","cid":"LV","name":"左室\rLV","color":"#b9ff94"},
    "RA":{"tool":"com","cid":"RA","name":"右房\rRA","color":"#f4e7d5"},
    "RV":{"tool":"com","cid":"RV","name":"右室\rRV","color":"#ff37d5"},
    "PAB1":{"tool":"com","cid":"PAB1","name":"肺动脉瓣\r1","color":"#ff7096"},
    "PAB2":{"tool":"com","cid":"PAB2","name":"肺动脉瓣\r2","color":"#5ce63e"},
    "PA":{"tool":"com","cid":"PA","name":"肺动脉\rPA","color":"#ffbd63"},
    "LPA":{"tool":"com","cid":"LPA","name":"左肺动脉\rLPA","color":"#ff7096"},
    "RPA":{"tool":"com","cid":"RPA","name":"右肺动脉\rRPA","color":"#5ce63e"},
    "AO":{"tool":"com","cid":"AO","name":"主动脉\rAO","color":"#E3FD5B"},
    "DA":{"tool":"com","cid":"DA","name":"动脉导管\r","color":"#ff7096"},
    "SJBGY":{"tool":"com","cid":"SJBGY","name":"三尖瓣\r隔叶","color":"#ff7096"},
    "SJBQY":{"tool":"com","cid":"SJBQY","name":"三尖瓣\r前叶","color":"#5ce63e"},
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
    "EC_PAB1":{"tool":"com","cid":"EC_PAB1","name":"异常\r肺动脉瓣","color":"#fcb1a3"}
};

let CRFButton_Common=PresetCommon;
let CRFButton_Q = PresetQuility;
let CRFButton_T = PresetTime;
let CRFButton_Spec = PresetTSC;

CRFButton_Common['default']={};
CRFButton_Common['default'].color="khari";