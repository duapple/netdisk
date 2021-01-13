/* 退出
*  @params
*  @return
*/
function logout() {
    window.location.href = logout_href;
}

/* 校验文件名
*  @params
*      fileName 文件名
*  @return
*      [boolean] true:是;false:否
*/
function validateFileName(fileName) {
    var reg = new RegExp('[\\\\/:*?\"<>|]');
    if (reg.test(fileName)) {
        return false;
    }
    return true;
}

/* 阻止冒泡
*  @params
*      e  
*  @return 
*/
function stopPropagation(e) {
    e = e || window.event;
    if(e.stopPropagation) { //W3C阻止冒泡方法
        e.stopPropagation();
    } 
    else {
        e.cancelBubble = true; //IE阻止冒泡方法
    }
}

/* 转换字节
*  @params
*      bytes [number] 字节数
*  @return
*      [string] 转换后的字符串
*/
function bytesToSize(bytes) {
    if (bytes === 0) return '0 B';
    let k = 1024,
    sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'],
    i = Math.floor(Math.log(bytes) / Math.log(k));
    return (bytes / Math.pow(k, i)). toFixed(2) + ' ' + sizes[i];
 }

 /* 转换字节
*  @params
*      e 父节点
*  @return
*/
 function deleteChild(e) { 
    let child = e.lastElementChild;  
    while (child) { 
        e.removeChild(child); 
        child = e.lastElementChild; 
    } 
} 
 
