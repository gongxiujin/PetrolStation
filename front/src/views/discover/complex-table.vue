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
      <el-table-column label="内容" align="center" width="280">
        <template slot-scope="scope">
          <router-link :to="'/moments/details/'+scope.row._id" class="link-type">
          <span>{{ scope.row.content }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="照片" align="center" width="100">
        <template slot-scope="scope">
          <el-button v-if="scope.row.images.length>0" type="text" @click="showDialog(scope.row)">
            查看照片
          </el-button>
        </template>
      </el-table-column>
      <el-table-column label="发布人" width="110px" align="center">
        <template slot-scope="scope">
          <router-link :to="'/users/'+scope.row.user_id" class="link-type">
            <span>{{ scope.row.nickname }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="发布时间" width="180px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.create_time }}</span>
        </template>
      </el-table-column>
      <el-table-column label="点赞" width="60px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.like }}</span>
        </template>
      </el-table-column>
      <el-table-column label="评论" width="60px" align="center">
        <template slot-scope="scope">
          <router-link v-if="scope.row.comment>0" :to="'/moments/comments?moment_id='+scope.row._id" class="link-type">
            <span>{{ scope.row.comment }}</span>
          </router-link>
          <span v-else>{{ scope.row.comment }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="180" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button size="mini" type="danger" @click="deleteMoment(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0"  :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog v-el-drag-dialog :visible.sync="dialogTableVisible" :title="dialogTitle" @dragDialog="handleDrag">
      <el-carousel :interval="6000" height="500px" type="card" >
        <el-carousel-item v-for="img in images">
          <img :src="img.thumbnail_url" class="image" style="width: 100%;height: auto;">
        </el-carousel-item>
      </el-carousel>
    </el-dialog>
  </div>
</template>

<script>
  import {moments, delete_moment} from '@/api/article'
  import waves from '@/directive/waves' // waves directive
  import {parseTime} from '@/utils'
  import Mallki from '@/components/TextHoverEffect/Mallki'
  import Pagination from '@/components/Pagination' // secondary package based on el-pagination
  import elDragDialog from '@/directive/el-drag-dialog'


  // arr to obj, such as { CN : "China", US : "USA" }


  export default {
    name: 'ComplexTable',
    components: {Pagination, Mallki},
    directives: {waves, elDragDialog},
    filters: {
      statusFilter(status) {
        const statusMap = {
          audit: '待审核',
          finish: '完成',
          certify: '待认证'
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
        dialogTableVisible: false,
        dialogTitle:'',
        tableKey: 0,
        list: null,
        total: 0,
        listLoading: true,
        images: [],
        listQuery: {
          page: 1,
          limit: 20,
          search: undefined
        }
      }
    },
    created() {
      this.getList()
    },
    methods: {
      getList() {
        this.listLoading = true;
        moments(this.listQuery).then(response => {
          this.list = response.data.moments;
          this.total = response.data.count;
          setTimeout(() => {
            this.listLoading = false;
          }, 0.5 * 1000)
        })
      },
      showDialog(row){
        this.dialogTitle=row.content;
        this.images=row.images;
        this.dialogTableVisible = true;
      },
      handleDrag() {
          console.log('handleDrag');
//        this.$refs.select.blur();
      },
      deleteMoment(row) {
        this.listLoading = true;
        delete_moment(row._id).then(response => {
          setTimeout(() => {
            this.$message({
              message: '操作成功',
              type: 'success'
            })

          }, 0.5 * 1000);
          this.getList();
        });

        this.listLoading = false;
      },
    }
  }
</script>
