<template>
  <div class="app-container">
    <!--<div class="filter-container">-->
    <!--<el-input v-model="listQuery.search" placeholder="用户昵称" style="width: 200px;" class="filter-item"-->
    <!--@keyup.enter.native="handleFilter"/>-->
    <!--<el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">-->
    <!--Search-->
    <!--</el-button>-->
    <!--</div>-->

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column label="所属动态" align="center" width="260">
        <template slot-scope="scope">
        <router-link :to="'/moments/details/'+scope.row.moment_id" class="link-type">
          <span>{{ scope.row.moment_content }}</span>
        </router-link>
        </template>
      </el-table-column>
      <el-table-column label="评论内容" align="center" width="200">
        <template slot-scope="scope">
          <span>{{ scope.row.content }}</span>
        </template>
      </el-table-column>
      <el-table-column label="评论人" width="110px" align="center">
        <template slot-scope="scope">
          <router-link :to="'/users/'+scope.row.user_id" class="link-type">
            <span>{{ scope.row.nickname }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="评论时间" width="180px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.create_time }}</span>
        </template>
      </el-table-column>
      <el-table-column label="回复人" width="110px" align="center">
        <template slot-scope="scope">
          <router-link v-if="scope.row.reply_user_id" :to="'/users/'+scope.row.reply_user_id" class="link-type">
            <span>{{ scope.row.reply_user_nickname }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="180" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button size="mini" type="danger" @click="deleteComment(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit"
                @pagination="getList"/>
  </div>
</template>
<script>
  import {get_comments, delete_comments} from '@/api/article'
  import waves from '@/directive/waves' // waves directive
  import {parseTime} from '@/utils'
  import Mallki from '@/components/TextHoverEffect/Mallki'
  import Pagination from '@/components/Pagination' // secondary package based on el-pagination


  // arr to obj, such as { CN : "China", US : "USA" }


  export default {
    name: 'Comments',
    components: {Pagination},
    directives: {waves},
    filters: {},
    data() {
      return {
        dialogTableVisible: false,
        dialogTitle: '',
        page: 1,
        tableKey: 0,
        list: null,
        listLoading: true,
        total: 0,
        listQuery: {
          page: 1,
          limit: 20,
          search: undefined,
          id: ""
        }
      }
    },
    created() {
      if (Object.keys(this.$route.query).length>0) {
          console.log(this.$route.query);
          this.listQuery.id = this.$route.query.moment_id;
      }

      this.getList()
    },
    methods: {
      getList() {
        this.listLoading = true;
        get_comments(this.listQuery).then(response => {
          this.list = response.data.comments;
          this.total = response.data.count;
          setTimeout(() => {
            this.listLoading = false;
          }, 0.5 * 1000)
        })
      },
      showDialog(row){
        this.dialogTitle = row.content;
        this.images = row.images;
        this.dialogTableVisible = true;
      },
      deleteComment(row) {
        this.listLoading = true;
        delete_comments(row._id).then(response => {
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
