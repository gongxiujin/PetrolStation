<template>
  <div v-if="user.playerOptions">
    <div :span="13" :xs="24" class="text-center">code：{{ user.video_code }}</div>
    <video-player  class="video-player-box"
                   ref="videoPlayer"
                   :options="user.playerOptions"
                   :playsinline="true"
                   customEventName="customstatechangedeventname"
                   >
    </video-player>
    </div>
  <div v-else>
    <p>用户未上传视频</p>
    <p>邀请人：{{user.be_guided}}</p>
  </div>

</template>

<script>
  import 'video.js/dist/video-js.css'

  import { videoPlayer } from 'vue-video-player'


  export default {
  props: {
    user: {
      type: Object,
      default: () => {
        return {
          name: '',
          email: ''
        }
      }
    }
  },
  components: {
    videoPlayer
  },
  methods: {
    submit() {
      this.$message({
        message: 'User information has been updated successfully',
        type: 'success',
        duration: 5 * 1000
      })
    }
  }
}
</script>
