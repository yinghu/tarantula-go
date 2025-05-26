var Admin = function(){

    let _hello = function(){
        return "hello";
    };
    let _login = function(){
        return "token"
    };
    return{
        hello : _hello,
        login : _login,
    };
}()