/* 文件查看刷新
*  @params
*      ret [string] 当前文件夹
*  @return 
*/
function fileShow(ret) {
    let container = document.getElementsByClassName("content")[0],  //文件目录表格所在的区域
        menu = document.getElementsByClassName("menu")[0],  //右键的菜单
        file_system = document.getElementsByClassName("file_system")[0];  //当前路径
    let index_data = "{\"Opt\"" + ":" + "0" + "," + "\"DirName\"" + ":" + "[\"" + ret + "\"]" + "}";
    let json_str,  //目录数据（字符串）
        message,   //转化为json格式
        key_word = [],  //存储所有文件夹，便于点击查看时获取当前点击的文件夹
        index = 0;  //key_word的索引

    // 显示当前目录
    $.ajax(
        {
            url: home_rpc,
            data: index_data,
            type: "POST",
            async: false,
            success: function (data)
            {
				if(data){
                    json_str = data;
                    return true;
                }
                else {
                    alert(data.description);
                    return false;
                }
            },
            error: function () {
                alert("Network error!")
            }
        });

    //创建tbody用于存储目录结构
    let json_table = document.getElementsByTagName("table")[0];
    json_table.style.borderCollapse = "collapse";
    let json_tbody = document.createElement("tbody"),
        json_tr = document.createElement("tr");
    json_table.appendChild(json_tbody);

    let table = document.getElementsByTagName("table")[0],
        infoList = document.getElementsByTagName("tbody"),
        infoLen = infoList.length,
        tobodyList = [];
    for (let i = infoLen - 1; i > 0; i--) {
        tobodyList[i] = document.getElementsByTagName("tbody")[i];
        table.removeChild(tobodyList[i]);
    }
    //此处必须创建tbody，否则无法加入到table中
    let infos = document.createElement("tbody");

    message = json_str;
    console.log(message);
    // 遍历对象的属性
    for (let prop in message) {
        if(prop !== "CurrentDir") {
            // 遍历对象的某个属性		
            for (let num in message[prop]) {
                let json_tr = document.createElement("tr");
                // 每行第一列为checkbox
                let checkbox_td = document.createElement("td"),
                    label = document.createElement("label"),
                    input = document.createElement("input"),
                    i = document.createElement("i");

                checkbox_td.className = "tdwidthbox";
                label.className = "checklabel";
                input.className = "checkbox";
                input.type = "checkbox";
                i.className = "check";

                label.appendChild(input);
                label.appendChild(i);
                checkbox_td.appendChild(label);
                json_tr.appendChild(checkbox_td);

                // 遍历对象的某个属性的数组
                for (let key in message[prop][num]) {
                    let json_td = document.createElement("td");
                    //获取键key
                    let td_txt = document.createTextNode(message[prop][num][key]);
                    // 大小和修改时间直接放入td
                    if (key != "FileName" && key != "DirName") {
                        json_td.appendChild(td_txt);
                    }
                    // 文件名的td添加样式
                    else {
                        let json_span = document.createElement("span");
                        let json_div = document.createElement("div");
                        json_div.className = "file_name";

                        json_span.appendChild(td_txt);
                        json_div.appendChild(json_span);
                        json_td.appendChild(json_div);
                    }
                    json_tr.appendChild(json_td);
                    infos.appendChild(json_tr);
                    json_tr.className = "trstyle";

                    if (key == "DirName" || key == "FileName") {
                        // 存储文件夹和文件的值放入key_word
                        key_word[index] = message[prop][num][key];
                        index++;
                        json_td.className = "tdwidth1";

                        // 添加文件的图标
                        let label_file = document.createElement("label"),
                            i_file = document.createElement("i");
                        label_file.className = "dir_label";
                        label_file.appendChild(i_file);
                        json_td.appendChild(label_file);

                        // 添加操作图标
						let div_icon = document.createElement("div"),
                            icon_share = document.createElement("i"),
                            icon_download = document.createElement("i"),
                            icon_more = document.createElement("i");
                        div_icon.className = "div_icon";
                        icon_share.className = "icon_share";
                        icon_download.className = "icon_download";
                        icon_more.className = "icon_more";
                        div_icon.appendChild(icon_share);
                        div_icon.appendChild(icon_download);
                        div_icon.appendChild(icon_more);
                        json_td.appendChild(div_icon);

                        if (key == "DirName") {  //文件夹的图标
                            i_file.className = "dir_i";
                        }
                        else{  //文件的图标
							i_file.className = "file_i";
							let arr = message[prop][num][key].split('.');
							let file_type = arr[arr.length-1]; //文件类型
							switch(file_type) {
								case "jpeg":
									i_file.className = "jpeg_i";
									break;
								case "jpg":
									i_file.className = "jpg_i";
									break;
								case "mp3":
									i_file.className = "mp3_i";
									break;
								case "mp4":
									i_file.className = "mp4_i";
									break;
								case "pdf":
									i_file.className = "pdf_i";
									break;
								case "png":
									i_file.className = "png_i";
									break;
								case "ppt":
									i_file.className = "ppt_i";
									break;
								case "rar":
									i_file.className = "rar_i";
									break;
								case "zip":
									i_file.className = "rar_i";
									break;
								case "txt":
									i_file.className = "txt_i";
									break;
								case "doc":
									i_file.className = "word_i";
									break;
								case "docx":
										i_file.className = "word_i";
										break;
								case "xls":
									i_file.className = "xls_i";
                                    break;
                                case "xlsx":
									i_file.className = "xls_i";
									break;
								default:
									i_file.className = "other_i";
									break;
							}
						}
                    }
                    else if (key == "Size") {
                        json_td.className = "tdwidth2";
                    }
                    else if (key == "ModTime") {
                        json_td.className = "tdwidth3";
                    }
                }
            }
        }
        else{
            deleteChild(file_system);
            current_dir = message[prop].slice(0, -1).split('/'); //以数组的形式存储路径
			for(let i = current_dir.length - 1; i >= 2 ; i--) {
				let a = document.createElement("a"),
					span = document.createElement("span");
				if(current_dir[i] === username) {
					a.innerText = "全部文件";
				}
				else{
					a.innerText = current_dir[i];
				}
				a.className = "file_system_a";
				span.innerText = '>';
				span.className = "file_system_span";
				file_system.insertBefore(a, file_system.children[0]);
				if(i !== (current_dir.length - 1)) {
					file_system.insertBefore(span, file_system.children[1]);
				}
			}
		}
    }
    table.appendChild(infos);

    // 左键点击表格某一行
	let trList = document.getElementsByTagName("tr"), //每一行文件
        checkList = document.getElementsByClassName("checkbox"),  //文件左侧的选择框
        fileList = document.getElementsByClassName("file_name"), //文件名集合
        more_show = document.getElementsByClassName("more")[0], //更多按钮
        trLen = trList.length,
        checkLen = checkList.length,
        fileLen = fileList.length,
        lastIndex_leftBtn = 0,  //左键的上一次点击
        lastIndex_rightBtn = 0;  //右键的上一次点击
    (function () {
        for (let i = 0; i < fileLen; i++) {
            fileList[i].onclick = function (e) {
                stopPropagation(e); //阻止冒泡
				this.index = i + 1; //去掉第一行
				// 清除所有选中框的样式
				clearBox();
				more_show.style.display = "block"; //显示更多按钮
                // 清除上一次右键点击的样式
                trList[lastIndex_rightBtn].style.background = "none";
                checkList[lastIndex_rightBtn].checked = false;
                if (!(trList[this.index])) {
                    return;
                }
                else {
                    trList[lastIndex_leftBtn].style.background = "none";
                    checkList[lastIndex_leftBtn].checked = false;
                    // 添加背景颜色
                    trList[this.index].style.background = "#e8f6fd";
                    // 选中方框
                    checkList[this.index].checked = true;
                    trList[this.index].isClick = true;
                    lastIndex_leftBtn = this.index; //保存当前的index
                }
            }
        }
    })();

    // 鼠标划过第一行不变换背景色
	trList[0].onmousemove = function () {
		trList[0].style.background = "none";
    }
    
    // 鼠标停留/离开时显示操作图标
	let iconList = document.getElementsByClassName("div_icon");
	(function () {
        for (let i = 1; i < trList.length; i++) {
            trList[i].onmouseenter = function (e) {
				stopPropagation(e); //阻止冒泡
				iconList[i-1].style.display = "block";
			}
			trList[i].onmouseleave = function (e) {
				stopPropagation(e); //阻止冒泡
				iconList[i-1].style.display = "none";
            }
        }
	})();

    // 点击选择框
	let labelList = document.getElementsByClassName("checklabel"),
        labelLen = labelList.length;
    (function() {
		for(let i = 1; i < labelLen; i++) {
			labelList[i].onclick = function (e) {
				stopPropagation(e);
				menu.style.display = 'none';
				if(checkList[i].checked) {
					checkList[i].checked = false;
					trList[i].style.background = "none";
					clearMore();
				}
				else{				
					checkList[i].checked = true;
					trList[i].style.background = "#e8f6fd";
					clearMore();
				}
			}
		}
	})();

    //左键点击查看文件
    let	filenameList = [];  //文件名集合
	let tdList = document.getElementsByClassName("tdwidth1"),
		td1Len = tdList.length;
		i_list = [];
	for (let i = 0; i < fileLen; i++) {
		filenameList[i] = fileList[i].getElementsByTagName("span")[0];
	}
	for (let i = 1; i < td1Len; i++) {
		i_list.push(tdList[i].getElementsByTagName("i")[0])
	}
	let nameLen = filenameList.length;
	(function () {
		for (let i = 0; i < nameLen; i++) {
			filenameList[i].onclick = function (e) {
                stopPropagation(e);
				current_file = key_word[i];
				if(i_list[i].className != "dir_i") { //是文件不可进入
					return;
				}
				else{ //文件夹可以进入
                    select_dir = (fileList[i].getElementsByTagName("span")[0]).innerText;  //当前点击的文件名
                    window.localStorage.setItem('select_dir',select_dir);
					fileShow(current_file);
				}
			}
		}
    })();

    // 鼠标双击
    (function () {
        for (let i = 0; i < fileLen; i++) {
            fileList[i].ondblclick = function (e) {
                stopPropagation(e);
                this.index = i + 1;
                clearBox();
                // 清除上一次右键点击的样式
                trList[lastIndex_rightBtn].style.background = "none";
                checkList[lastIndex_rightBtn].checked = false;
                if (!(trList[this.index])) {
                    return;
                }
                else {
                    trList[lastIndex_leftBtn].style.background = "none";
                    checkList[lastIndex_leftBtn].checked = false;
                    // 添加背景颜色
                    trList[this.index].style.background = "#e8f6fd";
                    // 选中方框
                    checkList[this.index].checked = true;
                    trList[this.index].isClick = true;
                    lastIndex_leftBtn = this.index; //保存当前的index
                }
                current_file = key_word[i];
				if(i_list[i].className == "file_i") {
					return;
				}
				else{
                    select_dir = (fileList[i].getElementsByTagName("span")[0]).innerText;  //当前点击的文件名
                    window.localStorage.setItem('select_dir',select_dir);
					fileShow(current_file);
				}
            }
        }
    })();

    // 右键文件弹出菜单
	(function() {
    	for(let i = 0; i < fileLen; i++) {
            fileList[i].index = i;  //自定义属性index保存索引
            fileList[i].isClick = false;   //定义点击开关
    		fileList[i].onmousedown = function(e) {
                stopPropagation(e);
				// 右键弹出菜单
    			if(e.button == 2) {
					// 清除上一次左键点击的样式
					trList[lastIndex_leftBtn].style.background = "none";
					checkList[lastIndex_leftBtn].checked = false;
                    container.style.overflow = "hidden";
                    this.index = i + 1;
                    if (!(fileList[i])) {
                        return;
                    }
                    else {
                        if (this.isClick) {
                            // 清除背景颜色
                            trList[this.index].style.background = "none";
                            // 不选中方框
                            checkList[this.index].checked = false;
                        }
                        else {
                            // 清除上一次点击的样式
                            trList[lastIndex_rightBtn].style.background = "none";
                            checkList[lastIndex_rightBtn].checked = false;
                            // 添加背景色
							trList[this.index].style.background = "#e8f6fd";
							// 选中方框
                            checkList[this.index].checked = true;
                            lastIndex_rightBtn = this.index; //保存当前的index
                        }
                    }
                    let menu = document.getElementsByClassName("menu")[0];  //右键的菜单
                    select_file = (fileList[i].getElementsByTagName("span")[0]).innerText;  //当前点击的文件名
					menu.style.display = 'block';
					// 根据鼠标点击位置和浏览器顶部的距离更改菜单的位置
					let h = mousePos(e);
					if(h < 700) {
						menu.style.top = h - 210 + "px";
					}
    				else{
						menu.style.top = "490px";
					}
				}
				// 左键关闭菜单
    			else if(e.button == 0) {
                    container.style.overflow = "auto";
                    menu.style.display = 'none';
    			}
    		}
    	}
    })();

    // 点击路径跳转
	let systemList = document.getElementsByClassName("file_system_a"),
        systemLen = systemList.length;
    (function () {
        for(let i = 0; i < systemLen; i++) {
            systemList[i].onclick = function () {
				let current = systemList[i].innerText,
					index_click = current_dir.indexOf(current),
					jump_num = 0;
                if(index_click !== -1) {
					jump_num = (current_dir.length - 1) - index_click;
					for(let i = 0; i < jump_num; i++) {
						returnFile();
					}
				}
                else{
					jump_num = (current_dir.length - 1) - 2;
					for(let i = 0; i < jump_num; i++) {
						returnFile();
					}
				}
            }
        }
    })();
}

