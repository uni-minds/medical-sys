function show3D(id,data,color){
    console.log(id);
    let refParent = document.getElementById(id);

    let points;
    let xa = -256;
    let ya = -300;
    let za = -256;
    let windowW=refParent.clientWidth;
    let windowH=refParent.clientHeight;



    function addPointsObj(scene, data) {
        let obj = buildLayer(data, color);
        scene.add(obj);
        return obj;
    }

    function addPlaneObj(scene) {
        let planeGeometry = new THREE.PlaneGeometry(512, 512, 0, 0);
        let planeMaterial = new THREE.MeshLambertMaterial({color: '#3d5cff'});
        let obj = new THREE.Mesh(planeGeometry, planeMaterial);
        obj.rotation.x = -0.5 * Math.PI;
        obj.position.x = 0;
        obj.position.y = ya;
        obj.position.z = 0;
        scene.add(obj);
        return obj
    }

    function buildLayer(data, color) {
        let geometry = new THREE.Geometry();
        data.forEach((item, index) => {
            geometry.vertices.push(new THREE.Vector3(item[0] + xa, item[2] + ya, item[1] + za));
        });
        let material = new THREE.PointsMaterial({color: color,size: 1.5});//材质对象
        let obj = new THREE.Points(geometry, material);//点模型对象
        return obj
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
    camera.position.x = 500;
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

function combinexyz(jsonpath){
    let datax=getJson(jsonpath+"center_x.json");
    let datay=getJson(jsonpath+"center_y.json");
    let dataz=getJson(jsonpath+"center_z.json");
    let data=[];
    for (let i=0;i<datax.length;i++) {
        data[i]=[datax[i],datay[i],dataz[i]]
    }
    console.log("combined:",data);
    return data;
}





