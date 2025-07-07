var Admin = (function(){
    let _cache = {};
    let _test = function(){
        return "hello";
    };
    return {
        test : _test,
    };
})();