/* 补0
*  @params
*      s [number] 数字
*  @return
*      s [string] 补0后的字符串
*/
function addZero(s) {
    return s < 10 ? '0' + s : s;
}

/* 添加文件至数组
*  @params
*      file [string] 文件名
*  @return
*/
function addList(file) {
    if(select_list.indexOf(file) == -1) {
        select_list.push(file);
    }
    else{
        return;
    }
}

/* 清除选中框样式
*  @params
*  @return
*/
function clearBox() {
    let trList = document.getElementsByTagName("tr"),
        checkList = document.getElementsByClassName("checkbox"),  //选择框
        checkLen = checkList.length;
    for (let i = 0; i < checkLen; i++) {
        checkList[i].checked = false;
        trList[i].style.background = "none";
    }
}

/* 隐藏更多按钮
*  @params
*  @return
*/
function clearMore() {
    let checkList = document.getElementsByClassName("checkbox"),  //文件左侧的选择框
        more_show = document.getElementsByClassName("more")[0]; //更多按钮
        checkLen = checkList.length,
        clicked_len = 0;
    
    for(let i = 0; i < checkLen; i++) {
		if(checkList[i].checked) {
			clicked_len++;
		}
	}
	if(clicked_len === 0) {
		more_show.style.display = "none"; //隐藏更多按钮
	}
	else{
		more_show.style.display = "block"; //显示更多按钮
    }
}

