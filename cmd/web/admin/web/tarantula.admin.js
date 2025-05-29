var Admin = function(){
    let presence = {}; 
    let _presence = function(){
        return presence;
    }
    let _login = function(payload,callback){
        let data = JSON.stringify(payload);
        let aj = new XMLHttpRequest();   
        aj.responseType = 'text';
        aj.onreadystatechange = function(){
            if(aj.status === 200 && aj.readyState === 4){
                let p = JSON.parse(aj.responseText);
                if(p.successful == true){
                    presence.token = p.token;
                    presence.systemId = p.systemId;
                    presence.stub = p.stub;
                    presence.login = p.login;
                    callback({successful:true});
                }else{
                    callback(p);
                }
            }
        };
        aj.open("POST","/admin/login",true);
        aj.setRequestHeader('Accept','application/json');
        aj.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        aj.setRequestHeader('Tarantula-payload-size',data.length);  
        aj.send(data);       
    };
    return{
        presence : _presence,
        login : _login,
    };
}()