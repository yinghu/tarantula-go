var Html = (function(){
    
    let _caption = function(word){
        return word[0].toUpperCase() + word.slice(1).toLowerCase()
    };
    let _messageWithId = function(id,message){
        document.querySelector(id).innerHTML = message;
    };
    let _openWithId = function(id){
        document.querySelector(id).style.display = "block";
    };
    let _closeWithId = function(id){
        document.querySelector(id).style.display = "none";
    };
    let _eventWithId = function(id,callback){
        document.querySelector(id).onclick = callback 
    };
    let _input = function(prop,prefix){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<label class='tx-text-12'>");
        tem.push(_caption(prop.Name));
        tem.push("</label>");
        tem.push("<input id='");
        tem.push(prefix+"-"+prop.Name+"' ");
        tem.push("class='w3-input w3-round w3-border tx-text-16' type='");
        tem.push(prop.Reference+"'/>");
        tem.push("</div>");
        return tem.join("");       
    };
    let _button = function(category,prefix){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<span id='");
        tem.push(prefix+"-"+category.Name+"' ");
        tem.push("class='w3-right w3-tag w3-green w3-round w3-border w3-border-red tx-text-18 tx-padding-button tx-margin-top-8 tx-margin-bottom-8'>");
        tem.push(category.Name);
        tem.push("</span>");
        tem.push("</div>");
        return tem.join("");
    };

    let _form = function(containerId,prefix,category,callback){
        console.log(category);
        document.querySelector(containerId).innerHTML="";
        let tem=[];
        tem.push("<fieldset>");
        tem.push("<legend class='tx-text-20'>");
        tem.push(category.Scope.toUpperCase()+"/"+category.Description.toUpperCase());
        tem.push("</legend>");
        category.Properties.forEach(prop=>{
            tem.push(_input(prop,prefix));    
        });
        tem.push(_button(category,prefix,callback));
        tem.push("</fieldset>");
        document.querySelector(containerId).innerHTML += tem.join("");
        _eventWithId("#"+prefix+"-"+category.Name,()=>{
            let data ={};
            category.Properties.forEach(p=>{
                data[p.Name] = document.querySelector("#"+prefix+"-"+p.Name).value;
            });
            callback(data);
        });
    };


    return {
        messageWithId : _messageWithId,
        openWithId : _openWithId,
        closeWithId : _closeWithId,
        eventWithId : _eventWithId,
        form : _form,
    };
})();