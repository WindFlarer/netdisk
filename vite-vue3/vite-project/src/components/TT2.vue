<template>
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
                        <el-menu-item index="2">
                            <el-icon>
                                <Menu />
                            </el-icon>
                            <span>新建文件夹</span>
                        </el-menu-item>
                        <el-menu-item index="3">
                            <el-icon>
                                <setting />
                            </el-icon>
                            <span>设置</span>
                        </el-menu-item>
                    </el-menu>
                </div>

            </el-col>
        </el-row>
    </div>
</template>
  
<script setup lang="ts">
import { ref } from 'vue';
import { ElUpload, ElProgress, ElMessage } from 'element-plus';
import axios from 'axios';
const showProgress = ref(false);
const uploadProgress = ref(0);
const uploadFileName = ref('');
const file = ref<File | null>(null)

const token = localStorage.getItem('token')
const config = {
    headers: {
        'Authorization': token,
    }
}


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
</style>