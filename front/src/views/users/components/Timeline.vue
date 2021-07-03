<template>
  <div class="block">
    <el-timeline>
      <el-timeline-item v-for="(item, index) of user.moments" :key="index" :timestamp="item.create_time|moment" placement="top">
        <el-card >
          <h4>{{ item.content }}</h4>
          <div v-for="img of item.images">
            <img :src="img.thumbnail_url">
          </div>
        </el-card>
      </el-timeline-item>
    </el-timeline>
  </div>
</template>

<script>
export default {
  props: {
    user: {
      type: Object,
      default: () => {
        return {
          nick_name: '',
          images: [],
          avatar: '',
        }
      }
    }
  }
}
</script>
