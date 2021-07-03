<template>
  <div class="app-container">
    <el-col :span="18" :offset="6">
      <el-form label-position="left" label-width="80px" :model="formLabelAlign">
        <el-form-item label="活动图片">
          <!--<el-input type="file" ref="file" @change="uploadFile($event.target.name, $event.target.files)"></el-input>-->
          <el-upload
            class="avatar-uploader"
            action="/api/management/load_image"
            :show-file-list="false"
            :on-success="handleAvatarSuccess"
            :before-upload="beforeAvatarUpload">
            <!--<img v-if="imageUrl" :src="imageUrl" class="avatar">-->
            <!--<i v-else class="el-icon-plus avatar-uploader-icon"></i>-->
            <img v-if="imageUrl" :src="imageUrl" class="avatar">
            <i v-else class="el-icon-plus avatar-uploader-icon"></i>
          </el-upload>
        </el-form-item>

        <el-form-item label="活动标题">
          <el-input v-model="formLabelAlign.text"></el-input>
        </el-form-item>
        <el-form-item label="指向地址">
          <el-input v-model="formLabelAlign.link"></el-input>
        </el-form-item>
        <el-button type="primary" @click="onSubmit">提交</el-button>
      </el-form>
    </el-col>
  </div>
</template>
<style>
  .avatar-uploader .el-upload {
    border: 1px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
  }
  .avatar-uploader .el-upload:hover {
    border-color: #409EFF;
  }
  .avatar-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    width: 178px;
    height: 178px;
    line-height: 178px;
    text-align: center;
  }
  .avatar {
    width: 690px;
    height: 250px;
    display: block;
  }
</style>
<script>
  import {createActivity} from '@/api/user'


  // arr to obj, such as { CN : "China", US : "USA" }


  export default {
    name: 'CreateActivity',
    data() {
      return {
        formLabelAlign: {
          text:'',
          image_url: null,
          link: ''
        },
        imageUrl: ''
      }
    },
    methods: {
      onSubmit() {
        createActivity(this.formLabelAlign).then(response => {
          setTimeout(() => {
            this.$message({
              message: '操作成功',
              type: 'success'
            })
          }, 0.5 * 1000);
          this.$router.push("/activity/index")
        });
      },
      handleAvatarSuccess(res, file) {
        this.imageUrl = res.data.image;
        this.formLabelAlign.image_url = res.data.image;
      },
      beforeAvatarUpload(file) {
        const isJPG = file.type === 'image/jpeg';
        const isLt2M = file.size / 1024 / 1024 < 2;

        if (!isJPG) {
          this.$message.error('上传头像图片只能是 JPG 格式!');
        }
//        if (!isLt2M) {
//          this.$message.error('上传头像图片大小不能超过 2MB!');
//        }
        return isJPG && isLt2M;
      }
    }
  }
</script>

