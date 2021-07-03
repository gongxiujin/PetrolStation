<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="listQuery.search" placeholder="用户昵称" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        Search
      </el-button>
    </div>
    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column label="ID" align="center" width="110">
        <template slot-scope="scope">
            <span>{{ scope.row.ID }}</span>
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
      <el-table-column label="已邀请人数" width="110px" align="center">
        <template slot-scope="scope">
          <el-button type="text" @click="showGuider(scope.row.guide)">{{ scope.row.guided_count }}</el-button>
        </template>
      </el-table-column>
      <!--<el-table-column label="操作" align="center" width="180" class-name="small-padding fixed-width">-->
        <!--<template slot-scope="{row}">-->
          <!--<router-link :to="'/guider/'+scope.row._id" class="link-type">-->
          <!--<span>查看推荐奖励</span>-->
          <!--</router-link>-->
        <!--</template>-->
      <!--</el-table-column>-->
    </el-table>

    <pagination v-show="total>0"  :total="total" :page.sync="listQuery.page" @pagination="getList" />

    <el-dialog v-el-drag-dialog :visible.sync="dialogTableVisible" title="推荐的用户" @dragDialog="handleDrag">
      <div >
        <el-table
          :key="tableKey"
          v-loading="listLoading"
          :data="guider"
          border
          fit
          highlight-current-row
          style="width: 100%;"
        >
          <el-table-column label="ID" align="center" width="110">
            <template slot-scope="scope">
              <span>{{ scope.row.ID }}</span>
            </template>
          </el-table-column>
          <el-table-column label="昵称" align="center" width="110">
            <template slot-scope="scope">
              <router-link :to="'/users/'+scope.row._id" class="link-type">
              <span>{{ scope.row.nick_name }}</span>
              </router-link>
            </template>
          </el-table-column>
          <el-table-column label="注册时间" align="center" width="110">
            <template slot-scope="scope">
              <span>{{ scope.row.create_time }}</span>
            </template>
          </el-table-column>
          <el-table-column label="认证" align="center" width="110">
            <template slot-scope="scope">
              <span>{{ scope.row.progress|statusFilter }}</span>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script>
  import {getGuider, beGuider} from '@/api/guider'
  import waves from '@/directive/waves' // waves directive
  import Mallki from '@/components/TextHoverEffect/Mallki'
  import Pagination from '@/components/Pagination' // secondary package based on el-pagination
  import elDragDialog from '@/directive/el-drag-dialog'




  export default {
    name: 'ComplexTable',
    components: {Pagination, Mallki},
    directives: {waves, elDragDialog},
    filters: {
      genderFilter(type) {
        const statusMap = {
          male: '男',
          female: '女'
        }
        return statusMap[type]
      },
      statusFilter(status) {
        const statusMap = {
          audit: '待审核',
          finish: '完成',
          certify: '待认证'
        }
        return statusMap[status]
      },
    },
    data() {
      return {
        dialogTableVisible: false,
        dialogTitle:'',
        tableKey: 0,
        guider: null,
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
      console.log(this.listQuery.limit);
        this.listLoading = true;
        getGuider(this.listQuery).then(response => {
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
      showGuider(guided){
        this.guider = guided;
        this.dialogTableVisible=true;
      },
      handleDrag() {
        console.log('handleDrag');
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
