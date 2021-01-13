/* 登录成功后进入用户主界面
*  @params
*  @return
*/
function loadPage() {
	// 显示用户名
	let user = document.getElementsByClassName("username")[0];
		user.innerHTML = username;  //用户名

	let container = document.getElementsByClassName("content")[0],  //文件目录表格所在的区域
		menu = document.getElementsByClassName("menu")[0],  //右键的菜单
		file_system = document.getElementsByClassName("file_system")[0];  //当前路径
	let index_data = "{\"Opt\"" + ":" + "0" + "," + "\"DirName\"" + ":" + "[\"" + current_file + "\"]" + "}";
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

	message = json_str; //转换为json格式
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

				// 遍历对象的某个属性
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
					json_tbody.appendChild(json_tr);
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

						if(key == "DirName") {  //文件夹的图标
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
	};

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
				current_file = key_word[i]; //存储当前点击的文件夹
				if(i_list[i].className != "dir_i") { //是文件不可进入
					return;
				}
				else{ //文件夹可以进入
					select_dir = (fileList[i].getElementsByTagName("span")[0]).innerText;  //存储当前点击的文件夹
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
					select_dir = (fileList[i].getElementsByTagName("span")[0]).innerText;  //存储当前点击的文件夹
					window.localStorage.setItem('select_dir',select_dir);
					fileShow(current_file);
				}
            }
        }
    })();

	// 屏蔽默认右键菜单
	container.oncontextmenu = function (event) {
		event.preventDefault();
	};

	// 右键文件弹出菜单
	(function() {
    	for(let i = 0; i < fileLen; i++) {
            fileList[i].index = i;  //自定义属性index保存索引
            fileList[i].isClick = false;   //定义点击开关
    		fileList[i].onmousedown = function(e) {
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

	// 整个页面点击鼠标左键关闭菜单
	let html = document.getElementsByTagName("html")[0];
	html.onclick = function(e) {
		stopPropagation(e);
		menu.style.display = 'none';
	}

	// 鼠标经过上传按钮
	let upload_btn = document.getElementsByClassName("upload")[0],
		upload_ul = document.getElementsByClassName("upload_file")[0];

	//鼠标经过
	upload_btn.onmouseenter = function () {
		upload_ul.style.display = 'block';
	}
	//鼠标离开
	upload_btn.onmouseleave = function () {
		upload_ul.style.display = 'none';
	}
}
