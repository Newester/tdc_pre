<template>
  <div id="app">
    <img src="./assets/logo.png" />
    <router-view />
    <el-upload
      class="avatar-uploader"
      action
      name="file1"
      :show-file-list="false"
      :http-request="uploadImg"
    >
      <el-button id="upload-btn">上传</el-button>
    </el-upload>
    <el-button id="download-btn" @click="downloadFile">下载</el-button>
  </div>
</template>

<script>
export default {
  name: 'App',
  methods: {
    uploadImg: function (fileObj) {
      let formData = new FormData()
      formData.set('file', fileObj.file)
      this.axios
        .post('http://127.0.0.1:8080/upload', formData, {
          headers: { 'Content-type': 'multipart/form-data' }
        })
        .then()
        .catch()
    },
    downloadFile: function () {
      this.axios({
        // 用axios发送post请求
        method: 'post',
        url: 'http://127.0.0.1:8080/download', // 请求地址
        data: { filename: 'hahah.txt' },
        responseType: 'blob' // 表明返回服务器返回的数据类型
      }).then(res => {
        // 处理返回的文件流
        console.log('back.......')
        const content = res
        const blob = new Blob([content])
        const fileName = 'README.md'
        if ('download' in document.createElement('a')) {
          // 非IE下载
          const elink = document.createElement('a')
          elink.download = fileName
          elink.style.display = 'none'
          elink.href = URL.createObjectURL(blob)
          document.body.appendChild(elink)
          elink.click()
          URL.revokeObjectURL(elink.href) // 释放URL 对象
          document.body.removeChild(elink)
        } else {
          // IE10+下载
          navigator.msSaveBlob(blob, fileName)
        }
      })
    }
  }
}
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
