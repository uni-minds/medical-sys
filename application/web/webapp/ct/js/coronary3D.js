function show3D(id,data,color){
    let refParent = document.getElementById(id);
    console.log(refParent,data,id);

    refParent.onmouseenter=function() {
        console.log("E",this,this.clientWidth,this.clientHeight)
    };

    let points;
    let xa = -256;
    let ya = 300;
    let za = -256;
    let windowW=refParent.clientWidth;
    let windowH=refParent.clientHeight;

    function addPointsObj(scene, data) {

        /*let obj = buildLayer(data["lad"], '#ff0006');
        scene.add(obj);
        obj = buildLayer(data["lcx"], '#0aff00');
        scene.add(obj);
        obj = buildLayer(data["lm"], '#000dff');
        scene.add(obj);
        obj = buildLayer(data["rca"],'#ff8900');
        scene.add(obj);
        */
        var obj = buildLayer(data,'#ff8900');
        scene.add(obj);
        return obj;
    }

    function addPlaneObj(scene) {
        let planeGeometry = new THREE.PlaneGeometry(512, 512, 0, 0);
        let planeMaterial = new THREE.MeshLambertMaterial({color: '#3d5cff'});
        let obj = new THREE.Mesh(planeGeometry, planeMaterial);
        obj.rotation.x = -0.5 * Math.PI;
        obj.position.x = 0;
        obj.position.y =0;
        obj.position.z = 0;
        scene.add(obj);
        return obj;
    }

    function buildLayer(data, color) {
        let geometry = new THREE.Geometry();
        data.forEach((item, index) => {
            geometry.vertices.push(new THREE.Vector3(item[0] + xa,  ya-item[2], item[1] + za));
        });
        let material = new THREE.PointsMaterial({color: color,size: 1.5});//材质对象
        let obj = new THREE.Points(geometry, material);//点模型对象
        return obj;
    }

    function animate() {
        requestAnimationFrame(animate);
        model.rotation.y += 0.01;
        renderer.render(scene, camera);
    }

    function render() {
        renderer.render(scene, camera);
    }

    let scene = new THREE.Scene();
    let camera = new THREE.PerspectiveCamera(80, windowW / windowH, 10, 20000);
    camera.position.x = 100;
    camera.position.y = 100;
    camera.position.z = 500;
    camera.lookAt(scene.position);

    let renderer = new THREE.WebGLRenderer();
    renderer.setClearColor('#FFFFFF');
    renderer.setSize(windowW , windowH);

    refParent.appendChild(renderer.domElement);

    controls = new THREE.OrbitControls(camera, renderer.domElement);
    controls.addEventListener('change', render);

    let model = addPointsObj(scene, data);
    addPlaneObj(scene);
    animate();
    render();
}

//data["lad"]=loadJson("webapp/output/coronary/json/lad_list.json");
//data["lcx"]=loadJson("webapp/output/coronary/json/lcx_list.json");
//data["lm"]=loadJson("webapp/output/coronary/json/lm_list.json");
//data["rca"]=loadJson("webapp/output/coronary/json/rca_list.json");

ccta_maskdata =loadJson("/webapp/temp/ccta/json/mask.json");

function upsideDown(data) {
    var temp=[];
    var l=data.length;
    for (var i=0;i<l;i++){
        temp.push(data[l-i]);
    }
    return temp;
}
show3D("area3D",ccta_maskdata,'#ff244e');
