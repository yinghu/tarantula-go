var Html = (function(){
    
    let _tasks = {};

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
        tem.push("class='w3-right w3-tag w3-green tx-text-18 tx-padding-button tx-margin-top-8 tx-margin-right-8'>");
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

    let _taskList = function(containerId,prefix,tbar,callback){
        console.log(tbar);
        document.querySelector(containerId).innerHTML="";
        let tem=[];
        tbar.Tasks.forEach(task=>{
            _tasks[task.Name]=task;
            tem.push("<div tx-task-name='");
            tem.push(task.Name+"' ")  
            tem.push("class='w3-display-container  w3-border-bottom w3-border-red tx-content-48 ");
            tem.push("tx-"+prefix+"-action'>");
            tem.push("<div class='w3-display-bottomright w3-padding'>");
            tem.push("<span><i class='material-symbols-outlined tx-margin-left-8 tx-orange-icon-24'>double_arrow</i></span>");
            tem.push("</div>");
            tem.push("<div class='w3-display-bottomleft w3-padding'>");
            tem.push("<span class='tx-margin-left-8 tx-text-24'>");
            tem.push(task.Name);
            tem.push("</span>");
            tem.push("</div>");
            tem.push("</div>");
        });
        document.querySelector(containerId).innerHTML = tem.join("");
        document.querySelectorAll(".tx-"+prefix+"-action").forEach(a=>{
            a.onclick = ()=>{
                callback(a.getAttribute("tx-task-name"));    
            };
        });
    };

    let _jobList = function(containerId,prefix,tn,callback){
        let task = _tasks[tn];
        console.log(task);
        document.querySelector(containerId).innerHTML="";
        let tem=[];
        tem.push("<span class='w3-bar-item w3-left w3-teal w3-tag tx-text-24 tx-padding-button'><i class='material-symbols-outlined tx-orange-icon-24 tx-margin-top-8'>settings</i>");
        tem.push(" "+task.Name);
        tem.push("</span>");
        task.Jobs.forEach(job=>{
            tem.push("<span tx-job-name='");
            tem.push(job.Callback+"' ");
            tem.push("class='w3-bar-item w3-right w3-green w3-tag tx-text-24 tx-padding-button tx-margin-left-4 ");
            tem.push("tx-"+prefix+"-action")
            tem.push("'>"+job.Name+"</span>");
        });
        document.querySelector(containerId).innerHTML = tem.join("");
        document.querySelectorAll(".tx-"+prefix+"-action").forEach(a=>{
            a.onclick = ()=>{
                callback(a.getAttribute("tx-job-name"));    
            };
        });
    };

    return {
        messageWithId : _messageWithId,
        openWithId : _openWithId,
        closeWithId : _closeWithId,
        eventWithId : _eventWithId,
        form : _form,
        taskList : _taskList,
        jobList : _jobList,
    };
})();