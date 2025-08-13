var Player = (function(){
    let _query = {authenticated:false,confirm:()=>{},job:{}};
    
    let _getQuery = function(){
        return _query;
    };

    let _login = function(auth){
        //if (_query.authenticated) return;
        if (auth.hasOwnProperty("token") && auth.successful){
            _query.authenticated = true;
            _query.token = auth.token;
        }
    };

    let _logout = function(){
        _query.authenticated = false;
        _query.token = "";
    };

    let _getJson = function(path,callback){
        let aj = new XMLHttpRequest();   
        aj.responseType = 'text';
        aj.onreadystatechange = function(){
            if(aj.status === 200 && aj.readyState === 4){
                callback(JSON.parse(aj.responseText));
            }
        };
        aj.open("GET",path,true);
        aj.setRequestHeader('Accept','application/json');
        aj.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        if (_query.authenticated){
            aj.setRequestHeader("Authorization","Bearer "+_query.token);
        }
        aj.send();
    };

    let _postJson = function(path,data,callback){
        let payload = JSON.stringify(data);
        let aj = new XMLHttpRequest();   
        aj.responseType = 'text';
        aj.onreadystatechange = function(){
            if(aj.status === 200 && aj.readyState === 4){
                callback(JSON.parse(aj.responseText));
            }
        };
        aj.open('POST',path,true);
        aj.setRequestHeader('Accept','application/json');
        aj.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        if (_query.authenticated){
            aj.setRequestHeader("Authorization","Bearer "+_query.token);
        }
        aj.send(payload);
    };
    return {
        query : _getQuery,
        login : _login,
        logout : _logout,
        getJson : _getJson,
        postJson : _postJson,
    };
})();