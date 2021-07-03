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
      <el-table-column label="昵称" align="center" width="130">
        <template slot-scope="scope">
          <router-link :to="'/users/'+scope.row.user_id" class="link-type">
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
          <span>{{ scope.row.create_time }}</span>
        </template>
      </el-table-column>
      <el-table-column label="提现微信号" width="110px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.wechat_code }}</span>
        </template>
      </el-table-column>
      <el-table-column label="转账时间" width="180px" align="center">
        <template slot-scope="scope">
          <span v-if="scope.row.finish">{{scope.row.finish_time}}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="230" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button v-if="!row.finish" type="primary" @click="handelMoney(row._id)">确认转账</el-button>
          <span v-else>已转账</span>
          <!--<el-button size="mini" type="danger" @click="showTooltip(row._id)">删除</el-button>-->
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" @pagination="getList" />
    <el-dialog v-el-drag-dialog width="30%" :visible.sync="dialogTooltip" :title="dialogTitle" @dragDialog="handleDrag">
      <div class="text-center">确定删除该用户的请求？</div>
      <div class="text-center" style="padding-top:20px">
        <el-button type="danger" @click="dialogTooltip=false" style="margin-right:50px">取消</el-button>
        <el-button type="primary" @click="deleteRow(deleteId)">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
  import {getRewards, handleReward, deleteReward} from '@/api/guider'
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
          'true': '已转账',
          finish: '完成',
          certify: '待认证'
        }
        return statusMap[status]
      },
    },
    data() {
      return {
        tableKey: 0,
        dialogTitle: '',
        deleteId: '',
        tooltipID: '',
        dialogTooltip: false,
        list: null,
        total: 0,
        listLoading: true,
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
        getRewards(this.listQuery).then(response => {
          this.list = response.data.record;
          this.total = response.data.count;
          setTimeout(() => {
            this.listLoading = false;
          }, 0.5 * 1000)
        })
      },
      showTooltip(id){
        this.dialogTitle='警告';
        this.deleteId=id;
        this.dialogTooltip = true;
      },
      deleteRow(id){
        deleteReward(id).then(response => {
          setTimeout(() => {
            this.$message({
              message: '操作成功',
              type: 'success'
            })

          }, 0.5 * 1000);
          this.getList();
        })
      },
      handelMoney(id){
        handleReward(id).then(response => {
          setTimeout(() => {
            this.$message({
              message: '操作成功',
              type: 'success'
            })

          }, 0.5 * 1000);
          this.getList();
        })
      },
      handleDrag() {
        console.log('handleDrag');
      },
    }
  }
</script>
<style lang="scss" scoped>
  .el-dialog{
    width: 30%;
  }
</style>