/* 全选文件
*  @params
*  @return
*/
function checkall() {
    let checkList = document.getElementsByClassName("checkbox"),  //选择框
        more_show = document.getElementsByClassName("more")[0]; //更多按钮
    
    more_show.style.display = checkList[0].checked ? "block" : "none";
    for (var i = 0; i < checkList.length; i++) {
        checkList[i].checked = checkList[0].checked;
    }
}

/* 将勾选的数据添加到数组中
*  @params
*  @return
*      select_list [array] 选择的文件数组
*/
function checkSelect() {
    let checkList = document.getElementsByClassName("checkbox"),  //选择框
		fileList = document.getElementsByClassName("file_name"),
		checkLen = checkList.length,
		fileLen = fileList.length;
    for (let i = 0; i < checkLen; i++) {
        if(checkList[i].checked) {
            addList(fileList[i-1].innerText);
        }
    }
    for(let i = 0; i < select_list.length; i++) {
        select_list[i] = "\"" + select_list[i] + "\"";
    }
    return select_list;
}

/* 鼠标点击位置到浏览器顶部的距离
*  @params
*      e 
*  @return
*      height [number] 高度
*/
function mousePos(e){
    e = e || window.event;
    let scrollY=document.documentElement.scrollTop||document.body.scrollTop;  //分别兼容ie和chrome
    let height = e.pageY || (e.clientY+scrollY);  //兼容火狐和其他浏览器
    return height;
}

/* 新建文件夹
*  @params
*  @return
*/
function isNew() {
    if(!newClick) {
        newFile();
    }
    else{
        let new_input = document.getElementsByClassName("new_input")[0];
        new_input.focus();
    }
}

/* 新建文件夹
*  @params
*  @return
*/
function newFile() {
    newClick = true;
    // 回到顶部
    $('html,body').animate({scrollTop: '0px'}, 800);
    let con = document.getElementsByClassName("content")[0];
    $(con).scrollTop(0);
    // 创建tbody格式
    let table = document.getElementsByTagName("table")[0];
    let tbody = document.createElement("tbody"),
        tr = document.createElement("tr"),
        td1 = document.createElement("td"),
        td2 = document.createElement("td"),
        td3 = document.createElement("td"),
        div = document.createElement("div"),
        new_input = document.createElement("input"),
        span = document.createElement("span"),
        i1 = document.createElement("i"),
        i2 = document.createElement("i");

    // 每行都加一个checkbox
    let checkbox_td = document.createElement("td"),
        label = document.createElement("label"),
        input = document.createElement("input"),
        i = document.createElement("i");

    checkbox_td.className = "tdwidthbox";
    label.className = "checklabel";
    input.className = "checkbox";
    input.type = "checkbox";
    i.className = "check";

    label.appendChild(input);
    label.appendChild(i);
    checkbox_td.appendChild(label);
    tr.appendChild(checkbox_td);

    tr.className = "trstyle";
    td1.className = "tdwidth1";
    td2.className = "tdwidth2";
    td2.innerText = "-";
    td3.className = "tdwidth3";

    // 文件创建时间
    let myDate = new Date();
    let month = addZero(myDate.getMonth() + 1),
        date = addZero(myDate.getDate()),
        hour = addZero(myDate.getHours()),
        min = addZero(myDate.getMinutes()),
        sec = addZero(myDate.getSeconds());

    // 创建时间
    td3.innerText = myDate.getFullYear() + "-" + month + "-" + date + " " + hour + ":" + min + ":" + sec;
    div.className = "filename";
    new_input.type = "text";
    new_input.className = "new_input";
    span.className = "icon";
    i1.className = "icon_1";
    i2.className = "icon_2";

    span.insertBefore(i1, div.children[1]); //添加兄弟节点
    div.insertBefore(new_input, div.children[0]);
    span.appendChild(i2); //添加子节点
    div.appendChild(span);
    td1.appendChild(div);

    // 添加文件夹标识
    let label_file = document.createElement("label"),
        i_file = document.createElement("i");
    label_file.className = "dir_label";
    label_file.appendChild(i_file);
    i_file.className = "dir_i";
    td1.appendChild(label_file);

    tr.appendChild(td1);
    tr.appendChild(td2);
    tr.appendChild(td3);
    tbody.appendChild(tr);
    // table.appendChild(tbody);
    table.insertBefore(tbody, table.children[1]);

    // 保存按钮
    i1.onclick = function (e) {
        stopPropagation(e);
        flag = true;
        let input_value = new_input.value;
        if (!input_value) {
            alert("文件名称不能为空，请重新输入！");
            new_input.focus(); //光标回到输入框内
        }
        else {
            // 验证文件名
            if (!validateFileName(input_value)) {
                alert("文件名不能包含以下字符:[\\\\/:*?\"<>|]");
                new_input.focus();  //光标定位到输入框中
            }
            else {
                let new_data = "{\"Opt\"" + ":" + "1" + "," + "\"DirName\"" + ":" + "[\"" + input_value + "\"]" + "}";
                $.ajax(
                    {
                        url: home_rpc,
                        data: new_data,
                        type: "POST",
                        async: false,
                        success: function (data)
                        {
                            if(data.code == 1000){
                                console.log(data.description);
                                let new_span = document.createElement("span"),
                                    file_txt = input_value;
                                // 隐藏新建文件夹的框，使添加的文件直接加入表格中
                                new_input.className = "hide";
                                i1.className = "hide";
                                i2.className = "hide";
                                span.className = "";
                                span.innerText = file_txt;
                                current_file = ".";
                                fileShow(current_file);
                                return true;
                            }
                            else {
                                alert(data.description);
                                return false;
                            }
                        },
                        error: function () {
                            alert("Network error!")
                        }
                    });
            }
        }
        newClick = false;
    }

    // 取消按钮
    i2.onclick = function () {
        fileShow(current_file);
        newClick = false;
    }
}

