<template>
  <div class="app-container">
    <div class="filter-container">
    <el-select v-model="listQuery.province_type" clearable placeholder="城市" style="width: 140px" class="filter-item" @change="handleFilter">
      <el-option v-for="item in provinces" :key="item" :label="item" :value="item" />
    </el-select>
    <el-select v-model="listQuery.constellation_type" clearable placeholder="星座" style="width: 140px" class="filter-item" @change="handleFilter">
      <el-option v-for="item in constellations" :key="item" :label="item" :value="item" />
    </el-select>
    <el-select v-model="listQuery.gender_type" clearable placeholder="性别" style="width: 140px" class="filter-item" @change="handleFilter">
      <el-option v-for="item in genders" :key="item.key" :label="item.label" :value="item.key" />
    </el-select>
      <el-button type="primary" @click="clearFilter">清除筛选</el-button>
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
      <el-table-column label="昵称" align="center" width="130">
        <template slot-scope="scope">
          <router-link :to="'/users/'+scope.row.user_id" class="link-type">
          <span>{{ scope.row.nick_name }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="性别" align="center" width="130">
        <template slot-scope="scope">
            <span>{{ scope.row.gender| genderFilter }}</span>
        </template>
      </el-table-column>
      <el-table-column label="颜值" align="center" width="50">
        <template slot-scope="scope">
          <span>{{ scope.row.face }}</span>
        </template>
      </el-table-column>
      <el-table-column label="图片" width="140px" align="center">
        <template slot-scope="scope">
          <img :src="scope.row.image_url+'?x-oss-process=image/resize,m_lfit,h_100,w_100/format,jpg'" alt="" style="width: 100px; height: 100px" @click="showDialog(scope.row)">
        </template>
      </el-table-column>
      <el-table-column label="上传时间" width="180px" align="center">
        <template slot-scope="scope">
          <span>{{scope.row.create_time}}</span>
        </template>
      </el-table-column>
      <el-table-column label="胜场次" width="80px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.result.win }}</span>
        </template>
      </el-table-column>
      <el-table-column label="总场次" width="80px" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.result.win+scope.row.result.fail }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="180" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button v-waves size="mini" type="danger" @click="showTooltip(row.image_id)">
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
      <div class="text-center">确定删除该图片？</div>
      <div class="text-center" style="padding-top:20px">
        <el-button type="danger" @click="dialogTooltip=false" style="margin-right:50px">取消</el-button>
        <el-button type="primary" @click="deleteImages()">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
  import {getFace, deleteImage, faceSelect} from '@/api/faceAudit'
  import waves from '@/directive/waves' // waves directive
  import Mallki from '@/components/TextHoverEffect/Mallki'
  import Pagination from '@/components/Pagination' // secondary package based on el-pagination
  import elDragDialog from '@/directive/el-drag-dialog'


  // arr to obj, such as { CN : "China", US : "USA" }


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
      }
    },
    data() {
      return {
        dialogTableVisible: false,
        dialogTitle:'',

        provinces: [],
        dialogtt: '',
        genders: [],
        constellations: [],
        dialogTooltip: false,
        tableKey: 0,
        tooltipID:'',
        showImage: "",
        list: null,
        total: 0,
        listLoading: true,
        images: [],
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
        getFace(this.listQuery).then(response => {
          this.list = response.data.users;
          this.total = response.data.count;
          this.listQuery.gender_type = response.data.match.gender;
          this.listQuery.constellation_type = response.data.match.constellation;
          this.listQuery.province_type = response.data.match.province;
        });
        faceSelect().then(response => {
          this.constellations = response.data.constellation;
          this.provinces = response.data.province;
          this.genders = [{"label": "男", "key": "male"}, {"label": "女", "key": "female"}]
          setTimeout(() => {
            this.listLoading = false;
            }, 0.5 * 1000)

        })

      },
      showDialog(row){
          this.showImage=row.image_url+'?x-oss-process=image/resize,m_lfit,h_300,w_300/format,jpg';
        this.dialogTableVisible = true;
      },
      handleDrag() {
        console.log('handleDrag');
      },
      clearFilter() {
        this.listQuery.province_type= '';
        this.listQuery.constellation_type= '';
        this.listQuery.gender_type= '';
        this.listQuery.page = 1;
        this.getList()
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
      deleteImages() {
        this.listLoading = true;
        deleteImage(this.tooltipID).then(response => {
          setTimeout(() => {
            this.$message({
              message: '操作成功',
              type: 'success'
            })

          }, 0.5 * 1000);
          this.dialogTooltip = false;
          this.getList();
        });

        this.listLoading = false;
      },
    }
  }
</script>
