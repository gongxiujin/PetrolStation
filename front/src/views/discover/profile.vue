<template>
  <div class="app-container">
    <div v-if="user">
      <el-row :gutter="20">

        <el-col :span="6" :xs="24">
          <user-card :user="user"/>
        </el-col>

        <el-col :span="18" :xs="24">
          <el-card>
            <el-tabs v-model="activeTab">
              <el-tab-pane label="个人图片" name="activity">
                <activity :user="user"/>
              </el-tab-pane>
              <el-tab-pane label="动态" name="timeline">
                <timeline :user="user"/>
              </el-tab-pane>
              <el-tab-pane label="认证视频" name="account">
                <account :user="user"/>
              </el-tab-pane>
            </el-tabs>
          </el-card>
        </el-col>

      </el-row>
    </div>
  </div>
</template>

<script>
  import {mapGetters} from 'vuex'
  import UserCard from './components/UserCard'
  import Activity from './components/Activity'
  import Timeline from './components/Timeline'
  import Account from './components/Account'
  import {user_profile} from '@/api/article'

  export default {
    name: 'Profile',
    components: {UserCard, Activity, Timeline, Account},
    data() {
      return {
        user: {},
        playerOptions: {

        },
        activeTab: 'activity'
      }
    },
    computed: {
      ...mapGetters([
        'nick_name',
        'avatar',
        'images',
      ])
    },
    created() {
      if (!this.$route.query) {
        this.$message({
          message: '参数错误',
          type: 'error'
        })
      }
      this.user_id = this.$route.params.id;
      this.get_details();
    },
    methods: {
      getUser() {
        this.user = {
          name: this.name,
          role: this.roles.join(' | '),
          email: 'admin@test.com',
          avatar: this.avatar
        }
      },
      get_details() {
        user_profile(this.user_id).then(response => {
          this.user = response.data;
          if (this.user.video_path){
            this.user.playerOptions={
              muted: true,
              height: '360',
              width: '500',
              language: 'en',
              playbackRates: [0.7, 1.0, 1.5, 2.0],
              sources: [{
                type: "video/mp4",
                src: this.user.video_path
              }],
            }
          }else{
//            this.user.playerOptions={
//              muted: true,
//              height: '360',
//              width: '500',
//              language: 'en',
//              playbackRates: [0.7, 1.0, 1.5, 2.0],
//              sources: [{
//                type: "video/mp4",
//                src: this.user.video_path
//              }],
//            }
          }
          setTimeout(() => {
            this.listLoading = false;
          }, 0.5 * 1000)
        })
      },
    }
  }
</script>
