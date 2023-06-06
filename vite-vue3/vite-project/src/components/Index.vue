<template>
    <div class="cloud-drive">

        <!-- 头部 -->
        <div class="header">
            <!-- 左侧 Logo -->
            <div class="logo">
                <!-- <img src="../assets/logo.png" alt="logo" /> -->
                <el-icon :size="30">
                    <cloudy />
                </el-icon>
                <span>农林云</span>
            </div>

            <!-- 右侧用户信息 -->
            <el-dropdown class="profile">
                <div class="profile_detail">
                    <!-- <el-avatar class="el-avatar" :src="avatarUrl" :size="40"></el-avatar> -->
                    <span>{{ nickname }}</span>


                </div>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item>设置</el-dropdown-item>
                        <el-dropdown-item>DownPath:{{ downPath }}</el-dropdown-item>
                        <el-dropdown-item>Action 3</el-dropdown-item>
                        <el-dropdown-item disabled>Action 4</el-dropdown-item>
                        <el-dropdown-item divided>Action 5</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </div>

        <div class="left">
            <!-- 左侧内容 -->
            <el-row class="side-menu-row">
                <el-col class="side-menu-col">
                    <div class="menu-wrapper">
                        <el-menu active-text-color="#ffd04b" background-color="#242424" class="el-menu-vertical-demo"
                            text-color="#fff">
                            <el-menu-item index="1">
                                <el-upload ref="upload" class="upload-demo" :action="''" :show-file-list="false"
                                    :before-upload="handleBeforeUpload" :on-progress="handleUploadProgress"
                                    :on-success="handleUploadSuccess">
                                    <el-icon>
                                        <document />
                                    </el-icon>
                                    <span>上传文件</span>
                                </el-upload>
                                <el-progress :percentage="uploadProgress" v-show="showProgress"></el-progress>
                            </el-menu-item>

                            <el-menu-item index="2" @click="createDir">
                                <el-icon>
                                    <Menu />
                                </el-icon>
                                <span>新建文件夹</span>
                            </el-menu-item>

                            <el-menu-item index="3" @click="setUp">
                                <el-icon>
                                    <setting />
                                </el-icon>
                                <span>设置下载地址</span>
                            </el-menu-item>


                        </el-menu>
                    </div>

                </el-col>
            </el-row>
        </div>
        <!-- 右侧内容 -->
        <div class="right">
            <!-- 路径 -->
            <div class="path">
                <div class="pathLeft">
                    <el-icon class="pathIcon" :size="13">
                        <DArrowRight />
                    </el-icon>
                    <span class="pathWords">path:{{ nowPathDisplay }}</span>
                </div>

                <div class="backIcon" @click="pathBack">
                    <el-icon>
                        <ArrowLeftBold />
                    </el-icon>
                </div>
            </div>
            <!-- 表单 -->
            <div class="table-wrapper">
                <el-table class="table" :data="pagedData" style="width: 100%" @row-click="handleRowClick">
                    <!-- 文件名称 -->
                    <el-table-column label="FileName" width="180">
                        <template #default="scope">
                            <div style="display: flex; align-items: center">
                                <template v-if="scope.row.isDir">
                                    <el-icon>
                                        <Folder />
                                    </el-icon>
                                </template>
                                <template v-else>
                                    <el-icon>
                                        <Document />
                                    </el-icon>
                                </template>
                                <span style="margin-left: 10px">{{ scope.row.fileName }}</span>
                            </div>
                        </template>
                    </el-table-column>

                    <!-- 大小 -->
                    <el-table-column label="Size" width="180" class="size-table">
                        <template #default="scope">
                            <div style="display: flex; align-items: center">
                                <span>{{ scope.row.fileSize }}</span>
                            </div>
                        </template>
                    </el-table-column>


                    <!-- 操作按钮 -->
                    <el-table-column label="Operations">
                        <template #default="scope">
                            <el-button size="small" @click.stop="handleRename(scope.$index, scope.row)">Rename</el-button>
                            <el-button size="small" type="danger"
                                @click.stop="handleDelete(scope.$index, scope.row)">Delete</el-button>
                            <el-button :disabled="scope.row.isDir" size="small" type="primary"
                                @click.stop="handleDownload(scope.$index, scope.row)">Download</el-button>

                        </template>
                    </el-table-column>
                </el-table>
                <el-pagination :current-page="page" :page-size="pageSize" :total="nowFile.length"
                    @current-change="handlePageChange" />
            </div>
        </div>



    </div>
