let PresetCommon = {
    "default":{"color": "khaki"},
    "XG":{"tool":"com","cid":"XG","name":"胸骨","color":"#ffe92c"},
    "JZ":{"tool":"com","cid":"JZ","name":"脊柱","color":"#9aff25"},
    "DA":{"tool":"com","cid":"DA","name":"降主动脉","color":"#ff3e10"},
    "LgZ1":{"tool":"com","cid":"LgZ1","name":"真肋骨1","color":"#FD843F"},
    "LgZ2":{"tool":"com","cid":"LgZ2","name":"真肋骨2","color":"#E3FD5B"},
    "LgJ1":{"tool":"com","cid":"LgJ1","name":"假肋骨1","color":"#709DFD"},
    "LgJ2":{"tool":"com","cid":"LgJ2","name":"假肋骨2","color":"#FB43FD"},
    "XJWM":{"tool":"com","cid":"XJWM","name":"心肌外膜","color":"#ffe1a3"},
    "FJMZ":{"tool":"com","cid":"FJMZ","name":"肺静脉左","color":"#dbff8e"},
    "FJM":{"tool":"com","cid":"FJM","name":"肺静脉右","color":"#6487ff"},
    "SZ":{"tool":"com","cid":"SZ","name":"房间隔\r（原发）","color":"#1bdaff"},
    "fjgJF":{"tool":"com","cid":"fjgJF","name":"房间隔\r（继发）","color":"#5BFF7F"},
    "sjg":{"tool":"com","cid":"sjg","name":"室间隔","color":"#FFB0F4"},
    "LA":{"tool":"com","cid":"LA","name":"左房\rLA","color":"#76ff3e"},
    "LV":{"tool":"com","cid":"LV","name":"左室\rLV","color":"#b9ff94"},
    "EJBQY":{"tool":"com","cid":"EJBQY","name":"二尖瓣\r前叶","color":"#a3ebff"},
    "EJBHY":{"tool":"com","cid":"EJBHY","name":"二尖瓣\r后叶","color":"#8bb9ff"},
    "RA":{"tool":"com","cid":"RA","name":"右房\rRA","color":"#fbaeff"},
    "RV":{"tool":"com","cid":"RV","name":"右室\rRV","color":"#ff37d5"},
    "SJBGY":{"tool":"com","cid":"SJBGY","name":"三尖瓣\r隔叶","color":"#ff7096"},
    "SJBQY":{"tool":"com","cid":"SJBQY","name":"三尖瓣\r前叶","color":"#5ce63e"},
    "RYKBM":{"tool":"com","cid":"RYKBM","name":"卵圆孔\r瓣膜","color":"#ffbd63"},
    "RYKKK":{"tool":"com","cid":"RYKKK","name":"卵圆孔\r开口","color":"#ff228f"},
    "Mark":{"tool":"text","cid":"Mark","name":"备注"}
};
let PresetQuility = {
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
    "TSC1":{"tool":"com","cid":"TSC1","name":"TSC\r1","color":"#FFA"},
    "TSC2":{"tool":"com","cid":"TSC2","name":"TSC\r2","color":"#FFC"},
    "TSC3":{"tool":"com","cid":"TSC3","name":"TSC\r3","color":"#FFE"},
    "CS":{"tool":"com","cid":"CS","name":"冠状静脉窦\rCS","color":"#F0E"},
    "SJBYC":{"tool":"com","cid":"SJBYC","name":"三尖瓣\r异常","color":"#a3eb70"},
    "EJBEC":{"tool":"com","cid":"EJBEC","name":"二尖瓣\r异常","color":"#a3eb00"},
    "GTFSB1":{"tool":"com","cid":"GTFSB1","name":"共同房室瓣\r1","color":"#F00"},
    "GTFSB2":{"tool":"com","cid":"GTFSB2","name":"共同房室瓣\r2","color":"#F01"},
    "FJMGTQ":{"tool":"com","cid":"FJMGTQ","name":"肺静脉共同腔","color":"#F02"},
    "DXS":{"tool":"com","cid":"DXS","name":"单心室","color":"#F03"},
    "DXF":{"tool":"com","cid":"DXF","name":"单心房","color":"#F04"}  
};


let CRFButton_Common=PresetCommon;
let CRFButton_Q = PresetQuility;
let CRFButton_T = PresetTime;
let CRFButton_Spec = PresetTSC;

CRFButton_Common['default']={};
CRFButton_Common['default'].color="khari";