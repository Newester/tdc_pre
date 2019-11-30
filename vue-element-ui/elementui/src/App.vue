<template>
  <div id="app">
    <img src="./assets/logo.png">
    <router-view/>
    <el-upload
      class="avatar-uploader"
      action
      name="file1"
      :show-file-list="false"
      :http-request="uploadImg">
      <el-button id="upload-btn">上传</el-button>
    </el-upload>
  </div>
</template>

<script>
export default {
  name: 'App',
  methods: {
    uploadImg: function (fileObj) {
      let formData = new FormData()
      formData.set('file', fileObj.file)
      this.axios.post('http://127.0.0.1:8080/upload', formData, {
        headers: {'Content-type': 'multipart/form-data'}
      }).then().catch()
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
