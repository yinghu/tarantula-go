var Html = (function(){
    
    let _tasks = {};
    let _enum ={};
    let _category ={};

    let _caption = function(word){
        return word[0].toUpperCase() + word.slice(1);
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
    let _checkbox = function(prop,prefix){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<label class='tx-text-12'>");
        tem.push(_caption(prop.Name));
        tem.push("</label>");
        tem.push("<input id='");
        tem.push(prefix+"-"+prop.Name+"' ");
        tem.push("class='w3-check w3-round w3-border tx-text-16 tx-margin-left-8' type='checkbox'/>");
        tem.push("</div>");
        return tem.join("");       
    };
    let _button = function(name,prefix){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<span id='");
        tem.push(prefix+"-"+name+"' ");
        tem.push("class='w3-right w3-tag w3-green tx-text-18 tx-padding-button tx-margin-top-8 tx-margin-right-8'>");
        tem.push(name);
        tem.push("</span>");
        tem.push("</div>");
        return tem.join("");
    };

    let _selectEnum = function(prop,prefix){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<label class='tx-text-12'>");
        tem.push(_caption(prop.Name));
        tem.push("</label>");
        tem.push("<select id='");
        tem.push(prefix+'-'+prop.name);
        tem.push("' name='chs' class='w3-round w3-border tx-text-16 w3-select tx-margin-left-4'>");
        tem.push("<option value='1'>"+"One"+"</option>");
        tem.push("<option value='2'>"+"Two"+"</option>");
        tem.push("</select>");
        tem.push("</div>");
        return tem.join('');    
    };

    let _selectCategory = function(prop,prefix){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<label class='tx-text-12'>");
        tem.push(_caption(prop.Name));
        tem.push("</label>");
        tem.push("<span tx-property-name='");
        tem.push(prop.Name);    
        tem.push("' class='w3-right tx-margin-right-4 ");
        tem.push(prefix+"-upload'>");
        tem.push("<i class='material-symbols-outlined tx-orange-icon-20'>add</i></span>");
        tem.push("<select id='");
        tem.push(prefix+'-'+prop.name);
        tem.push("' name='chs' class='w3-round w3-border tx-text-16 w3-select tx-margin-left-4'>");
        tem.push("<option value='' disabled selected>click plus to load</option>");
        tem.push("</select>");
        tem.push("</div>");
        return tem.join("");  
    }

    let _upload = function(prop,prefix){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<label class='tx-text-12'>");
        tem.push(_caption(prop.Name));
        tem.push("</label>");
        tem.push("<span tx-property-name='");
        tem.push(prop.Name);    
        tem.push("' class='w3-right tx-margin-right-4 ");
        tem.push(prefix+"-upload'>");
        tem.push("<i class='material-symbols-outlined tx-orange-icon-20'>file_open</i></span>");
        tem.push("<input disabled id='");
        tem.push(prefix+"-"+prop.Name+"' ");
        tem.push("class='w3-input w3-round w3-border tx-text-16' type='text'/>");
        tem.push("</div>");
        return tem.join("");  
    }

    let _icon = function(prefix,n,icon,color){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<span id='");
        tem.push(prefix+"-"+n+"-"+icon+"' class='w3-right'>");
        tem.push("<i class='material-symbols-outlined tx-margin-right-8 tx-"+color+"-icon-24'>"+icon+"</i></span>");
        tem.push("</div>");
        return tem.join("");
    }

    let _form = function(containerId,prefix,category,closeable,callback){
        console.log(category);
        document.querySelector(containerId).innerHTML="";
        let tem=[];
        tem.push("<fieldset>");
        tem.push("<legend class='tx-text-20'>");
        tem.push(category.Scope.toUpperCase()+"/"+category.Description.toUpperCase());
        tem.push("</legend>");
        if(closeable){
            tem.push(_icon(prefix,category.Name,"close","red"));
        }
        category.Properties.forEach(prop=>{
            tem.push(_input(prop,prefix));    
        });
        tem.push(_button(category.Name,prefix,callback));
        tem.push("</fieldset>");
        document.querySelector(containerId).innerHTML += tem.join("");
        if(closeable){
            _eventWithId("#"+prefix+"-"+category.Name+"-close",()=>{
                document.querySelector(containerId).style.display='none';
            });
        }
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

    let _enumForm = function(containerId,callback){
        document.querySelector(containerId).innerHTML="";
        _enum ={ix:0};
        let tem=[];
        tem.push("<fieldset>");
        tem.push("<legend class='tx-text-20'>");
        tem.push("Enum");
        tem.push("</legend>");
        tem.push(_icon("ee","enum","close","red"));
        tem.push(_input({Name:"Name",Reference:"text"},"enum"));
        tem.push(_input({Name:"Entry",Reference:"text"},"enum"));
        tem.push(_input({Name:"Value",Reference:"number"},"enum"));
        tem.push(_icon("ee","enum","add","orange"));
        tem.push("<div class='w3-card-4 w3-round w3-border tx-text-12 w3-ul tx-margin-left-4'>");
        tem.push("<ul id='tx-create-enum-properties' class='w3-ul'>");
        tem.push("</ul></div>");
        tem.push(_button("Save","ee",callback));
        tem.push("</fieldset>");
        document.querySelector(containerId).innerHTML += tem.join("");
        _eventWithId("#ee-enum-close",()=>{
            document.querySelector(containerId).style.display='none';
        });
        _eventWithId("#ee-enum-add",()=>{
            let id = "entry"+_enum.ix;
            _enum.ix++;
            let e = document.querySelector("#enum-Entry").value;
            let v = document.querySelector("#enum-Value").value/1;
            let prop = e+" : "+v; 
            let li = "<li id='"+id+"'>"+prop+"<scan class='w3-right tx-enum-entry-remove' "+"tx-enum-entry-id='"+id+"'><i class='material-symbols-outlined tx-red-icon-24'>remove</i></span></li>";
            document.querySelector("#tx-create-enum-properties").innerHTML += li;
            _enum[id]={Name:e,Value:v};
            document.querySelectorAll(".tx-enum-entry-remove").forEach(a=>{
                a.onclick = ()=>{
                    let removeId = a.getAttribute("tx-enum-entry-id");
                    document.querySelector("#"+removeId).style.display="none";
                    delete _enum[removeId];
                    console.log(a.getAttribute("tx-enum-entry-id"));
                };            
            });
        });
        _eventWithId("#ee-Save",()=>{
            let n = document.querySelector("#enum-Name").value;
            _enum.Name = n;
            callback(_enum);
        });       
    };

    let _categoryHeader = function(containerId, prefix){
        let tem=[];
        tem.push("<fieldset>");
        tem.push("<legend class='tx-text-20'>");
        tem.push("Category/Header");
        tem.push("</legend>");
        tem.push(_input({Name:"Name",Reference:"text"},prefix));
        tem.push(_input({Name:"Desctiption",Reference:"text"},prefix));
        tem.push(_input({Name:"Version",Reference:"text"},prefix));
        tem.push(_checkbox({Name:"Rechargeable"},prefix));
        tem.push(_checkbox({Name:"Downloadable"},prefix));
        tem.push("</fieldset>");
        document.querySelector(containerId).innerHTML += tem.join("");
    }

    let _categoryForm = function(containerId,callback){
        document.querySelector(containerId).innerHTML="";
        _category ={ix:0};
        let tem=[];
        tem.push(_icon("cc","category","close","red"));
        document.querySelector(containerId).innerHTML += tem.join("");
        _categoryHeader(containerId,"cc-header");
        tem = [];
        tem.push("<fieldset>");
        tem.push("<legend class='tx-text-20'>");
        tem.push("Category/Properties");
        tem.push("</legend>");
        tem.push(_input({Name:"Name",Reference:"text"},"category"));
        tem.push(_selectEnum({Name:"Type",Reference:"text"},"category"));
        tem.push(_checkbox({Name:"Nullable",Reference:"text"},"category"));
        tem.push(_checkbox({Name:"Downloadable",Reference:"text"},"category"));
        tem.push(_icon("cc","category","add","orange"));
        tem.push("<div class='w3-card-4 w3-round w3-border tx-text-12 w3-ul tx-margin-left-4'>");
        tem.push("<ul id='tx-create-category-properties' class='w3-ul'>");
        tem.push("</ul></div>");
        tem.push(_button("Save","cc",callback));
        tem.push("</fieldset>");
        document.querySelector(containerId).innerHTML += tem.join("");
        _eventWithId("#cc-category-close",()=>{
            document.querySelector(containerId).style.display='none';
        });
         _eventWithId("#cc-category-add",()=>{
            console.log("add property");
        }); 
         _eventWithId("#cc-Save",()=>{
            //let n = document.querySelector("#enum-Name").value;
            //_enum.Name = n;
            callback(_category);
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
        enumForm : _enumForm,
        categoryForm : _categoryForm,
    };
})();