/* 文件查看刷新
*  @params
*      ret [string] 当前文件夹
*  @return 
*/
function refresh() {
    let dir = ".";
    let icon_refresh = document.getElementsByClassName("iconfont-refresh")[0],
            rotateval = 0;
    function rot() {
        rotateval = rotateval + 1;
        if(rotateval === 360) {
            clearInterval(interval);
            rotateval = 0;
            fileShow(dir);
        }
        icon_refresh.style.transform = 'rotate('+rotateval+'deg)';
    }
    let interval = setInterval(rot, 5);
}

/* 创建上传/下载列表
*  @params
*      name,size,dir [string] 文件名，文件大小，上传目录
*  @return
*/
function newLoadli(name, size, dir) {
    let uploadList = document.getElementById("uploadList");
    let li = document.createElement("li"),
        div_pro = document.createElement("div"),
        div_info = document.createElement("div"),
        div_name = document.createElement("div"),
        div_size = document.createElement("div"),
        div_path = document.createElement("div"),
        div_sta = document.createElement("div"),
        div_ope = document.createElement("div"),
        em1 = document.createElement("em"),
        em2 = document.createElement("em");

    li.className = "status";
    div_pro.className = "process";
    div_info.className = "file-info";
    div_name.className = "file-name";
    div_size.className = "file-size";
    div_path.className = "file-path";
    div_sta.className = "file-status";
    div_ope.className = "file-operate";
    em1.className = "pause";
    em2.className = "remove";

    div_name.innerText = name;
    div_size.innerText = size;
    div_path.innerText = dir;
    div_sta.innerText = "正在上传";

    div_ope.appendChild(em1);
    div_ope.appendChild(em2);
    div_info.appendChild(div_name);
    div_info.appendChild(div_size);
    div_info.appendChild(div_path);
    div_info.appendChild(div_sta);
    div_info.appendChild(div_ope);
    li.appendChild(div_pro);
    li.appendChild(div_info);
    uploadList.appendChild(li);
}

/* 分片
*  @params
        fileSize 文件大小
*  @return
        chunkSize 分片的每片大小
*/
function chunk(fileSize) {
    let chunkSize = 0;
    //文件大小不大于10M
    if(fileSize <= (10 * 1024 * 1024)) {
        chunkSize = fileSize;
        console.log("0M-10M: " + chunkSize);
    }
    //文件大小范围：10M-1G
    else if(fileSize > (10 * 1024 * 1024) && fileSize <= (1024 * 1024 * 1024)) {
        chunkSize = 1024 * 1024 * 10; //10M
        console.log("10M-1G: " + chunkSize);
    }
    //文件大小范围：1G-10G
    else if(fileSize > (1024 * 1024 * 1024) && fileSize <= (10 * 1024 * 1024 * 1024)) {
        chunkSize = Math.ceil(fileSize / 100); //分为100份
        console.log("1G-10G: " + chunkSize);
    }
    //文件大小大于10G
    else{
        alert("文件大小超过10G，请分卷压缩后上传！");
        chunkSize = 0;
    }
    return chunkSize;
}

/* 文件上传进度
*  @params
*      progress [object] 上传进度对象
*  @return
*/
function updateProgress(progress) {
    let uploadList = document.getElementById("uploadList"),
        len = uploadList.children.length;
    let process = document.getElementsByClassName("process")[len-1],
        status = document.getElementsByClassName("file-status")[len],
        operate = document.getElementsByClassName("file-operate")[len],
        em1 = operate.getElementsByTagName("em")[0],
        em2 = operate.getElementsByTagName("em")[1],
        total = document.getElementsByClassName("total")[0];
    if (progress.lengthComputable) {
        console.log("loaded:" + progress.loaded, "total:" + progress.total);
        let current_progress = progress.loaded / progress.total; //当前片的进度
        process_global = (((chunkNum_uploaded - 1) / chunkNum) + (current_progress / chunkNum)) * 100; //每个文件总进度 = （已上传的片数/总片数 + 当前片的进度/总片数） * 100
        let percent = process_global.toFixed(2) + "%";
        console.log("percent:" + percent);
        process.style.width = percent; //每个文件的进度
        status.innerText = percent; //每个文件的进度值
        total_percent[file_index] = process_global.toFixed(2);
        let len = file_arr.length;
            total_proc = 0; //总进度
        for(let i = 0; i < total_percent.length; i++) {
            let sum = total_percent[i] / len;
            total_proc += sum;
        }
        console.log("total_percent:" + Math.round(total_proc));
        total.style.width = Math.round(total_proc) + "%"; //总进度
        if (process_global == 100) {
            status.innerText = "上传成功";
            status.style.color = "#9a079a";
            em1.className = "clear";
            em2.className = "";
        }
    }
}

