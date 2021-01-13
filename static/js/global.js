var host = "localhost:9090";

// 链接地址
var login_href = "http://" + host + "/login",
    regist_href = "http://" + host + "/regist",
    index_href = "http://" + host + "/index",
    share_href = "http://" + host + "/share",
    logout_href = "http://" + host + "/logout";

// 接口
var login_rpc = "http://" + host + "/login_auth",
    regist_rpc = "http://" + host + "/regist",
    home_rpc = "http://" + host + "/home",
    logout_rpc = "http://" + host + "/logout",
    upload_rpc = "http://" + host + "/upload",
    download_rpc = "http://" + host + "/download",
    uploadreq_rpc = "http://" + host + "/upload_request",
    upload_rpc = "http://" + host + "/upload";

var current_file = ".";  //当前所在的文件夹
    select_file = "",  //选中的文件
    select_dir = "",  //当前点击的目录
    current_dir = [],  //当前路径数组
    select_list = [],  //选中的文件数组，用于删除
    request = null,
    newClick = false, //新建文件夹调用标识
    fileObj = null, //上传文件的文件
    upload_type = null,
    process_global = 0, //总进度
    md5_file = null,
    chunkNum = 0, //分片数
    chunkNum_uploaded = 1, //已上传片数
    end = 0, //结束字节
    file_arr = null, //文件夹的文件数组
    file_index = 0, //文件夹里的文件数组的索引
    file_obj = null, //上传文件夹的每个文件
    endupload_flag = true, //上传结束标识
    total_percent = [],
    username = localStorage.getItem("user"); //用户名

var chunkSize = 0, //每片的大小
    fileSize = 0;  //文件大小