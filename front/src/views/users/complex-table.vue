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
      <el-table-column label="UID" align="center" width="80">
        <template slot-scope="scope">
          <router-link :to="'/users/'+scope.row._id" class="link-type">
            <span>{{ scope.row.ID }}</span>
          </router-link>
          <!--<el-button type="text" @click="userDetails(scope.row._id)">{{ scope.row.ID }}</el-button>-->
        </template>
      </el-table-column>
      <el-table-column label="邀请码" width="110px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.be_guided }}</span>
        </template>
      </el-table-column>
      <el-table-column label="注册时间" width="180px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.login }}</span>
        </template>
      </el-table-column>
      <el-table-column label="最后登录时间" width="180px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.last }}</span>
        </template>
      </el-table-column>
      <el-table-column label="昵称" width="180px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.nick_name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="性别" width="90px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.gender| genderFilter }}</span>
        </template>
      </el-table-column>
      <el-table-column label="当前状态" width="110px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.progress|statusFilter }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="180" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button size="mini" type="danger" @click="showDelete(row._id)" >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0"  :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
    <el-dialog v-el-drag-dialog width="30%" :visible.sync="dialogTooltip" title="警告" @dragDialog="handleDrag">
      <div class="text-center">如果删除该用户则该用户在本平台上发布的数据以及操作将被全部抹掉，确定删除该用户？</div>
      <div class="text-center" style="padding-top:20px">
        <el-button type="danger" @click="dialogTooltip=false" style="margin-right:50px">取消</el-button>
        <el-button type="primary" @click="deleteUser()">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { fetchList, delete_user } from '@/api/article'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Mallki from '@/components/TextHoverEffect/Mallki'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import elDragDialog from '@/directive/el-drag-dialog'



// arr to obj, such as { CN : "China", US : "USA" }


export default {
  name: 'ComplexTable',
  components: { Pagination, Mallki },
  directives: { waves, elDragDialog },
  filters: {
    statusFilter(status) {
      const statusMap = {
        audit: '待审核',
        certify: '未认证',
        finish: '已认证',
        refuse: '认证失败',
        profile: "待完善资料"
      }
      return statusMap[status]
    },
    genderFilter(type) {
      const statusMap = {
        male: '男',
        female: '女'
      }
      return statusMap[type]
    }
  },
  data() {
    return {
      tableKey: 0,
      list: null,
      total: 0,
      dialogTooltip: false,
      delete_id: '',
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        search: undefined,
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true;
      fetchList(this.listQuery).then(response => {
        this.list = response.data.users;
        this.total = response.data.count;
        setTimeout(() => {
          this.listLoading = false;
        }, 0.5 * 1000)
      })
    },
    showDelete(id){
      this.delete_id=id;
      this.dialogTooltip=true;
    },
    deleteUser() {
      this.listLoading = true;
      delete_user(this.delete_id).then(response => {
      setTimeout(() => {
        this.$message({
          message: '操作成功',
          type: 'success'
        })

    }, 0.5 * 1000);
        this.dialogTooltip=false;
        this.getList();
    });

      this.listLoading = false;
    },
    handleFilter() {
      this.getList()
    },
    sortByID(order) {
      if (order === 'ascending') {
        this.listQuery.sort = '+id'
      } else {
        this.listQuery.sort = '-id'
      }
      this.handleFilter()
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        importance: 1,
        remark: '',
        timestamp: new Date(),
        title: '',
        status: 'published',
        type: ''
      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row) // copy obj
      this.temp.timestamp = new Date(this.temp.timestamp)
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    handleDelete(row) {
      this.$notify({
        title: 'Success',
        message: 'Delete Successfully',
        type: 'success',
        duration: 2000
      })
      const index = this.list.indexOf(row)
      this.list.splice(index, 1)
    },
    formatJson(filterVal, jsonData) {
      return jsonData.map(v => filterVal.map(j => {
        if (j === 'timestamp') {
          return parseTime(v[j])
        } else {
          return v[j]
        }
      }))
    },
    handleDrag(){

    },
  }
}
</script>
