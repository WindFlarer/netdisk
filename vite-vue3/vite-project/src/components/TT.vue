<template>
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
                    <el-button size="small" type="primary"
                        @click.stop="handleDownload(scope.$index, scope.row)">Download</el-button>

                </template>
            </el-table-column>
        </el-table>
        <el-pagination :current-page="page" :page-size="pageSize" :total="nowFile.length"
            @current-change="handlePageChange" />
    </div>
</template>
  
<script lang="ts" setup>
import { useRoute, useRouter } from 'vue-router'
import { ref, computed, watch, onBeforeMount } from "vue";
import { ElMessage, ElMessageBox } from 'element-plus';
import axios from 'axios';
const route = useRoute()
const router = useRouter()

// 文件对象
interface UserFile {
    fileName: string
    fileSize: string
    path: string
    isDir: boolean
}
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
// const nowFile = ref<UserFile[]>([]);

const token = localStorage.getItem('token')
const config = {
    headers: {
        'Authorization': token,
    }
}
const downPath = ref('D:/netdisk-download')

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
        const addedPathsFile = new Set<string>()
        const addedPathsFolder = new Set<string>()
        for (let i = 0; i < tableData1.length; i++) {
            if (tableData1[i].path.startsWith(nowPath) && tableData1[i].path.length > nowPath.length && tableData1[i].path[nowPath.length] == '/') {
                if (!tableData1[i].isDir) {
                    if (!tableData1[i].path.substring(nowPath.length + 1).includes('/')) {
                        if (!addedPathsFile.has(tableData1[i].path)) {
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
        return nowFile1
    }
    var res: UserFile[] = []
    return select1(res, tableData.value)
}
);


// 在路由参数变化时调用 fetch 函数更新数据
watch(() => route.query.path, (newPath) => {
    console.log(tableData.value)
    nowPath = newPath as string
    tableData.value = [...tableData.value]
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
                window.alert('删除成功');
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
                window.alert('删除成功');
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


</script>
<style>
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
    --el-input-bg-color: #242424;
    --el-input-border-color: #404040;
}
</style>