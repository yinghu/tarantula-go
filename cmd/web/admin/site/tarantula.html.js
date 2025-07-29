var Html = (function(){
    
    let _tasks = {};
    let _types = {};
    let _headers ={};
    let _enum ={};
    let _category ={};
    let _task ={};
    let _instance ={};
    let _upload ={};
    
    let _uploadTask = (data,callback)=>{}
    let _taskChanged = (tsk)=>{};
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
    let _currentTask = function(){
        return _task;
    };
    let _closeWithClass = function(cls){
        document.querySelectorAll(cls).forEach(a=>{
            a.style.display = "none";
        });
    }
    let _eventWithId = function(id,callback){
        document.querySelector(id).onclick = callback; 
    };
    let _registerTaskChangeListener = function(listener){
        _taskChanged = listener;
    };
    let _registerUploadTask = function(task){
        _uploadTask = task;
    };
    let _setup = function(tlist){
        _types = tlist.Types;
        _headers = tlist.Headers;
    };
    let _typeList = function(tlist){
        tlist.forEach(e=>{
            e.Type="category";
            _types[e.Name]=e;
        });
    };
    let _enumList = function(elist){
        elist.forEach(e=>{
            e.Type="enum";
            _types[e.Name]=e;
        });
    };
    let _taskList = function(conf,tbar,callback){
        document.querySelector(conf.id).innerHTML="";
        let tem=[];
        tbar.Tasks.forEach(task=>{
            _tasks[task.Name]=task;
            tem.push("<div tx-task-name='");
            tem.push(task.Name+"' ")  
            tem.push("class='w3-display-container  w3-border-bottom w3-border-red tx-content-48 ");
            tem.push("tx-"+conf.prefix+"-action'>");
            tem.push("<div class='w3-display-bottomright w3-padding'><span id='");
            tem.push(conf.prefix+"-icon-"+task.Name+"' class='");
            tem.push(conf.prefix+"-task-icon'>");
            tem.push("<i class='material-symbols-outlined tx-margin-left-8 tx-orange-icon-24'>double_arrow</i></span>");
            tem.push("</div>");
            tem.push("<div class='w3-display-bottomleft w3-padding'>");
            tem.push("<span class='tx-margin-left-8 tx-text-20'>");
            tem.push(task.Name);
            tem.push("</span>");
            tem.push("</div>");
            tem.push("</div>");
        });
        document.querySelector(conf.id).innerHTML = tem.join("");
        _closeWithClass("."+conf.prefix+"-task-icon");
        _openWithId("#"+conf.prefix+"-icon-"+tbar.Tasks[0].Name);
        document.querySelectorAll(".tx-"+conf.prefix+"-action").forEach(a=>{
            a.onclick = ()=>{
                _closeWithClass("."+conf.prefix+"-task-icon");
                let tn = a.getAttribute("tx-task-name");
                callback(tn);
                _openWithId("#"+conf.prefix+"-icon-"+tn);    
            };
        });
    };
    let _jobList = function(conf,tn,callback){
        let task = _tasks[tn];
        _task.Name = tn;
        _task.ScopeSequence = task.ScopeSequence;
        _task.Icon = task.Icon;
        _taskChanged(_task);
        document.querySelector(conf.id).innerHTML="";
        let tem=[];
        task.Jobs.forEach(job=>{
            tem.push("<span tx-job-name='");
            tem.push(job.Callback+"' ");
            tem.push("class='w3-bar-item w3-green w3-tag tx-text-20 tx-padding-button tx-margin-right-4 tx-margin-bottom-4 ");
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
            if(c.ScopeSequence == _task.ScopeSequence){
                tem.push("<span tx-category-id='"+c.Id+"' ");
                tem.push("class='w3-bar-item w3-green w3-tag tx-text-20 tx-padding-button tx-margin-right-4 tx-margin-bottom-4 tx-"+conf.prefix+"-opt'>");
                tem.push(c.Name);
                tem.push("</span>");
            }    
        });
        document.querySelector(conf.id).innerHTML = tem.join("");
        document.querySelectorAll(".tx-"+conf.prefix+"-opt").forEach(a=>{
            a.onclick = ()=>{
                callback(a.getAttribute("tx-category-id"));    
            };
        });
    };
    let _categorySelect = function(conf,clist,callback){
        let ctn = document.querySelector(conf.id);
        ctn.innerHTML = "";
        let tem =[];
        tem.push("<option value='' disabled selected>select to load</option>");
        clist.forEach(c=>{
            if(c.ScopeSequence==_task.ScopeSequence){
                tem.push("<option tx-category-id='"+c.Id+"'>");
                tem.push(c.Name);
                tem.push("</option>");
            }    
        });
        ctn.innerHTML = tem.join("");
        ctn.onchange = ()=>{
            callback(ctn.options[ctn.selectedIndex].text);
        };
    };
    let _instanceList = function(conf,clist,callback){
        document.querySelector(conf.id).innerHTML = "";
        let tem =[];
        clist.forEach(c=>{
            tem.push("<span tx-instance-id='"+c.ItemId+"' ");
            tem.push("class='w3-bar-item w3-green w3-tag tx-text-20 tx-padding-button tx-margin-right-4 tx-margin-bottom-4 tx-"+conf.prefix+"-opt'>");
            tem.push(c.ConfigurationName+" : "+c.ConfigurationVersion);
            tem.push("</span>");    
        });
        document.querySelector(conf.id).innerHTML = tem.join("");
        document.querySelectorAll(".tx-"+conf.prefix+"-opt").forEach(a=>{
            a.onclick = ()=>{
                callback(a.getAttribute("tx-instance-id"));    
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
        tem.push(_button("Save","ee"));
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

    let _textarea = function(prop,prefix){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<label class='tx-text-12'>");
        tem.push(_caption(prop.Name));
        tem.push("</label>");
        tem.push("<textarea id='");
        tem.push(prefix+"-"+prop.Name+"' ");
        tem.push("class='w3-round w3-border tx-text-16' type='");
        tem.push(prop.Reference+"'></textarea>");
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
        tem.push(prefix+"-delete' ");
        tem.push("class='w3-left w3-tag w3-green tx-text-18 tx-padding-button tx-margin-top-8 tx-margin-right-8 tx-hidden'>");
        tem.push("Delete");
        tem.push("</span>");
        tem.push("<span id='");
        tem.push(prefix+"-register' ");
        tem.push("class='w3-left w3-tag w3-green tx-text-18 tx-padding-button tx-margin-top-8 tx-margin-right-8 tx-hidden'>");
        tem.push("Register");
        tem.push("</span>");

        tem.push("<span id='");
        tem.push(prefix+"-"+name+"' ");
        tem.push("class='w3-right w3-tag w3-green tx-text-18 tx-padding-button tx-margin-top-8 tx-margin-right-8'>");
        tem.push(name);
        tem.push("</span>");
        tem.push("</div>");
        return tem.join("");
    };

    let _uploadFile = function(prop,prefix){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<label class='tx-text-12'>");
        tem.push(_caption(prop.Name));
        tem.push("</label>");
        tem.push("<span tx-property-name='");
        tem.push(prop.Name+"' id='");
        tem.push(prefix+"-"+prop.Name+"-upload");    
        tem.push("' class='w3-right tx-hidden tx-margin-right-4 ");
        tem.push(prefix+"-upload'>");
        tem.push("<i class='material-symbols-outlined tx-orange-icon-20'>upload</i></span>");
        tem.push("<input id='");
        tem.push(prefix+"-"+prop.Name+"' ");
        tem.push("class='w3-input w3-round w3-border tx-text-16 ");
        tem.push(prefix+"-file-ready' type='file'/>");
        tem.push("</div>");
        return tem.join("");  
    };

    let _icon = function(prefix,n,icon,color){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<span id='");
        tem.push(prefix+"-"+n+"-"+icon+"' class='w3-right'>");
        tem.push("<i class='material-symbols-outlined tx-margin-right-8 tx-"+color+"-icon-24'>"+icon+"</i></span>");
        tem.push("</div>");
        return tem.join("");
    };


    let _selectEnum = function(prop,prefix){
        let t = _category.types[prop.Reference];
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<label class='tx-text-12'>");
        tem.push(_caption(prop.Name));
        tem.push("</label>");
        tem.push("<select id='");
        tem.push(prefix+'-'+prop.Name);
        tem.push("' name='chs' class='w3-round w3-border tx-text-16 w3-select tx-margin-left-4'>");
        t.Values.forEach(e=>{
            tem.push("<option value='"+e.Value+"'>"+e.Name+"</option>");    
        });
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
        for(const [k,v] of Object.entries(types)){
            if(v.Type != "category" || v.ScopeSequence <= _task.ScopeSequence){ 
                tem.push("<option value='"+k+"'>"+k+"</option>");
            }
        }
        tem.push("</select>");
        tem.push("</div>");
        return tem.join("");  
    };

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
    };

    let _selectCategory = function(prop,prefix,listing){
        let tem=[];
        tem.push("<div class='w3-panel'>");
        tem.push("<label class='tx-text-12'>");
        tem.push(_caption(prop.Name));
        tem.push("</label>");
        tem.push("<span tx-property-category='");
        tem.push(prop.Reference.split(":")[1]+"' tx-property-name='");
        tem.push(prop.Name);    
        tem.push("' class='w3-right tx-margin-right-4 ");
        tem.push(prefix+"-load'>");
        tem.push("<i class='material-symbols-outlined tx-orange-icon-20'>add</i></span>");
        tem.push("<select id='");
        tem.push(prefix+'-'+prop.Name);
        tem.push("' name='"+prop.Name+"' class='w3-round w3-border tx-text-16 w3-select tx-margin-left-4");
        if (listing){
            tem.push(" "+prefix+"-select");
        }
        tem.push("'>");
        tem.push("<option value='' disabled selected>click plus to load</option>");
        tem.push("</select>");
        tem.push("</div>");
        if (listing){
            tem.push("<div class='w3-card-4 w3-round w3-border tx-text-12 w3-ul tx-margin-left-4'>");
            tem.push("<ul id='");
            tem.push(prefix+"-select-"+prop.Name);
            tem.push("' class='w3-ul'>");
            tem.push("</ul></div>");
        }
        return tem.join("");  
    };

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
        tem.push(_button(category.Name,conf.prefix));
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
        if(_task.ScopeSequence != 3){
            document.querySelector("#"+prefix+"-Rechargeable").disabled = true;
        }
    };

    let _categoryForm = function(conf,callback){
        document.querySelector(conf.id).innerHTML="";
        _category = {ix:0,types:_types,headers:_headers};
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
        tem.push(_selectType({Name:"Type",Reference:"text"},"category",_category.types));
        tem.push(_checkbox({Name:"Nullable",Reference:"text"},"category"));
        tem.push(_selectReference({Name:"Reference",Reference:"text"},"category"));
        tem.push(_icon("cc","category","add","orange"));
        tem.push("<div class='w3-card-4 w3-round w3-border tx-text-12 w3-ul tx-margin-left-4'>");
        tem.push("<ul id='create-category-properties' class='w3-ul'>");
        tem.push("</ul></div>");
        tem.push(_button("Save","cc"));
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
                for(const [k,v] of Object.entries(_category.types)){
                    if(v.Type=="category"){
                        ref.push("<option value='"+k+"'>"+k+"</option>");
                    }
                }
                _category.referenceSelect.innerHTML = ref.join("");
                _category.referenceSelectBox.style.display = "block";
            }
            else{
                let cat = _category.types[ty];
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
            let tp = _category.types[_category.typeSelect.options[_category.typeSelect.selectedIndex].text];
            if( tp.Type == "list" || tp.Type =="set"){
               v = _category.referenceSelect.options[_category.referenceSelect.selectedIndex].getAttribute("value");
               let cat = _category.types[v];
               let item = cn+":"+tp.Type+"&lt"+"category:"+cat.Name+"&gt";
               let id = "c-entry"+_category.ix;
               _addItem(id,item);
               _category.ix++;
               _category[id]={Name:cn,Type:tp.Type,Reference:"category:"+cat.Name,Nullable:nullable}
               return; 
            }
            if(tp.Type== "enum"){
               v = _category.referenceSelect.options[_category.referenceSelect.selectedIndex].getAttribute("value");
               let cat = _category.types[v];
               let item = (cn+":"+tp.Type+"&lt"+cat.Name+"&gt");
               let id = "c-entry"+_category.ix;
               _addItem(id,item);
               _category.ix++;
               _category[id]={Name:cn,Type:tp.Type,Reference:cat.Name,Nullable:nullable}
               return;
            }
            if(tp.Type =="category"){
                let item = (cn+":category"+"&lt"+tp.Name+"&gt");
                let id = "c-entry"+_category.ix;
               _addItem(id,item);
               _category.ix++;
               _category[id]={Name:cn,Type:tp.Type,Reference:tp.Type+":"+tp.Name,Nullable:nullable}   
                return;
            }
            let  item = cn+":"+tp.Type+"&lt"+tp.Name+"&gt";
            let id = "c-entry"+_category.ix;
            _addItem(id,item);
            _category.ix++;
            _category[id]={Name:cn,Type:tp.Type,Reference:tp.Reference,Nullable:nullable}
                
        }); 
         _eventWithId("#cc-Save",()=>{
            let nm = document.querySelector("#cc-header-Name").value;
            let des = document.querySelector("#cc-header-Description").value;
            let rc = document.querySelector("#cc-header-Rechargeable").checked;
            _category.Name = nm;
            _category.Description = des;
            _category.Rechargeable = rc;
            _category.Scope = _task.Name;
            _category.ScopeSequence = _task.ScopeSequence;
            callback(_category);
        });      
    };

    let _addItem = function(id,item){
        let li = "<li id='"+id+"'>"+item+"<scan class='w3-right tx-text-16 category-entry-remove' "+"category-entry-id='"+id+"'><i class='material-symbols-outlined tx-red-icon-24'>remove</i></span></li>";
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
    };

    let _readInstanceHeader = function(conf){
        _instance.header.Properties.forEach(c=>{
            _instance.save[c.Name]=document.querySelector("#"+conf.prefix+"-"+c.Name).value;
        });
    };

    let _readInstance = function(conf){
        _instance.save.header = {};
        _instance.category.Properties.forEach(c=>{
            if(c.Type=="boolean"){
                _instance.save.header[c.Name]=document.querySelector("#"+conf.prefix+"-"+c.Name).checked;
            }else if(c.Type=="dateTime"){
                _instance.save.header[c.Name]=new Date(Date.parse(document.querySelector("#"+conf.prefix+"-"+c.Name).value)).toISOString();
            }else if(c.Type=="enum"){
                let ctn = document.querySelector("#"+conf.prefix+"-"+c.Name);
                _instance.save.header[c.Name]= ctn.options[ctn.selectedIndex].value/1;    
            }else if(c.Type=="list" ||c.Type=="set"){
                _instance.save.application[c.Name]=_instance.build[c.Name];    
            }else if(c.Type=="category"){
                let ctn = document.querySelector("#"+conf.prefix+"-"+c.Name);
                let tem =[];
                tem.push(ctn.options[ctn.selectedIndex].value);
                _instance.save.application[c.Name]= tem;
            }else{
                if (c.Reference == "number"){
                    _instance.save.header[c.Name]=document.querySelector("#"+conf.prefix+"-"+c.Name).value/1;    
                }
                else{
                    if(c.Reference == "file"){
                        _instance.save.header[c.Name] = _instance.build[c.Name];
                    }else{
                        _instance.save.header[c.Name]=document.querySelector("#"+conf.prefix+"-"+c.Name).value;
                    }
                }
            }
        });
    };

    let _addCategories = function(containerId,cats){
        let cnt = document.querySelector(containerId);
        cnt.innerHTML = "";
        let tem =[];
        tem.push("<option value='0' disabled selected>drop down to choose</option>");
        cats.forEach(c=>{
            tem.push("<option value='"+c.ItemId+"'>"+c.ConfigurationCategory+":"+c.ConfigurationName+"/"+c.ConfigurationVersion+"</option>");     
        });
        cnt.innerHTML = tem.join("");
    };

    let _addInstance = function(item){
        let tem=[];
        tem.push("<li id='"+item.id+"' tx-select-name='"+item.selected+"'>");
        tem.push(item.name);
        tem.push("<scan class='w3-right instance-entry-remove' "+"instance-entry-id='"+item.id+"'><i class='material-symbols-outlined tx-red-icon-24'>remove</i></span></li>");
        item.build.innerHTML += tem.join("");
        document.querySelectorAll(".instance-entry-remove").forEach(a=>{
            a.onclick = ()=>{
                let removeId = a.getAttribute("instance-entry-id");
                let removed = document.querySelector("#"+removeId);
                removed.style.display="none";
                let selected = removed.getAttribute("tx-select-name").split(":");
                let ids = _instance.build[selected[0]].filter(id=> id!==selected[1]);
                _instance.build[selected[0]]=ids;
            };
        });
    };
    
    let _instanceForm = function(conf,data,save,load){
        if (_headers[data.Scope] === undefined){
            _instance = {ix:0,save:{header:{},application:{}},category:data,header:_headers["Application"],build:{}};
        }else{
            _instance = {ix:0,save:{header:{},application:{}},category:data,header:_headers[data.Scope],build:{}};
        }
        document.querySelector(conf.id).innerHTML="";
        let tem=[];
        tem.push(_icon("ins","category","close","red"));
        document.querySelector(conf.id).innerHTML += tem.join("");
        _instanceHeader({id:conf.id,prefix:"ins-header",category:data.Name},_instance.header);
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
            }else if(p.Type=="list" || p.Type=="set"){
                _instance.build[p.Name]=[];
                tem.push(_selectCategory(p,"ins",true));
            }else if(p.Type=="category"){
                tem.push(_selectCategory(p,"ins",false));
            }else{
                if(p.Reference == "textarea"){
                    tem.push(_textarea(p,"ins"));
                }else if(p.Reference =="file"){
                    tem.push(_uploadFile(p,"ins"));
                }
                else{
                    tem.push(_input(p,"ins"));
                }
            }
        });
        tem.push(_button(conf.name,"ins"));
        tem.push("</fieldset>");
        document.querySelector(conf.id).innerHTML += tem.join("");
        document.querySelectorAll(".ins-load").forEach(a=>{
            a.onclick = ()=>{
                let cid = "#ins-"+a.getAttribute("tx-property-name");
                load(a.getAttribute("tx-property-category"),cats=>{
                    _addCategories(cid,cats);
                });
            };
        });
        document.querySelectorAll(".ins-select").forEach(a=>{
            a.onchange = ()=>{
                let pname = a.getAttribute("name");
                let selected = document.querySelector("#ins-"+pname);
                let sid = selected.options[selected.selectedIndex].value;
                let notExisting = true;
                _instance.build[pname].forEach(i=>{
                    if(i===sid && notExisting){
                        notExisting = false;
                    }
                });
                if (notExisting){
                    _instance.build[pname].push(sid);
                    let item ={selected:pname+":"+sid,id:"ins-"+pname+"-"+_instance.ix,name:sid,build:document.querySelector("#ins-select-"+pname)};
                    _addInstance(item);
                    _instance.ix++;
                }
            };
        });
        document.querySelectorAll(".ins-upload").forEach(a=>{
            a.onclick = ()=>{
                _uploadTask(_upload,(resp)=>{
                    pname = a.getAttribute("tx-property-name");
                    _closeWithId("#ins-"+pname+"-upload");
                    _instance.build[pname]=resp.file;
                });
            };
        });
        document.querySelectorAll(".ins-file-ready").forEach(a=>{
            a.onchange = (e)=>{
                _upload = {};
                let _fd = e.target.files[0];
                _upload.file = _fd;
                _upload.name = _fd.name;
                let reader = new FileReader();
                reader.onloadend = ()=>{
                    _upload.data = reader.result;
                    _upload.ready = true;
                    document.querySelector("#"+a.getAttribute("id")+"-upload").style.display = "block";
                } ;
                reader.readAsArrayBuffer(_fd);
            };
        });
        _eventWithId("#ins-category-close",()=>{
            document.querySelector(conf.id).style.display='none';
        });
        _eventWithId("#ins-"+conf.name,()=>{
            _readInstanceHeader({prefix:"ins-header"});
            _instance.save.ConfigurationCategory = data.Name;
            _readInstance({prefix:"ins"});
            save(_instance.save);
        });
    };

    let _populateInstance = function(ins,load,deleteCall,registerCall){
        _instance.header.Properties.forEach(p=>{
            document.querySelector("#ins-header-"+p.Name).value = ins[p.Name];
        });
        _instance.category.Properties.forEach(p=>{
            let ctn = document.querySelector("#ins-"+p.Name);
            if(p.Type =="boolean"){
                ctn.checked = ins.header[p.Name];    
            }else if(p.Type=="enum"){
                ctn.options[ins.header[p.Name]].selected = true;    
            }else if(p.Type=='dateTime'){
                const date = new Date(Date.parse(ins.header[p.Name]));
                let parts = date.toISOString().split('T');
                ctn.value = parts[0]+"T"+parts[1].split(".")[0];    
            }else if(p.Type=="list" || p.Type=="set"){
                let cid = "#ins-"+p.Name;
                load(p.Reference.split(":")[1],cats=>{
                    _addCategories(cid,cats);
                    ins.application[p.Name].forEach(c=>{    
                        if(c!==""){
                            _instance.build[p.Name].push(c);
                            let item ={selected:p.Name+":"+c,id:"ins-"+p.Name+"-"+_instance.ix,name:c,build:document.querySelector("#ins-select-"+p.Name)};
                            _addInstance(item);
                            _instance.ix++;
                        }
                    });
                });
            }else if(p.Type=="category"){
                let cid = "#ins-"+p.Name;
                let ctn = document.querySelector(cid);
                load(p.Reference.split(":")[1],cats=>{
                    _addCategories(cid,cats);
                    for(let i=0;i<ctn.options.length;i++){
                        if(ctn.options[i].value === ins.application[p.Name][0]){
                            ctn.options[i].selected = true;
                            break;
                        }
                    }
                });   
            }else{
                if(p.Reference!="file"){
                    ctn.value = ins.header[p.Name];
                }    
            }
        });
        _eventWithId("#ins-delete",()=>{
            deleteCall(ins.ItemId);     
        });
        _openWithId("#ins-delete");
        _eventWithId("#ins-register",()=>{
            registerCall(ins.ItemId);     
        });
        _openWithId("#ins-register");
    };

    _editForm = function(conf,clist,edit,preview){
        document.querySelector(conf.id).innerHTML = "";
        let tem =[];
        tem.push("<fieldset>");
        tem.push("<legend class='tx-text-20'>");
        tem.push("Edit Form");
        tem.push("</legend>");
        tem.push(_icon(conf.prefix,"category","close","red"));
        
        clist.forEach(c=>{
            if(c.ScopeSequence == _task.ScopeSequence){
                tem.push("<div class='w3-panel w3-padding w3-border-bottom w3-border-red'>");
                tem.push("<span class='w3-green w3-tag tx-text-20 tx-padding-button'>");
                tem.push(c.Name);
                tem.push("</span>");
                tem.push("<span tx-category-id='"+c.Id+"' ");
                tem.push("class='w3-green w3-right w3-tag tx-text-20 tx-padding-button tx-margin-left-8 tx-"+conf.prefix+"-delete-opt'>Delete");
                tem.push("</span>");
                tem.push("<span tx-category-id='"+c.Id+"' ");
                tem.push("class='w3-green w3-right w3-tag tx-text-20 tx-padding-button tx-margin-left-8 tx-"+conf.prefix+"-preview-opt'>Preview");
                tem.push("</span>");
                tem.push("<span tx-category-id='"+c.Id+"' ");
                tem.push("class='w3-green w3-right w3-tag tx-text-20 tx-padding-button tx-margin-left-8 tx-"+conf.prefix+"-register-opt'>Register");
                tem.push("</span>");
                tem.push("<span tx-category-id='"+c.Id+"' ");
                tem.push("class='w3-green w3-right w3-tag tx-text-20 tx-padding-button tx-margin-left-8 tx-"+conf.prefix+"-release-opt'>Release");
                tem.push("</span>");
                tem.push("</div>");
            }    
        });
        tem.push("</fieldset>");
        document.querySelector(conf.id).innerHTML = tem.join("");
        document.querySelectorAll(".tx-"+conf.prefix+"-delete-opt").forEach(a=>{
            a.onclick = ()=>{
                edit(a.getAttribute("tx-category-id"));    
            };
        });
        document.querySelectorAll(".tx-"+conf.prefix+"-preview-opt").forEach(a=>{
            a.onclick = ()=>{
                preview(a.getAttribute("tx-category-id"));    
            };
        });
        _eventWithId("#"+conf.prefix+"-category-close",()=>{
            document.querySelector(conf.id).style.display='none';
        });       
    };

    return {
        registerTaskChangeListener : _registerTaskChangeListener,
        registerUploadTask : _registerUploadTask,
        messageWithId : _messageWithId,
        openWithId : _openWithId,
        closeWithId : _closeWithId,
        eventWithId : _eventWithId,
        form : _form,
        currentTask : _currentTask,
        taskList : _taskList,
        jobList : _jobList,
        setup : _setup,
        typeList : _typeList,
        enumList : _enumList,
        categoryList : _categoryList,
        categorySelect : _categorySelect,
        instanceList : _instanceList,
        enumForm : _enumForm,
        categoryForm : _categoryForm,
        instanceForm : _instanceForm,
        populateInstance : _populateInstance,
        editForm : _editForm,
    };
})();