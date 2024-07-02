<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="应用代码">
          <el-input v-model.trim="params.appName" clearable placeholder="应用名" @clear="search" />
        </el-form-item>
        <el-form-item label="应用别名">
          <el-input v-model.trim="params.category" clearable placeholder="应用别名" @clear="search" />
        </el-form-item>
        <el-form-item label="中文名">
          <el-input v-model.trim="params.title" clearable placeholder="中文名" @clear="search" />
        </el-form-item>
        <el-form-item label="baseUrl">
          <el-input v-model.trim="params.baseUrl" clearable placeholder="baseUrl" @clear="search" />
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">新增</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger" @click="batchDelete">批量删除</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="appName" label="应用名" />
        <el-table-column show-overflow-tooltip sortable prop="alias" label="应用别名" />
        <el-table-column show-overflow-tooltip sortable prop="title" label="中文名" />
        <el-table-column show-overflow-tooltip sortable prop="baseUrl" label="baseUrl" />
        <el-table-column show-overflow-tooltip sortable prop="introduction" label="说明" />
        <el-table-column fixed="right" label="操作" align="center" width="120">
          <template slot-scope="scope">
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip class="delete-popover" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.ID)">
                <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" />
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        :current-page="params.pageNum"
        :page-size="params.pageSize"
        :total="total"
        :page-sizes="[1, 5, 10, 30]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px;float:right;margin-bottom: 10px;"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />

      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="120px">
          <el-form-item label="访问路径" prop="path">
            <el-input v-model.trim="dialogFormData.appName" placeholder="应用代码" />
          </el-form-item>
          <el-form-item label="别名" prop="category">
            <el-input v-model.trim="dialogFormData.alias" placeholder="别名" />
          </el-form-item>
          <el-form-item label="所属类别" prop="desc">
            <el-select v-model.trim="dialogFormData.typeId" multiple placeholder="请选择类别" style="width:100%">
              <el-option
                v-for="item in typeList"
                :key="item.ID"
                :label="item.desc"
                :value="item.ID"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="请求方式" prop="method">
            <el-input v-model.trim="dialogFormData.title" placeholder="中文名" />
          </el-form-item>
          <el-form-item label="URL" prop="method">
            <el-input v-model.trim="dialogFormData.baseUrl" placeholder="baseUrl" />
          </el-form-item>
          <el-form-item label="说明" prop="desc">
            <el-input v-model.trim="dialogFormData.introduction" type="textarea" placeholder="说明" show-word-limit maxlength="100" />
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelForm()">取 消</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div>
      </el-dialog>

    </el-card>
  </div>
</template>

<script>
import { getApps, createApp, updateAppById, batchDeleteAppByIds, getAppTypes } from '@/api/system/app'

export default {
  name: 'App',
  filters: {
    methodTagFilter(val) {
      if (val === 'GET') {
        return ''
      } else if (val === 'POST') {
        return 'success'
      } else if (val === 'PUT') {
        return 'info'
      } else if (val === 'PATCH') {
        return 'warning'
      } else if (val === 'DELETE') {
        return 'danger'
      } else {
        return 'info'
      }
    }
  },
  data() {
    return {
      // 查询参数
      params: {
        appName: '',
        alias: '',
        title: '',
        baseUrl: '',
        pageNum: 1,
        pageSize: 10
      },
      // 表格数据
      tableData: [],
      total: 0,
      loading: false,

      typeList: [],

      // dialog对话框
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        appName: '',
        alias: '',
        title: '',
        baseUrl: '',
        typeId: '',
        introduction: ''
      },
      dialogFormRules: {
        appName: [
          { required: true, message: '请输入应用名', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        alias: [
          { required: true, message: '请输入应用别名', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        title: [
          { required: true, message: '请输入中文名', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        baseUrl: [
          { required: true, message: '请输入Url地址', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        introduction: [
          { required: false, message: '说明', trigger: 'blur' },
          { min: 0, max: 100, message: '长度在 0 到 128 个字符', trigger: 'blur' }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: []
    }
  },
  created() {
    this.getAppTypes()
    this.getTableData()
  },
  methods: {
    // 查询
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const resp = await getApps(this.params)
        console.log('resp====', JSON.stringify(resp))
        this.tableData = resp.data.list
        this.total = resp.data.total
      } finally {
        this.loading = false
      }
    },

    // 获取表格数据
    async getAppTypes() {
      this.loading = true
      this.params.typeId = 0
      try {
        const { data } = await getAppTypes()
        this.typeList = data
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增接口'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData.ID = row.ID
      this.dialogFormData.appName = row.appName
      this.dialogFormData.alias = row.alias
      this.dialogFormData.title = row.title
      this.dialogFormData.baseUrl = row.baseUrl
      this.dialogFormData.introduction = row.introduction

      this.dialogFormTitle = '修改应用'
      this.dialogType = 'update'
      this.dialogFormVisible = true
    },

    // 提交表单
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          let msg = ''
          this.submitLoading = true
          try {
            if (this.dialogType === 'create') {
              const { message } = await createApp(this.dialogFormData)
              msg = message
            } else {
              const { message } = await updateAppById(this.dialogFormData.ID, this.dialogFormData)
              msg = message
            }
          } finally {
            this.submitLoading = false
          }

          this.resetForm()
          this.getTableData()
          this.$message({
            showClose: true,
            message: msg,
            type: 'success'
          })
        } else {
          this.$message({
            showClose: true,
            message: '表单校验失败',
            type: 'error'
          })
          return false
        }
      })
    },

    // 提交表单
    cancelForm() {
      this.resetForm()
    },

    resetForm() {
      this.dialogFormVisible = false
      this.$refs['dialogForm'].resetFields()
      this.dialogFormData = {
        path: '',
        category: '',
        method: '',
        desc: ''
      }
    },

    // 批量删除
    batchDelete() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const apiIds = []
        this.multipleSelection.forEach(x => {
          apiIds.push(x.ID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteAppByIds({ apiIds: apiIds })
          msg = message
        } finally {
          this.loading = false
        }

        this.getTableData()
        this.$message({
          showClose: true,
          message: msg,
          type: 'success'
        })
      }).catch(() => {
        this.$message({
          showClose: true,
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
    },

    // 单个删除
    async singleDelete(Id) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await batchDeleteAppByIds({ apiIds: [Id] })
        msg = message
      } finally {
        this.loading = false
      }

      this.getTableData()
      this.$message({
        showClose: true,
        message: msg,
        type: 'success'
      })
    },

    // 分页
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getTableData()
    }
  }
}
</script>

<style scoped>
  .container-card{
    margin: 10px;
  }

  .delete-popover{
    margin-left: 10px;
  }
</style>
