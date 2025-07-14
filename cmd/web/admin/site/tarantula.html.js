var Html = (function(){
    
    let _tasks = {};
    let _enum ={};
    let _category ={};
    let _task ={};
    
    let _taskChanged = tn=>{
        console.log("Task changed :"+tn);
    };
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
    let _registerTaskChangeListener = function(listener){
        _taskChanged = listener;
    };
    let _typeList = function(tlist){
        _category ={ix:0,types:tlist.Types,headers:tlist.Headers};
    };
    let _taskList = function(conf,tbar,callback){
        console.log(tbar);
        document.querySelector(conf.id).innerHTML="";
        let tem=[];
        tbar.Tasks.forEach(task=>{
            _tasks[task.Name]=task;
            tem.push("<div tx-task-name='");
            tem.push(task.Name+"' ")  
            tem.push("class='w3-display-container  w3-border-bottom w3-border-red tx-content-48 ");
            tem.push("tx-"+conf.prefix+"-action'>");
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
        document.querySelector(conf.id).innerHTML = tem.join("");
        document.querySelectorAll(".tx-"+conf.prefix+"-action").forEach(a=>{
            a.onclick = ()=>{
                callback(a.getAttribute("tx-task-name"));    
            };
        });
    };
    let _jobList = function(conf,tn,callback){
        let task = _tasks[tn];
        _task.Name = tn;
        console.log(task);
        _taskChanged(_task.Name);
        document.querySelector(conf.id).innerHTML="";
        let tem=[];
        tem.push("<span class='w3-bar-item w3-left w3-teal w3-tag tx-text-24 tx-padding-button'><i class='material-symbols-outlined tx-orange-icon-24 tx-margin-top-8'>settings</i>");
        tem.push(" "+task.Name);
        tem.push("</span>");
        task.Jobs.forEach(job=>{
            tem.push("<span tx-job-name='");
            tem.push(job.Callback+"' ");
            tem.push("class='w3-bar-item w3-right w3-green w3-tag tx-text-24 tx-padding-button tx-margin-left-4 ");
            tem.push("tx-"+conf.prefix+"-action")
            tem.push("'>"+job.Name+"</span>");
        });
        document.querySelector(conf.id).innerHTML = tem.join("");
        document.querySelectorAll(".tx-"+conf.prefix+"-action").forEach(a=>{
            a.onclick = ()=>{
                callback(a.getAttribute("tx-job-name"));    
            };
        });
    };
     let _categoryList = function(conf,clist,callback){
        document.querySelector(conf.id).innerHTML = "";
        let tem =[];
        clist.forEach(c=>{
            tem.push("<span tx-category-id='"+c.Id+"' ");
            tem.push("class='w3-bar-item w3-green w3-tag tx-text-20 tx-padding-button tx-margin-right-4 tx-margin-bottom-4 tx-"+conf.prefix+"-opt'>");
            tem.push(c.Name);
            tem.push("</span>");    
        });
        document.querySelector(conf.id).innerHTML = tem.join("");
        document.querySelectorAll(".tx-"+conf.prefix+"-opt").forEach(a=>{
            a.onclick = ()=>{
                callback(a.getAttribute("tx-category-id"));    
            };
        });
    };
    let _enumForm = function(conf,callback){
        document.querySelector(conf.id).innerHTML="";
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
        tem.push("<ul id='create-enum-properties' class='w3-ul'>");
        tem.push("</ul></div>");
        tem.push(_button("Save","ee",callback));
        tem.push("</fieldset>");
        document.querySelector(conf.id).innerHTML += tem.join("");
        _eventWithId("#ee-enum-close",()=>{
            document.querySelector(conf.id).style.display='none';
        });
        _eventWithId("#ee-enum-add",()=>{
            let id = "entry"+_enum.ix;
            _enum.ix++;
            let e = document.querySelector("#enum-Entry").value;
            let v = document.querySelector("#enum-Value").value/1;
            let prop = e+" : "+v; 
            let li = "<li id='"+id+"'>"+prop+"<scan class='w3-right enum-entry-remove' "+"enum-entry-id='"+id+"'><i class='material-symbols-outlined tx-red-icon-24'>remove</i></span></li>";
            document.querySelector("#create-enum-properties").innerHTML += li;
            _enum[id]={Name:e,Value:v};
            document.querySelectorAll(".enum-entry-remove").forEach(a=>{
                a.onclick = ()=>{
                    let removeId = a.getAttribute("enum-entry-id");
                    document.querySelector("#"+removeId).style.display="none";
                    delete _enum[removeId];
                };            
            });
        });
        _eventWithId("#ee-Save",()=>{
            let n = document.querySelector("#enum-Name").value;
            _enum.Name = n;
            callback(_enum);
        });       
    };
    //END

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

    let _selectType = function(prop,prefix,types){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<label class='tx-text-12'>");
        tem.push(_caption(prop.Name));
        tem.push("</label>");
        tem.push("<select id='");
        tem.push(prefix+'-'+prop.Name);
        tem.push("' name='chs' class='w3-round w3-border tx-text-16 w3-select tx-margin-left-4'>");
        for(const k of Object.keys(types)){
            tem.push("<option value='"+k+"'>"+k+"</option>");
        }
        tem.push("</select>");
        tem.push("</div>");
        return tem.join("");  
    }

    let _selectReference = function(prop,prefix){
        let tem=[];
        tem.push("<div id='"+prefix+"-"+prop.Name+"-box' "+"class='w3-panel'>");
        tem.push("<label class='tx-text-12'>");
        tem.push(_caption(prop.Name));
        tem.push("</label>");
        tem.push("<select id='");
        tem.push(prefix+'-'+prop.Name);
        tem.push("' name='chs' class='w3-round w3-border tx-text-16 w3-select tx-margin-left-4'>");
        tem.push("</select>");
        tem.push("</div>");
        return tem.join("");  
    }

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

    let _form = function(conf,category,callback){
        document.querySelector(conf.id).innerHTML="";
        let tem=[];
        tem.push("<fieldset>");
        tem.push("<legend class='tx-text-20'>");
        tem.push(category.Scope.toUpperCase()+"/"+category.Description.toUpperCase());
        tem.push("</legend>");
        if(conf.closeable){
            tem.push(_icon(conf.prefix,category.Name,"close","red"));
        }
        category.Properties.forEach(prop=>{
            tem.push(_input(prop,conf.prefix));    
        });
        tem.push(_button(category.Name,conf.prefix,callback));
        tem.push("</fieldset>");
        document.querySelector(conf.id).innerHTML += tem.join("");
        if(conf.closeable){
            _eventWithId("#"+conf.prefix+"-"+category.Name+"-close",()=>{
                document.querySelector(conf.id).style.display='none';
            });
        }
        _eventWithId("#"+conf.prefix+"-"+category.Name,()=>{
            let data ={};
            category.Properties.forEach(p=>{
                data[p.Name] = document.querySelector("#"+conf.prefix+"-"+p.Name).value;
            });
            callback(data);
        });
    };

    let _categoryHeader = function(containerId, prefix){
        let tem=[];
        tem.push("<fieldset>");
        tem.push("<legend class='tx-text-20'>");
        tem.push("Category : Header");
        tem.push("</legend>");
        tem.push(_input({Name:"Name",Reference:"text"},prefix));
        tem.push(_input({Name:"Description",Reference:"text"},prefix));
        tem.push(_checkbox({Name:"Rechargeable"},prefix));
        tem.push("</fieldset>");
        document.querySelector(containerId).innerHTML += tem.join("");
    }

    let _categoryForm = function(conf,data,callback){
        document.querySelector(conf.id).innerHTML="";
        _category ={ix:0,types:data.Types,headers:data.Headers};
        let tem=[];
        tem.push(_icon("cc","category","close","red"));
        document.querySelector(conf.id).innerHTML += tem.join("");
        _categoryHeader(conf.id,"cc-header");
        tem = [];
        tem.push("<fieldset>");
        tem.push("<legend class='tx-text-20'>");
        tem.push("Category : Properties");
        tem.push("</legend>");
        tem.push(_input({Name:"Name",Reference:"text"},"category"));
        tem.push(_selectType({Name:"Type",Reference:"text"},"category",data.Types));
        tem.push(_checkbox({Name:"Nullable",Reference:"text"},"category"));
        tem.push(_checkbox({Name:"Downloadable",Reference:"text"},"category"));
        tem.push(_selectReference({Name:"Reference",Reference:"text"},"category"));
        tem.push(_icon("cc","category","add","orange"));
        tem.push("<div class='w3-card-4 w3-round w3-border tx-text-12 w3-ul tx-margin-left-4'>");
        tem.push("<ul id='create-category-properties' class='w3-ul'>");
        tem.push("</ul></div>");
        tem.push(_button("Save","cc",callback));
        tem.push("</fieldset>");
        document.querySelector(conf.id).innerHTML += tem.join("");
        _category.typeSelect = document.querySelector("#category-Type");
        _category.referenceSelectBox = document.querySelector("#category-Reference-box");
        _category.referenceSelectBox.style.display = "none";
        _category.referenceSelect = document.querySelector("#category-Reference");
        _category.build = document.querySelector("#create-category-properties");
        _category.typeSelect.onchange =()=>{
             _category.referenceSelectBox.style.display = "none";
            let ty = _category.typeSelect.options[_category.typeSelect.selectedIndex].text;
            if (ty=="List" || ty=="Set"){
                let ref =[];
                for(const k of Object.keys(data.Types)){
                    ref.push("<option value='"+k+"'>"+k+"</option>");
                }
                _category.referenceSelect.innerHTML = ref.join("");
                _category.referenceSelectBox.style.display = "block";
            }
            else{
                let cat = data.Types[ty];
                if(cat.Type == "enum"){
                    let ref =[];
                    cat.Values.forEach(v=>{
                        ref.push("<option value='"+cat.Name+"'>"+v.Name+"</option>");
                    });
                    _category.referenceSelect.innerHTML = ref.join("");
                    _category.referenceSelectBox.style.display = "block";        
                }
            }
        };
        _category.referenceSelect.onchange=()=>{
            console.log(_category.referenceSelect.options[_category.referenceSelect.selectedIndex].getAttribute("value"));
        };
        _eventWithId("#cc-category-close",()=>{
            document.querySelector(conf.id).style.display='none';
        });
         _eventWithId("#cc-category-add",()=>{
            let cn = document.querySelector("#category-Name").value;
            let nullable = document.querySelector("#category-Nullable").checked;
            let downloadable = document.querySelector("#category-Downloadable").checked;
            let tp = data.Types[_category.typeSelect.options[_category.typeSelect.selectedIndex].text];
            if( tp.Type == "list" || tp.Type =="set"){
               v = _category.referenceSelect.options[_category.referenceSelect.selectedIndex].getAttribute("value");
               let cat = data.Types[v];
               let item = cn+":"+tp.Type+"&lt"+"category:"+cat.Name+"&gt";
               let id = "c-entry"+_category.ix;
               _addItem(id,item);
               _category.ix++;
               _category[id]={Name:cn,Type:tp.Type,Reference:"category:"+cat.Name,Nullable:nullable,Downloadable:downloadable}
               return; 
            }
            if(tp.Type== "enum"){
               v = _category.referenceSelect.options[_category.referenceSelect.selectedIndex].getAttribute("value");
               let cat = data.Types[v];
               let item = (cn+":"+tp.Type+"&lt"+cat.Name+"&gt");
               let id = "c-entry"+_category.ix;
               _addItem(id,item);
               _category.ix++;
               _category[id]={Name:cn,Type:tp.Type,Reference:cat.Name,Nullable:nullable,Downloadable:downloadable}
               return;
            }
            if(tp.Type =="category"){
                let item = (cn+":category"+"&lt"+tp.Name+"&gt");
                let id = "c-entry"+_category.ix;
               _addItem(id,item);
               _category.ix++;
               _category[id]={Name:cn,Type:tp.Type,Reference:tp.Type+":"+tp.Name,Nullable:nullable,Downloadable:downloadable}   
                return;
            }
            let  item = cn+":"+tp.Type+"&lt"+tp.Name+"&gt";
            let id = "c-entry"+_category.ix;
            _addItem(id,item);
            _category.ix++;
            _category[id]={Name:cn,Type:tp.Type,Reference:tp.Reference,Nullable:nullable,Downloadable:downloadable}
                
        }); 
         _eventWithId("#cc-Save",()=>{
            let nm = document.querySelector("#cc-header-Name").value;
            let des = document.querySelector("#cc-header-Description").value;
            let rc = document.querySelector("#cc-header-Rechargeable").checked;
            _category.Name = nm;
            _category.Description = des;
            _category.Rechargeable = rc;
            _category.Scope = _task.Name;
            callback(_category);
        });      
    };

    let _addItem = function(id,item){
        let li = "<li id='"+id+"'>"+item+"<scan class='w3-right category-entry-remove' "+"category-entry-id='"+id+"'><i class='material-symbols-outlined tx-red-icon-24'>remove</i></span></li>";
        _category.build.innerHTML += li;
        document.querySelectorAll(".category-entry-remove").forEach(a=>{
            a.onclick = ()=>{
                let removeId = a.getAttribute("category-entry-id");
                document.querySelector("#"+removeId).style.display="none";
                delete _category[removeId];
            };
        });
    };

    let _instanceHeader = function(conf,header){
        let tem=[];
        tem.push("<fieldset>");
        tem.push("<legend class='tx-text-20'>");
        tem.push(conf.category+" : Header");
        tem.push("</legend>");
        header.Properties.forEach(c=>{
            tem.push(_input(c,conf.prefix));
        });
        tem.push("</fieldset>");
        document.querySelector(conf.id).innerHTML += tem.join("");
    }
    
    let _instanceForm = function(conf,data,callback){
        document.querySelector(conf.id).innerHTML="";
        let tem=[];
        tem.push(_icon("ins","category","close","red"));
        document.querySelector(conf.id).innerHTML += tem.join("");
        _instanceHeader({id:conf.id,prefix:"ins-header",category:data.Name},_category.headers[data.Scope]);
        tem = [];
        tem.push("<fieldset>");
        tem.push("<legend class='tx-text-20'>");
        tem.push(data.Name+" : Properties");
        tem.push("</legend>");
        data.Properties.forEach(p=>{
            if(p.Type =="boolean"){
                tem.push(_checkbox(p,"ins"));
            }else if(p.Type=="enum"){
                tem.push(_selectEnum(p,"ins"));
            }else if(p.Type=="list"){
                tem.push(_selectCategory(p,"ins"));
            }else if(p.Type=="category"){
                tem.push(_selectCategory(p,"ins"));
            }else{
                tem.push(_input(p,conf.prefix));
            }
        })
        tem.push(_button("Save","ins",callback));
        tem.push("</fieldset>");
        document.querySelector(conf.id).innerHTML += tem.join("");
        _eventWithId("#ins-category-close",()=>{
            document.querySelector(conf.id).style.display='none';
        });
        _eventWithId("#ins-Save",()=>{
            callback({});
        });
    }

    return {
        registerTaskChangeListener : _registerTaskChangeListener,
        messageWithId : _messageWithId,
        openWithId : _openWithId,
        closeWithId : _closeWithId,
        eventWithId : _eventWithId,
        form : _form,
        taskList : _taskList,
        jobList : _jobList,
        typeList : _typeList,
        categoryList : _categoryList,
        enumForm : _enumForm,
        categoryForm : _categoryForm,
        instanceForm : _instanceForm,
    };
})();