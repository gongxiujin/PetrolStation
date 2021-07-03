<template>
  <div class="app-container">

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column label="ID" align="center" width="130">
        <template slot-scope="scope">
            <span>{{ scope.row.ID }}</span>
        </template>
      </el-table-column>
      <el-table-column label="微信号" width="80px" align="center">
        <template slot-scope="scope">
            <span>{{ scope.row.wechat }}</span>
        </template>
      </el-table-column>
      <el-table-column label="昵称" align="center" width="130">
        <template slot-scope="scope">
          <router-link :to="'/users/'+scope.row._id" class="link-type">
          <span>{{ scope.row.nick_name }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="性别" align="center" width="50">
        <template slot-scope="scope">
          <span>{{ scope.row.gender|genderFilter }}</span>
        </template>
      </el-table-column>
      <el-table-column label="申请时间" width="180px" align="center">
        <template slot-scope="scope">
          <span>{{scope.row.create_time}}</span>
        </template>
      </el-table-column>
      <el-table-column label="图片" width="80px" align="center">
        <template slot-scope="scope">
          <el-button type="text" @click="showImage(scope.row.avatar)">查看头像</el-button>
        </template>
      </el-table-column>
      <el-table-column label="视频" width="80px" align="center">
        <template slot-scope="scope">
          <el-button v-if="scope.row.video_url" type="text" @click="showVideo(scope.row.video_url)">查看视频</el-button>
        </template>
      </el-table-column>
      <el-table-column label="推荐人" width="80px" align="center">
        <template slot-scope="scope">
          <router-link :to="'/users/'+scope.row.guided_id" class="link-type">
          <span>{{ scope.row.guided_nickname }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="180" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button type="text" @click="showTooltip(row._id, 'success')">
            通过
          </el-button>
          <el-button type="text" @click="showTooltip(row._id, 'refuse')">
          拒绝
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0"  :total="total" :page.sync="listQuery.page" @pagination="getList" />

    <el-dialog v-el-drag-dialog :visible.sync="dialogTableVisible" :title="dialogTitle" @dragDialog="handleDrag">
      <img v-if="showWindows.type==='image'" :src="showWindows.url" alt="" style="width: 400px;height:400px">
      <div v-else>
        <div :span="13" :xs="24" class="text-center">code：{{ video_code }}</div>
        <video-player  class="video-player-box"
                       ref="videoPlayer"
                       :options="playerOptions"
                       :playsinline="true"
                       customEventName="customstatechangedeventname">
        </video-player>
      </div>
    </el-dialog>
    <el-dialog v-el-drag-dialog width="30%" :visible.sync="dialogTooltip" :title="dialogtt" @dragDialog="handleDrag">
      <div class="text-center">确定{{tooltipStatus|genderStatus}}该用户的请求？</div>
      <div class="text-center" style="padding-top:20px">
      <el-button type="danger" @click="dialogTooltip=false" style="margin-right:50px">取消</el-button>
      <el-button type="primary" @click="operation(tooltipID, tooltipStatus)">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
  import {getAudit, handleVideo} from '@/api/faceAudit'
  import waves from '@/directive/waves' // waves directive
  import Mallki from '@/components/TextHoverEffect/Mallki'
  import Pagination from '@/components/Pagination' // secondary package based on el-pagination
  import elDragDialog from '@/directive/el-drag-dialog'
  import { videoPlayer } from 'vue-video-player'

  import 'video.js/dist/video-js.css'

  // arr to obj, such as { CN : "China", US : "USA" }


  export default {
    name: 'ComplexTable',
    components: {Pagination, Mallki, videoPlayer},
    directives: {waves, elDragDialog},
    filters: {
      genderFilter(type) {
        const statusMap = {
          male: '男',
          female: '女'
        }
        return statusMap[type]
      },
      genderStatus(type) {
        const statusMap = {
          success: '通过',
          refuse: '拒绝'
        }
        return statusMap[type]
      }
    },
    data() {
      return {
        dialogTableVisible: false,
        dialogTitle:'',
        tableKey: 0,
        dialogtt: '',
        video_code:'',
        video_url:'',
        tooltipStatus: '',
        tooltipID: '',
        dialogTooltip: false,
        showWindows: {
            "type":'',
            "url":"",
        },
        playerOptions:{

        },
        list: null,
        total: 0,
        listLoading: true,
        images: [],
        listQuery: {
          page: 1,
          limit: 20,
          search: ''
        }
      }
    },
    created() {
      this.getList()
    },
    methods: {
      getList() {
        this.listLoading = true;
        getAudit(this.listQuery).then(response => {
          this.list = response.data.users;
          this.total = response.data.count;
          setTimeout(() => {
            this.listLoading = false;
          }, 0.5 * 1000)
        })
      },
      showDialog(row){
        this.showImage=row.image_url;
        this.dialogTableVisible = true;
      },
      showVideo(row){
          this.dialogTitle='审核视频';
          this.showWindows={"type": "video", "url": row.video_url}
          this.video_code=row.video_code
          this.playerOptions={
          muted: true,
            height: '360',
            width: '500',
            language: 'en',
            playbackRates: [0.7, 1.0, 1.5, 2.0],
            sources: [{
            type: "video/mp4",
            src: row.video_url
          }],
        },
        this.dialogTableVisible = true;
      },
      showImage(url){
        this.dialogTitle='用户图片';
        this.showWindows={"type": "image", "url": url}
        this.dialogTableVisible = true;
      },
      showTooltip(id, status){
          this.dialogtt='警告';
        this.tooltipStatus = status;
          this.tooltipID = id;
        this.dialogTooltip = true;
      },
      handleDrag() {
        console.log('handleDrag');
      },
      operation(user_id, status) {
        this.listLoading = true;

        handleVideo(user_id, status).then(response => {
          setTimeout(() => {
            this.$message({
              message: '操作成功',
              type: 'success'
            })

          }, 0.5 * 1000);
          this.getList();
        });

        this.listLoading = false;
        this.dialogTooltip=false;
      },
      handleFilter(){

      },
    }
  }
</script>
<style lang="scss" scoped>
.el-dialog{
  width: 30%;
}
</style>
