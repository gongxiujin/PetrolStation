<template>
  <div class="app-container">
    <el-col :span="18" :offset="6">
      <el-col :span="12">
        <el-col>
          <img class="img-circle" :src="details.avatar+'?x-oss-process=image/resize,m_lfit,h_40,w_40/format,jpg'">
          <router-link :to="'/users/'+details.user_id" class="link-type">
            <span class="username text-muted">{{details.nickname}}</span>
          </router-link>
        </el-col>
        <el-col :offset="2">
          <el-col style="padding-top: 30px">
            <el-col>{{details.content}}</el-col>
          </el-col>
          <el-col :span="8" v-for="img in details.images" style="padding-top: 30px">
            <img :src="img.thumbnail_url" @click="showImage(img.image_url)" width="100%">
          </el-col>
        </el-col>
      </el-col>
      <el-col :span="18" style="padding-top: 30px">
        <img v-for="(like, index) of details.likes" v-if="index<10" class="img-circle" @click="userDetails(like.user_id)" :src="like.avatar+'?x-oss-process=image/resize,h_30,w_30/format,jpg'" alt="" style="height:30px;width:30px">
        <p v-if="details.like>10" style="color: #777;font-size:14px">等{{details.like}}人觉得挺赞</p>
      </el-col>
      <el-col :span="18" v-for="comment in details.comments" style="padding-top: 30px">
        <el-col :span="3" style="padding-right:0; text-align: right">
          <router-link :to="'/users/'+comment.user_id" class="link-type">
            <span class="">{{comment.user_nick_name}}</span>
          </router-link>
        </el-col>
        <el-col :span="21" style="text-align: left;padding-left:5px">
            <span v-if="comment.replay_user_nick_name"> 回复
              <router-link :to="'/users/'+comment.replay_user_id" class="link-type">
                <span class="">{{comment.replay_user_nick_name}}</span>
              </router-link>
              :{{comment.content}}
            </span>
          <span v-else :span="18">
            :{{comment.content}}
            </span>
          <span :span="3">{{comment.create_time}}</span>
        </el-col>
        <el-col :span="12" style="text-align: right">
          <el-button type="text" @click="deleteComment(comment._id)">删除</el-button>
        </el-col>

      </el-col>
    </el-col>
    <el-dialog v-el-drag-dialog :center="true" :visible.sync="dialogTableVisible" title="图片详情" width="70%" @dragDialog="handleDrag">
      <div style="text-align: center">
        <img :src="dialogImage" class="image">
      </div>
    </el-dialog>
  </div>

</template>
<script>
  import {moments, delete_comments} from '@/api/article'
  import waves from '@/directive/waves' // waves directive
  import {parseTime} from '@/utils'
  import Mallki from '@/components/TextHoverEffect/Mallki'
  import PanThumb from '@/components/PanThumb'
  import elDragDialog from '@/directive/el-drag-dialog'

  // arr to obj, such as { CN : "China", US : "USA" }


  export default {
    name: 'Comments',
    components: {PanThumb},
    directives: {waves, elDragDialog},
    filters: {},
    data() {
      return {
        dialogTableVisible: false,
        dialogImage: '',
        page: 1,
        tableKey: 0,
        details: {
          avatar: '',
          nickname: "",
          images: [],
          gender: ""

        },
        listLoading: true,
        total: 0,
        query: {
          id: ""
        }
      }
    },
    created() {
      this.query.id = this.$route.params.id;
      this.getList()
    },
    methods: {
      getList() {
        moments(this.query).then(response => {
          this.details = response.data.moments;
          setTimeout(() => {
            this.listLoading = false;
          }, 0.5 * 1000)
        })
      },
      showImage(image_url){
        this.dialogImage = image_url+'?x-oss-process=image/resize,m_lfit,h_400,w_400/format,jpg';
        this.dialogTableVisible = true;
      },
      userDetails(user_id){
//          this.$router.push('/users/'+user_id);
      },
      handleDrag(){

      },
      deleteComment(id) {
        this.listLoading = true;
        delete_comments(id).then(response => {
          setTimeout(() => {
            this.$message({
              message: '操作成功',
              type: 'success'
            })

          }, 0.5 * 1000);
          this.getList();
        });

        this.listLoading = false;
      }
    }
  }
</script>
<style lang="scss" scoped>
  .app-container {

  .img-circle {
    border-radius: 50%;
    width: 40px;
    height: 40px;
    float: left;
  }

  .username {
    display: block;
    margin-left: 50px;
    padding-top: 22px;
    font-size: 16px;
    color: #000;
  }

  .text-muted {
    color: #777;
  }
  .icon-front {
    font-family: "iconfont" !important;
    font-size: 16px;
    font-style: normal;
    -webkit-font-smoothing: antialiased;
    -webkit-text-stroke-width: 0.2px;
    -moz-osx-font-smoothing: grayscale;
  }
  }
</style>
