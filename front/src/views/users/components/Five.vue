<template>
    <el-form :inline="true" :model="formInline" class="demo-form-inline" v-loading="listLoading" >
      <el-form-item :label="user.gender==='male'?'帅气':'美丽'">
        <el-select v-model="user.five.cool" placeholder="数值" style="width: 90px;">
          <el-option
            v-for="option in options"
            :key="option.value"
            :label="option.label"
            :value="option.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="性感">
        <el-select v-model="user.five.sexy" placeholder="数值" style="width: 90px;">
          <el-option
            v-for="option in options"
            :key="option.value"
            :label="option.label"
            :value="option.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="user.gender==='male'?'阳光':'可爱'">
        <el-select v-model="user.five.character" placeholder="数值" style="width: 90px;">
          <el-option
            v-for="option in options"
            :key="option.value"
            :label="option.label"
            :value="option.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="智慧">
        <el-select v-model="user.five.wisdom" placeholder="数值" style="width: 90px;">
          <el-option
            v-for="option in options"
            :key="option.value"
            :label="option.label"
            :value="option.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="user.gender==='male'?'多金':'高贵'">
        <el-select v-model="user.five.wealth" placeholder="数值" style="width: 90px;">
          <el-option
            v-for="option in options"
            :key="option.value"
            :label="option.label"
            :value="option.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">提交</el-button>
      </el-form-item>
    </el-form>

</template>

<script>
  import {update_user_figure} from '@/api/article'

  export default {
    props: {
      user: {
        type: Object,
        default: () => {
          return {
            _id: '',
            email: '',
            five:{
                cool: 0,
              sexy: 0,
              character: 0,
              wisdom: 0,
              wealth: 0}
          }
        }
      }
    },
    data(){
      return {
        options:  [{label:'1', value:1},
          {label:'2', value:2}, {label:'3', value:3},
          {label:'4', value:4}, {label:'5', value:5}],
        listLoading: false,
        formInline: {
          cool: 0,
          sexy: 0,
          character: 0,
          wisdom: 0,
          wealth: 0
        }
      }
    },
    methods: {
      onSubmit() {
        this.listLoading = true;
        update_user_figure(this.user.user_id, this.user.five).then(response => {
          setTimeout(() => {
            this.$message({
              message: '操作成功',
              type: 'success'
            });
            this.listLoading = false;
          }, 0.5 * 1000);
          this.$router.push(`/users/${this.user.user_id}`);
        });

      }
    }
  }
</script>