/* 上传文件
*  @params
*  @return
*/
function upload() {
    upload_type = 1; //上传文件
    fileObj = document.getElementById('file').files[0];
    let file_name = fileObj.name,
        file_size = bytesToSize(fileObj.size),
        dir = localStorage.getItem("select_dir");
    fileSize = fileObj.size;
    console.log(fileObj);
    console.log("file_name:" + file_name);
    console.log("file_size:" + file_size);
    console.log("dir:" + dir);
    newLoadli(file_name, file_size, dir);
    
    chunkSize = chunk(fileSize);
    chunkNum = Math.ceil(fileSize / chunkSize);
    md5_file = b64_md5(fileObj);
    let upload_data = "{\"Option\"" + ":" + "\"" + "uploadFile" + "\"" + "," + "\"FileName\"" + ":" + "\"" + file_name + "\"" + "," + "\"Size\"" + ":" + "\"" + fileSize + "\"" + "," + "\"ChunkNum\"" + ":" + "\"" + chunkNum + "\"" + "," + "\"MD5\"" + ":" + "\"" + md5_file + "\"" + "," + "\"ChunkPos\"" + ":" + "\"" + 1 + "\"" + "}";
    console.log(upload_data);

    let ret = confirm("是否将该文件上传至当前目录？");
    if (ret) {
        $.ajax(
            {
                url: uploadreq_rpc,
                data: upload_data,
                type: "POST",
                async: false,
                success: function (data)
                {
                    if(data.code == 1000){
                        console.log(data.description);
                        if(endupload_flag) {
                            uploadFile(0);
                        }
                        return true;
                    }
                    else {
                        alert(data.description);
                        return false;
                    }
                },
                error: function () {
                    alert("Network error!")
                }
            });
    }
    else {
        return;
    }
}

/* 续传文件
*  @params
*  @return
*/
function reUpload() {
    let file_name = fileObj.name;

    let upload_data = "{\"Option\"" + ":" + "\"" + "reUploadFile" + "\"" + "," + "\"FileName\"" + ":" + "\"" + file_name + "\"" + "," + "\"Size\"" + ":" + "\"" + fileSize + "\"" + "," + "\"ChunkNum\"" + ":" + "\"" + chunkNum + "\"" + "," + "\"MD5\"" + ":" + "\"" + md5_file + "\"" + "," + "\"ChunkPos\"" + ":" + "\"" + chunkNum_uploaded + "\"" + "}";
    console.log(upload_data);
    $.ajax(
        {
            url: uploadreq_rpc,
            data: upload_data,
            type: "POST",
            async: false,
            success: function (data) 
            {
				if(data.code == 1000){
					console.log(data.description);
                    uploadFile(end-chunkSize);
                    return true;
                }
                else {
                    alert(data.description);
                    return false;
                }
            },
            error: function () {
                alert("Network error!")
            }
        });
}

/* 取消上传文件
*  @params
*  @return
*/
function cancelUpload() {
    let upload_data = "{\"Option\"" + ":" + "\"" + "uploadCancel" + "\"" + "," + "\"FileName\"" + ":" + "\"" + fileObj.name + "\"" + "," + "\"Size\"" + ":" + "\"" + "" + "\"" + "," + "\"ChunkNum\"" + ":" + "\"" + "" + "\"" + "," + "\"MD5\"" + ":" + "\"" + "" + "\"" + "," + "\"ChunkPos\"" + ":" + "\"" + "" + "\"" + "}";
    console.log(upload_data);
    $.ajax(
        {
            url: uploadreq_rpc,
            data: upload_data,
            type: "POST",
            async: false,
            success: function (data) 
            {
				if(data.code == 1000){
					console.log(data.description);
                    return true;
                }
                else {
                    alert(data.description);
                    return false;
                }
            },
            error: function () {
                alert("Network error!")
            }
        });
}