</template>
  
<script setup lang="ts">
import { ElAvatar, ElButton, ElDialog, ElDropdown, ElDropdownItem, ElDropdownMenu, ElForm, ElFormItem, ElInput, ElTable, ElTableColumn, ElUpload, ElProgress, ElMessageBox, ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router'
import { ref, computed, watch, onBeforeMount } from "vue";
import axios from 'axios';
const route = useRoute()
const router = useRouter()


// 初始化
// 文件对象
interface UserFile {
    fileName: string
    fileSize: string
    path: string
    isDir: boolean
}

// 上传文件
const showProgress = ref(false);
const uploadProgress = ref(0);
const uploadFileName = ref('');
const file = ref<File | null>(null)


// 分页初始化
const page = ref(1)
const pageSize = 10
const pagedData = computed(() => {
    const start = (page.value - 1) * pageSize
    const end = start + pageSize
    return nowFile.value.slice(start, end)
})

// 分页handler
const handlePageChange = (newPage: number) => {
    page.value = newPage
}
// 初始化nowPath
var nowPath = route.query.path as string
if (nowPath == null || nowPath == '/') {
    nowPath = ''
}
//初始化 所有路径 以及 当前目录文件和文件夹
const tableData = ref<UserFile[]>([]);

const token = localStorage.getItem('token')
const config = {
    headers: {
        'Authorization': token,
    }
}
const downPath = ref('D:/netdisk-download')

// 头像
const nickname = ref<string>('windflare');

// 当前路径显示
const nowPathDisplay = ref('')
nowPathDisplay.value = nowPath
if (nowPath == '') {
    nowPathDisplay.value = '/'
}



//去后端拿数据
const fetch = async () => {

    try {
        const response = await axios.get('http://localhost:8080/file/fileList', config);
        if (response.data.code == 0) {
            // 拿到数据, 存入tabledata
            tableData.value = response.data.data;
        } else {
            window.alert('获取失败');
        }

    } catch (error) {
        console.error(error);
    }
}
fetch();




//根据当前路径筛选文件
// nowFile.value = tableData.value.filter((data) => {
//     // 如果数据的 path 以 nowPath 开头，长度大于 nowPath（说明后续），路径中下一级不包含 '/' 字符（即在当前目录下）
//     return data.path.startsWith(nowPath) && data.path.length > nowPath.length && data.path[nowPath.length] == '/' && !data.path.substring(nowPath.length + 1).includes('/')

// });

//根据当前路径筛选文件与文件夹
const nowFile = computed(() => {
    // 对tabledata进行筛选
    // 所有nowPath下的文件都能直接筛选出, 但是nowPath下的文件夹需要根据给定的数据进行添加
    const select1 = (nowFile1: UserFile[], tableData1: UserFile[]) => {
        // 防止重复添加文件夹
        const addedPathsFile = new Set<string>()
        const addedPathsFolder = new Set<string>()

        for (let i = 0; i < tableData1.length; i++) {
            // 如果以nowpath开头的文件夹,并且长度大于nowpath, 并且下一个字符是'/
            if (tableData1[i].path.startsWith(nowPath) && tableData1[i].path.length > nowPath.length && tableData1[i].path[nowPath.length] == '/') {
                // 如果不是文件夹
                if (!tableData1[i].isDir) {
                    // 后续没有'/'就判定为文件
                    if (!tableData1[i].path.substring(nowPath.length + 1).includes('/')) {
                        if (!addedPathsFile.has(tableData1[i].path)) {
                            tableData1[i].fileSize = formatFileSize(tableData1[i].fileSize)
                            nowFile1.push(tableData1[i])
                            addedPathsFile.add(tableData1[i].path)
                        }
                    }
                } else {
                    if (tableData1[i].path[nowPath.length] == '/') {
                        const temp: UserFile = {
                            fileName: '',
                            fileSize: '',
                            path: '',
                            isDir: true,
                        }
                        let index = tableData1[i].path.indexOf('/', nowPath.length + 1)
                        if (index == -1) {
                            index = tableData1[i].path.length
                        }
                        temp.path = tableData1[i].path.substring(0, index)
                        temp.fileName = tableData1[i].path.substring(nowPath.length + 1, index)

                        // 如果
                        if (!(addedPathsFolder.has(temp.path))) {
                            nowFile1.push(temp)
                            addedPathsFolder.add(temp.path)
                        }
                    }
                }
            }
        }
        nowFile1.sort((a, b) => {
            if (a.isDir && !b.isDir) {
                return -1; // a should come first
            } else if (!a.isDir && b.isDir) {
                return 1; // b should come first
            } else {
                return 0; // no change in order
            }
        });
        return nowFile1
    }
    var res: UserFile[] = []
    return select1(res, tableData.value)
}
);

// 计算容量
function formatFileSize(sizeStr: string) {
    var size = parseFloat(sizeStr);
    if (size < 1024) {
        return size + 'B';
    } else if (size < 1048576) {
        return (size / 1024).toFixed(2) + 'KB';
    } else if (size < 1073741824) {
        return (size / 1048576).toFixed(2) + 'MB';
    } else if (size < 1099511627776) {
        return (size / 1073741824).toFixed(2) + 'GB';
    } else {
        return (size / 1099511627776).toFixed(2) + 'TB';
    }
}


// 在路由参数变化时调用 fetch 函数更新数据
watch(() => route.query.path, (newPath) => {
    //更新目录
    nowPath = newPath as string
    tableData.value = [...tableData.value]
    //修改path展示
    nowPathDisplay.value = nowPath
    if (nowPath == '') {
        nowPathDisplay.value = '/'
    }


    if (newPath == '/' || newPath == null) {
        location.reload()
    }
})



const handleRowClick = (row: UserFile) => {
    if (row.isDir) {
        let jumpTo = route.path + '?path=' + row.path
        router.push(jumpTo)
    }

}

const newFileName = ref('')

// 重命名
const handleRename = (index: number, row: UserFile) => {
    let promptPromise = ElMessageBox.prompt('Please input new name', 'Tip', {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
    });
    promptPromise.then(({ value }) => {
        newFileName.value = value;
        ElMessage({
            type: 'success',
            message: `success`,
        });
        if (row.isDir) {
            dirRenameFunc();
        } else {
            fileRenameFunc();
        }

    }).catch(() => {
        ElMessage({
            type: 'info',
            message: 'Input canceled',
        });
    });

    const dirRenameFunc = async () => {
        try {
            const response = await axios.put('http://localhost:8080/file/modifyDir', {
                NewDirName: newFileName.value,
                OldDirName: row.fileName,
                Path: row.path,
            }, config);

            if (response.data.code === 0) {
                window.alert('修改成功');
                fetch();
            } else if (response.data.code === 1) {
                window.alert('文件夹已存在');
            } else {
                window.alert('修改失败');
            }
        } catch (error) {
            console.error(error);
            window.alert('网络错误，请稍后重试');
        }
    };
    const fileRenameFunc = async () => {
        try {
            const response = await axios.put('http://localhost:8080/file/modifyFile', {
                NewFileName: newFileName.value,
                Path: row.path,
            }, config);

            if (response.data.code === 0) {
                window.alert('修改成功');
                fetch();
            } else if (response.data.code === 1) {
                window.alert('文件已存在');
            } else {
                window.alert('修改失败');
            }
        } catch (error) {
            console.error(error);
            window.alert('网络错误，请稍后重试');
        }
    };
};

// 删除
const handleDelete = (index: number, row: UserFile) => {
    ElMessageBox.confirm(
        'proxy will permanently delete the file. Continue?',
        'Warning',
        {
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            type: 'warning',
        }
    )
        .then(() => {
            ElMessage({
                type: 'success',
                message: 'Delete completed',
            })

            if (row.isDir) {
                dirDeleteFunc()
            } else {
                fileDeleteFunc()
            }
        })
        .catch(() => {
            ElMessage({
                type: 'info',
                message: 'Delete canceled',
            });
        });
    const fileDeleteFunc = async () => {
        try {
            const response = await axios.delete('http://localhost:8080/file/deleteFile', {
                data: { Path: row.path, }, headers: { Authorization: token }
            });

            if (response.data.code === 0) {
                fetch();
            } else if (response.data.code === 1) {
                window.alert('删除失败');
            } else {
                window.alert('修改失败');
            }
        } catch (error) {
            console.error(error);
            window.alert('网络错误，请稍后重试');
        }
    };
    const dirDeleteFunc = async () => {
        try {
            const response = await axios.delete('http://localhost:8080/file/deleteDir', {
                data: {
                    Path: row.path,
                    FileName: row.fileName,
                }, headers: { Authorization: token }
            });

            if (response.data.code === 0) {
                fetch();
            } else if (response.data.code === 1) {
                window.alert('删除失败');
            } else {
                window.alert('修改失败');
            }
        } catch (error) {
            console.error(error);
            window.alert('网络错误，请稍后重试');
        }
    };
}

//下载
const handleDownload = (index: number, row: UserFile) => {
    ElMessageBox.confirm(
        'proxy will download the file. Continue?',
        'Info',
        {
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            type: 'info',
        }
    )
        .then(() => {
            ElMessage({
                type: 'success',
                message: 'Delete completed',
            })

            fileDownloadFunc()
        })
        .catch(() => {
            ElMessage({
                type: 'info',
                message: 'Delete canceled',
            });
        });
    const fileDownloadFunc = async () => {
        try {
            const response = await axios.post('http://localhost:8080/file/download', {
                path: row.path,
                fileName: row.fileName,
                downPath: downPath.value,
            }, config)


            if (response.data.code === 0) {
                window.alert('下载成功');
                fetch();
            } else if (response.data.code === 1) {
                window.alert('下载失败');
            } else {
                window.alert('下载失败');
            }
        } catch (error) {
            console.error(error);
            window.alert('网络错误，请稍后重试');
        }
    };
}


// 上传文件
const uploadFile = () => {

}


const newDirPath = ref('')
const newDirName = ref('')
// 新建文件夹
const createDir = () => {
    let promptPromise = ElMessageBox.prompt('Please input dir name', 'Tip', {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
    });
    promptPromise.then(({ value }) => {
        if (nowPath == '/' || nowPath == null) {
            newDirPath.value = '/' + value;
        } else {
            newDirPath.value = nowPath + '/' + value;
        }

        newDirName.value = value;
        ElMessage({
            type: 'success',
            message: `success`,
        });
        createDirFunc();

    }).catch(() => {
        ElMessage({
            type: 'info',
            message: 'Input canceled',
        });
    });

    const createDirFunc = async () => {
        try {
            const response = await axios.post('http://localhost:8080/file/createDir', {
                Path: newDirPath.value,
                FileName: newDirPath.value,
            }, config);

            if (response.data.code === 0) {
                fetch();
            } else if (response.data.code === 1) {
                window.alert('文件夹已存在');
            } else {
                window.alert('修改失败');
            }
        } catch (error) {
            console.error(error);
            window.alert('网络错误，请稍后重试');
        }
    };
}

// 路径回退
const pathBack = () => {
    if (nowPath != '/' && nowPath != '') {
        nowPath = nowPath.substring(0, nowPath.lastIndexOf('/'))
        let next = route.path + '?path=' + nowPath
        router.push(next)
    }
}

const setUp = () => {
    // 点击弹出提示框输入下载地址
    let promptPromise = ElMessageBox.prompt('Please input new address', 'Config', {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
    });
    promptPromise.then(({ value }) => {
        newFileName.value = value;
        ElMessage({
            type: 'success',
            message: `success`,
        });
        downPath.value = value;

    }).catch(() => {
        ElMessage({
            type: 'info',
            message: 'Input canceled',
        });
    });

}

// 文件上传
const handleBeforeUpload = (uploadedFile: File) => {
    // 保存上传的文件并获取文件名
    file.value = uploadedFile;
    //获取文件名
    // const dotIndex = uploadedFile.name.lastIndexOf('.');
    // if (dotIndex >= 0) {
    //     fileName.value = uploadedFile.name.slice(0, dotIndex);
    // } else {
    //     fileName.value = uploadedFile.name;
    // }
    uploadFileName.value = uploadedFile.name;
    showProgress.value = true; // 显示上传进度条

    console.log(uploadFileName.value)
    // 连接后端
    const handleSubmit = async () => {
        const formData = new FormData()
        formData.append('path', '/' + uploadFileName.value)
        formData.append('fileName', uploadFileName.value)
        if (file.value) {
            formData.append('file', file.value)
        }
        try {
            const response = await axios.post('http://localhost:8080/file/upload', formData, config)
            if (response.data.code === 0) {
                fetch();
                ElMessage({
                    type: 'success',
                    message: 'upload completed',
                })
            } else if (response.data.code === 1) {
                ElMessage({
                    type: 'error',
                    message: 'already Existed',
                })
            } else if (response.data.code === 2) {
                ElMessage({
                    type: 'error',
                    message: 'cos error',
                })
            } else if (response.data.code === 3) {
                ElMessage({
                    type: 'error',
                    message: 'database error',
                })
            }
            console.log(response.data)
        } catch (error) {
            console.error(error)
        }
    };
    handleSubmit()
    return true; // 允许上传
}

const handleUploadProgress = (event: ProgressEvent) => {
    // 更新上传进度
    if (event.lengthComputable) {
        uploadProgress.value = Math.round((event.loaded / event.total) * 100);
    }
};

const handleUploadSuccess = (response: any) => {
    // 处理上传成功后的响应
    console.log('上传成功:', response);
    ElMessage.success('文件上传成功');
    showProgress.value = false; // 隐藏上传进度条
    // 在这里发送文件和文件名到后端
};


</script>
  
<style>
@font-face {
    font-family: 'CTWeiBeiSJ';
    src:
        url('../assets/CTWeiBeiSJ.woff') format('woff');
    font-weight: normal;
    font-style: normal;
}

@font-face {
    font-family: 'Facon';
    src: url('../assets/Facon.woff') format('woff');
    font-weight: normal;
    font-style: normal;
}

@font-face {
    font-family: 'ReturnPolicyDEMO-Regular';
    src: url('../assets/ReturnPolicyDEMO-Regular.woff') format('woff');
    font-weight: normal;
    font-style: normal;
}

.cloud-drive {
    padding: 20px;
    font-size: 16px;
}

/* 头部样式 */
.header {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 50px;
    background-color: #333333;

    display: flex;
    justify-content: space-between;
    align-items: center;
    /* 底部阴影 */
    box-shadow: 0 2px 4px rgba(177, 174, 174, 0.2);
}

/* 左侧logo */
.logo {
    margin-right: auto;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 350px;
    height: 40px;
    left: 100px;
}

.logo span {
    font-family: 'CTWeiBeiSJ', cursive;
    font-size: 25px;
    margin-left: 25px;
    user-select: none;
}

.logo img {
    width: 40px;
    height: 40px;
    margin-left: 50px;
}

/* 右侧用户信息 */

.path {
    width: 800px;
    height: 100px;
    display: flex;
    justify-content: space-between;
    align-items: center;
}


.pathWords {
    margin-left: 5px;
    width: 100px;

}

.backIcon {
    cursor: pointer;
}


.profile {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    width: 150px;
    height: 50px;
    cursor: pointer;
}

.profile_detail {
    display: flex;
    justify-content: center;
    align-items: center;
}

.profile .profile_detail span {
    letter-spacing: 1px;
    font-size: 17px;
    margin-right: 10px;
    color: azure;
    font-weight: normal;
}

/* 左侧菜单 */
.el-menu-item.is-active {
    color: #ffffff;
}

.side-menu-row {
    position: absolute;
    left: 100px;
    top: 55px;
    height: calc(100% - 55px);
}

.side-menu-col {
    background-color: #242424;
    width: 200px;
    height: 100%;
}

.menu-wrapper {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    width: 100%;
}

.el-menu-vertical-demo {
    width: 100%;
}


.actions {
    display: flex;
    justify-content: space-between;
    margin-bottom: 20px;
}

.upload {
    flex: 1;
    margin-right: 20px;
}

.file-list {
    max-width: 800px;
    margin: 0 auto;
}

/* 表格样式 */
.el-table {
    --el-table-tr-bg-color: #242424;
    --el-table-header-bg-color: #242424;
    --el-table-row-hover-bg-color: #aaaaaa;
    cursor: pointer;
}

.table-wrapper {
    width: 800px;
}

.el-table thead {
    cursor: default;
}

/* 分页样式 */
.el-pagination {
    --el-pagination-bg-color: #242424;
    --el-pagination-button-color: #576a8f;
    --el-pagination-button-disabled-bg-color: #333131;
}

.el-input {
    --el-input-bg-color: #ffffff;
    --el-input-border-color: #404040;
}
</style>