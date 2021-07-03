<template>
  <div class="user-activity">
    <div class="post">
      <div class="user-block">
        <img class="img-circle" :src="user.avatar">
        <span class="username">{{user.nick_name}}</span>
      </div>
      <div class="user-images">
        <el-carousel :interval="6000" type="card" height="200px">
          <el-carousel-item v-for="img in user.images">
            <img :src="img.image" class="image">
          </el-carousel-item>
        </el-carousel>
      </div>
    </div>
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
    },
  }
</script>

<style lang="scss" scoped>
  .user-activity {
  .user-block {

  .username,
  .description {
    display: block;
    margin-left: 50px;
    padding: 2px 0;
  }

  .username{
    font-size: 16px;
    color: #000;
  }

  :after {
    clear: both;
  }

  .img-circle {
    border-radius: 50%;
    width: 40px;
    height: 40px;
    float: left;
  }

  span {
    font-weight: 500;
    font-size: 12px;
  }
  }

  .post {
    font-size: 14px;
    border-bottom: 1px solid #d2d6de;
    margin-bottom: 15px;
    padding-bottom: 15px;
    color: #666;

  .image {
    width: 200px;
    height: 200px;

  }

  .user-images {
    padding-top: 20px;
  }
  }

  .list-inline {
    padding-left: 0;
    margin-left: -5px;
    list-style: none;

  li {
    display: inline-block;
    padding-right: 5px;
    padding-left: 5px;
    font-size: 13px;
  }

  .link-black {

  &:hover,
  &:focus {
     color: #999;
   }
  }
  }

  }

  .box-center {
    margin: 0 auto;
    display: table;
  }

  .text-muted {
    color: #777;
  }
</style>