/* 分片上传文件
*  @params
        start [number] 起始字节
*  @return
*/
function uploadFile(start) {
    current_file = ".";
    endupload_flag = false;
    // 上传完成 
    if (start >= fileSize) {
        console.log("上传完成......");
        endupload_flag = true;
        process_global = 0;
        chunkNum_uploaded = 1;
        if(upload_type === 1) {
            fileShow(current_file);
        }
        else{
            file_index++;
            if(file_index >= file_arr.length) {
                fileShow(current_file);
            }
            else{
                uploadEver(file_index);
            }
        }
        return;
    }
    // 获取文件块的终止字节
    end = (start + chunkSize > fileSize) ? fileSize : (start + chunkSize);

    // 将文件切块上传
    let form_data = new FormData(document.getElementById('filename')); //获取表单信息
    let formData = new FormData();
    if(!form_data.get('uploadfile').name) {
        formData.append("uploadfile", file_obj.slice(start, end)) //将获取的文件分片赋给新的对象
    }
    else{
        formData.append("uploadfile", form_data.get("uploadfile").slice(start, end)) //将获取的文件分片赋给新的对象
    }

    request = $.ajax({
        url: upload_rpc,
        data: formData,
        type: "POST",
        cache: false,
        processData: false,
        contentType: false, //必须false才会自动加上正确的Content-Type
        //这里我们先拿到jQuery产生的 XMLHttpRequest对象，为其增加 progress 事件绑定，然后再返回交给ajax使用
        xhr: function () {
            let xhr = $.ajaxSettings.xhr();
            if (xhr.upload) {
                xhr.upload.onprogress = function (progress) {
                    self.updateProgress(progress);
                };
            }
            return xhr;
        },
        success: function(data) 
        {
            if(data.code == 1000){
                console.log(data.description);
                chunkNum_uploaded ++;
                console.log("准备上传第" + chunkNum_uploaded + "片......");
                uploadFile(end);
            }
            else{
                alert(data.description);
                return false;
            }
        }
    });
}

/* 上传文件夹
*  @params
*  @return
*/
function uploadDir() {
    upload_type = 2; //上传文件夹
    $('#folder').change(function(e){
        let folder_name = null; //文件夹名
        let files = e.target.files; //所有文件
        file_arr = files;
        folder_name = (files[0].webkitRelativePath).split('/')[0];
        
        //新建上传的同名文件夹
        let new_data = "{\"Opt\"" + ":" + "1" + "," + "\"DirName\"" + ":" + "[\"" + folder_name + "\"]" + "}";
        console.log(new_data);
        $.ajax(
        {
            url: home_rpc,
            data: new_data,
            type: "POST",
            async: false,
            success: function (data)
            {
                if(data.code == 1000){
                    console.log(data.description);
                    fileShow(folder_name);
                    return true;
                }
                else {
                    alert(data.description);
                    return false;
                }
            },
            error: function () {
                alert("Network error!")
            }
        });
        uploadEver(0);
    });
}

/* 排队上传单个文件
*  @params
        index 文件数组索引
*  @return
*/
function uploadEver(index) {
    file_obj = file_arr[index];
    let file_name = file_obj.name,
        file_size = bytesToSize(file_obj.size),
        dir = localStorage.getItem("select_dir");
        fileSize = file_obj.size;
    console.log(file_obj);
    console.log("file_name:" + file_name);
    console.log("fileSize:" + fileSize);
    console.log("dir:" + dir);
    newLoadli(file_name, file_size, dir);

    chunkSize = chunk(fileSize);
    chunkNum = Math.ceil(fileSize / chunkSize);
    md5_file = b64_md5(file_obj);
    let upload_data = "{\"Option\"" + ":" + "\"" + "uploadFile" + "\"" + "," + "\"FileName\"" + ":" + "\"" + file_name + "\"" + "," + "\"Size\"" + ":" + "\"" + fileSize + "\"" + "," + "\"ChunkNum\"" + ":" + "\"" + chunkNum + "\"" + "," + "\"MD5\"" + ":" + "\"" + md5_file + "\"" + "," + "\"ChunkPos\"" + ":" + "\"" + 1 + "\"" + "}";
    console.log(upload_data);
    $.ajax(
        {
            url: uploadreq_rpc,
            data: upload_data,
            type: "POST",
            async: false,
            success: function (data)
            {
                if(data.code == 1000){
                    console.log(data.description);
                    uploadFile(0);
                    return true;
                }
                else {
                    alert(data.description);
                    return false;
                }
            },
            error: function () {
                alert("Network error!")
            }
        });
}

/* 返回上一级目录
*  @params
*  @return
*/
function returnFile() {
    current_file = "..";
    fileShow(current_file);
}

/* 删除文件
*  @params
*  @return
*/
function deleteFile() {
    current_file = ".";
    checkSelect();
    let del_data = "{\"Opt\"" + ":" + "2" + "," + "\"DirName\"" + ":" + "[" + select_list + "]" + "}";
    console.log(del_data);
    $.ajax(
        {
            url: home_rpc,
            data: del_data,
            type: "POST",
            async: false,
            success: function (data) 
            {
				if(data.code == 1000){
					console.log(data.description);
                    let menu = document.getElementsByClassName("menu")[0];  //右键的菜单
                    menu.style.display = "none";
                    select_list = [];
                    fileShow(current_file);
                    clearMore();
                    return true;
                }
                else {
                    alert(data.description);
                    return false;
                }
            },
            error: function () {
                alert("Network error!")
            }
        });
}

/* 下载文件
*  @params
*  @return
*/
function downloadFile() {
    console.log(select_file);
    let form = document.createElement("form"),
        input = document.createElement("input");
    form.style.display = "none";
    form.method = "post";
    form.action = download_rpc;
    form.enctype = "multipart/form-data";
    input.type = "hidden";
    input.name = "downloadfile";
    input.value = select_file;
    form.appendChild(input);
    document.body.appendChild(form);

    let form_data = new FormData(form);
    form.submit();
}

