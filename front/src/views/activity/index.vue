<template>
  <div class="app-container">
      <el-button type="primary" @click="createActivity">创建活动</el-button>

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column label="活动图片" align="center" width="230">
        <template slot-scope="scope">
          <img :src="scope.row.image_url" alt="" style="width: 100px; height: 100px" @click="showDialog(scope.row)">
        </template>
      </el-table-column>
      <el-table-column label="活动标题" align="center" width="180">
        <template slot-scope="scope">
          <span>{{scope.row.text}}</span>
        </template>
      </el-table-column>
      <el-table-column label="链接地址" align="center" width="230">
        <template slot-scope="scope">
          <a :href="scope.row.link" target="_blank" class="link-type">
            <span>{{ scope.row.link }}</span>
          </a>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" width="180">
        <template slot-scope="scope">
          <span>{{ scope.row.create_time }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="180" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button v-waves size="mini" type="danger" @click="showTooltip(row._id)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0"  :total="total" :page.sync="listQuery.page" @pagination="getList" />

    <el-dialog v-el-drag-dialog :visible.sync="dialogTableVisible" :title="dialogTitle" @dragDialog="handleDrag">
          <img :src="showImage" width="500px" height="500px">
    </el-dialog>
    <el-dialog v-el-drag-dialog width="30%" :visible.sync="dialogTooltip" :title="dialogtt" @dragDialog="handleDrag">
      <div class="text-center">确定删除该活动？</div>
      <div class="text-center" style="padding-top:20px">
        <el-button type="danger" @click="dialogTooltip=false" style="margin-right:50px">取消</el-button>
        <el-button type="primary" @click="deleteImages(tooltipID)">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
  import {getActivity, deleteActivity} from '@/api/user'
  import Pagination from '@/components/Pagination' // secondary package based on el-pagination
  import waves from '@/directive/waves' // waves directive
  import elDragDialog from '@/directive/el-drag-dialog'


  // arr to obj, such as { CN : "China", US : "USA" }


  export default {
    name: 'ComplexTable',
    components: {Pagination},
    directives: {elDragDialog, waves},
    data() {
      return {
        dialogTableVisible: false,
        dialogTitle:'',
        dialogtt: '',
        dialogTooltip: false,
        tableKey: 0,
        tooltipID:'',
        showImage: "",
        list: null,
        total: 0,
        listLoading: true,
        listQuery: {
          page: 1,
          limit: 20,
          search: '',
          province_type: '',
          constellation_type: '',
          gender_type: '',
        }
      }
    },
    created() {
      this.getList()
    },
    methods: {
      getList() {
        this.listLoading = true;
        getActivity(this.listQuery).then(response => {
          this.list = response.data.record;
          this.total = response.data.count;
          this.listLoading = false;
        });
      },
      createActivity(){
        this.$router.push("/activity/create")
      },
      showDialog(row){
          this.showImage=row.image_url+'?x-oss-process=image/resize,m_lfit,h_200,w_200/format,jpg';
        this.dialogTableVisible = true;
      },
      handleDrag() {
        console.log('handleDrag');
      },
      handleFilter() {
        this.listQuery.page = 1
        this.getList()
      },
      showTooltip(id){
        this.dialogtt='警告';
        this.tooltipID = id;
        this.dialogTooltip = true;
      },
      deleteImages(id) {
        this.listLoading = true;
        deleteActivity(id).then(response => {
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
    }
  }
</script>
