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
      <el-table-column label="举报原因内容" align="center" width="200">
        <template slot-scope="scope">
          <span></span>
        </template>
      </el-table-column>
      <el-table-column label="举报人" width="110px" align="center">
        <template slot-scope="scope">
          <router-link :to="'/users/'+scope.row.report_user_id" class="link-type">
            <span>{{ scope.row.nickname }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="举报内容" width="110px" align="center">
        <template slot-scope="scope">
          <router-link :to="'/moments/details/'+scope.row.moment_id" class="link-type">
            <span>{{ scope.row.content }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="举报时间" width="180px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.create_time }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="110px" align="center">
        <template slot-scope="scope">
          <router-link v-if="scope.row.reply_user_id" :to="'/users/'+scope.row.reply_user_id" class="link-type">
            <span>{{ scope.row.reply_user_nickname }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="操作/处理结果" align="center" width="180" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button v-if="row.status===1" size="mini" type="text" @click="showDialog(row)">
            处理
          </el-button>
          <span v-else>
            <span>{{row.result}}</span>
          </span>
        </template>
      </el-table-column>
    </el-table>
    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit"
                @pagination="getList"/>
    <el-dialog v-el-drag-dialog :visible.sync="dialogTableVisible" :title="dialogTitle" @dragDialog="handleDrag">
      <el-form>
      <el-form-item v-if="!status" label="处理意见">
        <el-input
          v-model="opinion"
          :autosize="{ minRows: 2, maxRows: 4}"
          type="textarea"
          placeholder=""
        />
          <el-row type="flex" class="row-bg" justify="center" style="padding-top: 20px">
            <el-col :span="6"><el-button v-waves type="primary" @click="handelReports">提交</el-button></el-col>
          </el-row>
      </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>
<script>
  import {get_reports, handle_reports} from '@/api/article'
  import waves from '@/directive/waves' // waves directive
  import {parseTime} from '@/utils'
  import Mallki from '@/components/TextHoverEffect/Mallki'
  import Pagination from '@/components/Pagination' // secondary package based on el-pagination
  import elDragDialog from '@/directive/el-drag-dialog'


  export default {
    name: 'Comments',
    components: {Pagination},
    directives: {waves, elDragDialog},
    filters: {},
    data() {
      return {
        dialogTableVisible: false,
        dialogTitle: '',
        page: 1,
        tableKey: 0,
        list: null,
        opinion: null,
        status:true,
        select_id: '',
        listLoading: true,
        total: 0,
        listQuery: {
          page: 1,
          limit: 20,
        }
      }
    },
    created() {
      this.getList()
    },
    methods: {
      getList() {
        this.listLoading = true;
        get_reports(this.listQuery).then(response => {
          this.list = response.data.reports;
          this.total = response.data.count;
          setTimeout(() => {
            this.listLoading = false;
          }, 0.5 * 1000)
        })
      },

      showDialog(row){
        this.select_id = row._id;
        this.status= row.result===null
        this.opinion = row.result;
        this.dialogTableVisible = true;
      },
      handelReports() {
        this.listLoading = true;
        handle_reports(this.select_id, this.opinion).then(response => {
          setTimeout(() => {
            this.$message({
              message: '操作成功',
              type: 'success'
            })

          }, 0.5 * 1000);
          this.dialogTableVisible = false;
          this.getList();
        });

        this.listLoading = false;
      },
      handleDrag(){

      }
    }
  }
</script>