/* 判断上传任务列表是否为空
*  @params
*  @return
*/
function isEmptyUpload() {
    let uploadList = document.getElementById("uploadList"),
        fileInfoList = uploadList.getElementsByClassName("file-info"),
        len = fileInfoList.length,
        nothing = document.getElementsByClassName("nothing")[0],
        progress = document.getElementsByClassName("upload-progress")[0],
        upload_img = nothing.getElementsByTagName("img")[0],
        download_img = nothing.getElementsByTagName("img")[1],
        info = nothing.getElementsByClassName("info")[0];
    // 判断列表内是否有任务
    if(len == 0) {
        progress.style.display = "none";
        nothing.style.display = "block";
        upload_img.style.display = "";
        download_img.style.display = "none";
        info.innerText = "当前没有上传任务喔~";
    }
    else{
        nothing.style.display = "none";
        progress.style.display = "block";
    }
}

/* 判断下载任务列表是否为空
*  @params
*  @return
*/
function isEmptyDownload() {
    let downloadList = document.getElementById("downloadList"),
        fileInfoList = downloadList.getElementsByClassName("file-info"),
        len = fileInfoList.length,
        nothing = document.getElementsByClassName("nothing")[0],
        progress = document.getElementsByClassName("download-progress")[0],
        upload_img = nothing.getElementsByTagName("img")[0],
        download_img = nothing.getElementsByTagName("img")[1],
        info = nothing.getElementsByClassName("info")[0];
    // 判断列表内是否有任务
    if(len == 0) {
        progress.style.display = "none";
        nothing.style.display = "block";
        upload_img.style.display = "none";
        download_img.style.display = "";
        info.innerText = "当前没有下载任务喔~";
    }
    else{
        nothing.style.display = "none";
        progress.style.display = "block";
    }
}

/* 跳转到传输列表
*  @params
*  @return
*/
function toTransport() {
    current_file = ".";
    let upload_module = document.getElementsByClassName("upload-progress")[0], //上传
        uploadList = document.getElementById("uploadList"),
        download_module = document.getElementsByClassName("download-progress")[0], //下载
        downloadList = document.getElementById("downloadList"), 
        nav_title = document.getElementsByClassName("nav-title")[0],
        netdisk = nav_title.getElementsByTagName("li")[0],
        transport = nav_title.getElementsByTagName("li")[1],
        transport_content = document.getElementsByClassName("transport-content")[0],
        main_content = document.getElementsByClassName("main-content")[0],
        disk = document.getElementsByClassName("disk")[0],
        trans = document.getElementsByClassName("trans")[0],
        download = trans.getElementsByTagName("div")[0], //左侧下载菜单
        upload = trans.getElementsByTagName("div")[1]; //左侧上传菜单
    // 顶部导航的显示
    main_content.style.display = "none";
    disk.style.display = "none";
    netdisk.className = "";
    transport_content.style.display = "block";
    trans.style.display = "block";
    transport.className = "active";

    isEmptyUpload();

    // 点击下载
    download.onclick = function () {
        download.style.background = "#e2ddec";
        upload.style.background = "#f8f7f7";
        upload_module.style.display = "none"
        download_module.style.display = "block";
        isEmptyDownload();
    }

    // 点击上传
    upload.onclick = function () {
        upload.style.background = "#e2ddec";
        download.style.background = "#f8f7f7";
        download_module.style.display = "none";
        upload_module.style.display = "block";
        isEmptyUpload();
    }

    let liList = uploadList.getElementsByTagName("li"),
        total = document.getElementsByClassName("total")[0],
        operationList = document.getElementsByClassName("file-operate"),
        opeLen = operationList.length;
    (function () {
        for (let i = 1; i < opeLen-1; i++) {
            let em_btn = operationList[i].getElementsByTagName("em")[0],
                em_cancel = operationList[i].getElementsByTagName("em")[1];
            
            em_btn.onclick = function () {
                if(em_btn.className != 'clear') {
                    // 如果当前为暂停图标
                    if(em_btn.className == "pause") {
                        em_btn.className = "continue";
                        request.abort();
                        fileShow(current_file);
                    }
                    // 如果当前为继续图标
                    else{
                        em_btn.className = "pause";
                        reUpload();
                    }
                }
                else{ // 如果当前为清除图标
                    liList[i-1].style.display = "none";
                    isEmptyUpload();
                }
            }
            // 点击移除图标
            em_cancel.onclick = function () {
                request.abort();
                cancelUpload();
                uploadList.removeChild(uploadList.children[i-1]);
                total.style.width = 0;
                isEmptyUpload();
                fileShow(current_file);
            }
        }
    })();
}

/* 跳转到我的网盘
*  @params
*  @return
*/
function toDisk() {
    let nav_title = document.getElementsByClassName("nav-title")[0],
        netdisk = nav_title.getElementsByTagName("li")[0],
        transport = nav_title.getElementsByTagName("li")[1],
        transport_content = document.getElementsByClassName("transport-content")[0],
        main_content = document.getElementsByClassName("main-content")[0],
        disk = document.getElementsByClassName("disk")[0],
        trans = document.getElementsByClassName("trans")[0];
    main_content.style.display = "block";
    disk.style.display = "block";
    netdisk.className = "active";
    transport_content.style.display = "none";
    trans.style.display = "none";
    transport.className = "";
}

/* 全部暂停下载
*  @params
*  @return
*/
function pauseList() {
    let pause = document.getElementsByClassName("total-pause")[0];
    pause.innerText = pause.innerText == "全部暂停" ? "全部开始" : "全部暂停";
}

/* 全部取消下载
*  @params
*  @return
*/
function cancelList() {

}