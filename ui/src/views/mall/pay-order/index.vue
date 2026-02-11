<template>
  <div class="pay-order-container p-6">
    <n-card title="支付订单管理" class="order-card">
      <!-- 头部操作栏 -->
      <div class="header-section">
        <div class="search-section">
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索订单号、商户订单号或标题"
            clearable
            style="width: 250px"
            @input="handleSearch"
          >
            <template #prefix>
              <n-icon><search-outline /></n-icon>
            </template>
          </n-input>
          <n-select
            v-model:value="filterStatus"
            placeholder="筛选状态"
            clearable
            style="width: 150px; margin-left: 12px"
            :options="statusOptions"
            @update:value="handleFilterChange"
          />
          <n-select
            v-model:value="filterChannel"
            placeholder="筛选支付渠道"
            clearable
            style="width: 150px; margin-left: 12px"
            :options="channelOptions"
            @update:value="handleFilterChange"
          />
          <n-date-picker
            v-model:value="dateRange"
            type="daterange"
            clearable
            style="width: 240px; margin-left: 12px"
            placeholder="选择日期范围"
            @update:value="handleFilterChange"
          />
        </div>
        <div class="action-section">
          <n-button type="primary" @click="exportOrders" style="margin-right: 12px">
            <template #icon>
              <n-icon><download-outline /></n-icon>
            </template>
            导出订单
          </n-button>
          <n-button @click="handleRefresh">
            <template #icon>
              <n-icon><refresh-outline /></n-icon>
            </template>
            刷新
          </n-button>
        </div>
      </div>

      <!-- 统计信息卡片 -->
      <div class="stats-section">
        <n-grid :cols="4" :x-gap="16" :y-gap="16">
          <n-gi>
            <n-statistic label="今日订单" :value="todayStats.total">
              <template #prefix>
                <n-icon color="#52c41a"><document-text-outline /></n-icon>
              </template>
            </n-statistic>
          </n-gi>
          <n-gi>
            <n-statistic label="今日收入" :value="todayStats.amount" precision="2">
              <template #prefix>
                <n-icon color="#1890ff"><cash-outline /></n-icon>
              </template>
              <template #suffix>元</template>
            </n-statistic>
          </n-gi>
          <n-gi>
            <n-statistic label="成功率" :value="todayStats.successRate" precision="2">
              <template #prefix>
                <n-icon color="#722ed1"><trending-up-outline /></n-icon>
              </template>
              <template #suffix>%</template>
            </n-statistic>
          </n-gi>
          <n-gi>
            <n-statistic label="待处理" :value="todayStats.pending">
              <template #prefix>
                <n-icon color="#fa8c16"><time-outline /></n-icon>
              </template>
            </n-statistic>
          </n-gi>
        </n-grid>
      </div>

      <!-- 订单列表表格 -->
      <n-data-table
        :columns="columns"
        :data="filteredOrderList"
        :loading="loading"
        :pagination="pagination"
        :row-key="(row) => row.id"
        :scroll-x="1500"
        :remote="true"
        class="order-table"
      />
    </n-card>

    <!-- 订单详情模态框 -->
    <n-modal v-model:show="showDetailModal" preset="dialog" title="订单详情" style="width: 800px">
      <div class="order-detail-content" v-if="currentOrder">
        <n-tabs type="line" animated>
          <n-tab-pane name="basic" tab="基本信息">
            <n-descriptions :columns="2" bordered>
              <n-descriptions-item label="订单ID">{{ currentOrder.id }}</n-descriptions-item>
              <n-descriptions-item label="商户订单号">{{
                currentOrder.merchantOrderNo
              }}</n-descriptions-item>
              <n-descriptions-item label="订单标题">{{ currentOrder.title }}</n-descriptions-item>
              <n-descriptions-item label="订单金额">
                <n-text type="primary" strong>¥{{ currentOrder.amount }}</n-text>
              </n-descriptions-item>
              <n-descriptions-item label="实际支付">
                <n-text type="success" strong>¥{{ currentOrder.paidAmount }}</n-text>
              </n-descriptions-item>
              <n-descriptions-item label="支付渠道">{{
                currentOrder.channelName
              }}</n-descriptions-item>
              <n-descriptions-item label="订单状态">
                <n-tag :type="getStatusType(currentOrder.status)">{{
                  getStatusText(currentOrder.status)
                }}</n-tag>
              </n-descriptions-item>
              <n-descriptions-item label="创建时间">{{
                currentOrder.createdAt
              }}</n-descriptions-item>
              <n-descriptions-item label="支付时间" v-if="currentOrder.paidAt">{{
                currentOrder.paidAt
              }}</n-descriptions-item>
              <n-descriptions-item label="过期时间">{{
                currentOrder.expireAt
              }}</n-descriptions-item>
              <n-descriptions-item label="用户ID">{{ currentOrder.userId }}</n-descriptions-item>
            </n-descriptions>
          </n-tab-pane>

          <n-tab-pane name="payment" tab="支付信息">
            <n-descriptions :columns="2" bordered>
              <n-descriptions-item label="支付流水号" v-if="currentOrder.tradeNo">{{
                currentOrder.tradeNo
              }}</n-descriptions-item>
              <n-descriptions-item label="第三方订单号" v-if="currentOrder.thirdPartyOrderNo">{{
                currentOrder.thirdPartyOrderNo
              }}</n-descriptions-item>
              <n-descriptions-item label="支付状态">{{
                getStatusText(currentOrder.status)
              }}</n-descriptions-item>
              <n-descriptions-item label="支付时间" v-if="currentOrder.paidAt">{{
                currentOrder.paidAt
              }}</n-descriptions-item>
              <n-descriptions-item label="支付金额" v-if="currentOrder.paidAmount"
                >¥{{ currentOrder.paidAmount }}</n-descriptions-item
              >
              <n-descriptions-item label="手续费" v-if="currentOrder.fee"
                >¥{{ currentOrder.fee }}</n-descriptions-item
              >
              <n-descriptions-item label="失败原因" v-if="currentOrder.failReason" :span="2">
                <n-text type="error">{{ currentOrder.failReason }}</n-text>
              </n-descriptions-item>
            </n-descriptions>
          </n-tab-pane>

          <n-tab-pane name="extra" tab="扩展信息">
            <n-descriptions :columns="1" bordered>
              <n-descriptions-item label="客户端IP">{{
                currentOrder.clientIp
              }}</n-descriptions-item>
              <n-descriptions-item label="用户代理">{{
                currentOrder.userAgent
              }}</n-descriptions-item>
              <n-descriptions-item label="回调地址">{{
                currentOrder.notifyUrl
              }}</n-descriptions-item>
              <n-descriptions-item label="跳转地址">{{
                currentOrder.returnUrl
              }}</n-descriptions-item>
              <n-descriptions-item label="扩展参数" v-if="currentOrder.extra">
                <n-code :code="formatExtra(currentOrder.extra)" language="json" :word-wrap="true" />
              </n-descriptions-item>
            </n-descriptions>
          </n-tab-pane>
        </n-tabs>
      </div>
      <template #action>
        <n-space justify="end">
          <n-button @click="showDetailModal = false">关闭</n-button>
          <n-button
            v-if="currentOrder?.status === 'pending'"
            type="warning"
            @click="handleCancelOrder"
          >
            取消订单
          </n-button>
          <n-button v-if="currentOrder?.status === 'paid'" type="error" @click="handleRefundOrder">
            退款
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { NButton, NIcon, NTag, NText, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import {
  SearchOutline,
  RefreshOutline,
  DownloadOutline,
  DocumentTextOutline,
  CashOutline,
  TrendingUpOutline,
  TimeOutline,
  EyeOutline,
} from '@vicons/ionicons5'
import * as payOrderApi from '@/api/mall/payOrder'

const message = useMessage()

// 搜索和筛选
const searchKeyword = ref('')
const filterStatus = ref<string | null>(null)
const filterChannel = ref<string | null>(null)
const dateRange = ref<[number, number] | null>(null)

// 表格数据
const loading = ref(false)
const orderList = ref<PayOrderItem[]>([])

// 模态框状态
const showDetailModal = ref(false)
const currentOrder = ref<PayOrderItem | null>(null)

// 分页配置
const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  total: 0,
  onChange: (page: number) => {
    pagination.page = page
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
  },
})

// 订单状态选项
const statusOptions = [
  { label: '全部状态', value: null as any },
  { label: '待支付', value: 'pending' },
  { label: '支付成功', value: 'paid' },
  { label: '支付失败', value: 'failed' },
  { label: '已取消', value: 'cancelled' },
  { label: '已退款', value: 'refunded' },
]

// 支付渠道选项
const channelOptions = [
  { label: '全部渠道', value: null as any },
  { label: '支付宝', value: 'alipay' },
  { label: '微信支付', value: 'wechatpay' },
  { label: '银联支付', value: 'unionpay' },
  { label: 'PayPal', value: 'paypal' },
  { label: 'Stripe', value: 'stripe' },
]

// 今日统计
const todayStats = reactive({
  total: 0,
  amount: 0,
  successRate: 0,
  pending: 0,
})

// 订单数据类型
interface PayOrderItem {
  id: number
  merchantOrderNo: string
  title: string
  amount: number
  paidAmount: number
  fee: number
  channel: string
  channelName: string
  status: 'pending' | 'paid' | 'failed' | 'cancelled' | 'refunded'
  tradeNo?: string
  thirdPartyOrderNo?: string
  userId: string
  clientIp: string
  userAgent: string
  notifyUrl: string
  returnUrl: string
  extra?: string
  failReason?: string
  createdAt: string
  paidAt?: string
  expireAt: string
  updatedAt: string
}

// 后端订单数据类型
interface BackendPayOrder {
  id: number
  created_at: string
  updated_at: string
  channel_type: string
  order_id: string
  out_trade_no: string
  total_fee: string
  subject: string
  body: string
  notify_url: string
  return_url?: string
  extra?: string
  pay_url?: string
  state: string
  error_msg?: string
  raw?: string
}

// 状态映射
const stateToStatusMap: Record<string, PayOrderItem['status']> = {
  '0': 'cancelled',
  '1': 'pending',
  '2': 'paid',
  '3': 'failed',
  '4': 'refunded'
}

// 渠道映射
const channelToNameMap: Record<string, string> = {
  'alipay': '支付宝',
  'wechatpay': '微信支付',
  'unionpay': '银联支付',
  'paypal': 'PayPal',
  'stripe': 'Stripe'
}

// 将后端数据映射到前端数据
const mapBackendToFrontend = (backendOrder: BackendPayOrder): PayOrderItem => {
  const status = stateToStatusMap[backendOrder.state] || 'pending'
  const channelName = channelToNameMap[backendOrder.channel_type] || backendOrder.channel_type

  return {
    id: backendOrder.id,
    merchantOrderNo: backendOrder.out_trade_no || '',
    title: backendOrder.subject || '',
    amount: parseFloat(backendOrder.total_fee || '0'),
    paidAmount: status === 'paid' ? parseFloat(backendOrder.total_fee || '0') : 0,
    fee: 0,
    channel: backendOrder.channel_type,
    channelName: channelName,
    status: status,
    tradeNo: backendOrder.order_id,
    thirdPartyOrderNo: '',
    userId: '',
    clientIp: '',
    userAgent: '',
    notifyUrl: backendOrder.notify_url || '',
    returnUrl: backendOrder.return_url || '',
    extra: backendOrder.extra,
    failReason: backendOrder.error_msg,
    createdAt: backendOrder.created_at,
    paidAt: status === 'paid' ? backendOrder.updated_at : undefined,
    expireAt: '',
    updatedAt: backendOrder.updated_at
  }
}

// 筛选后的订单列表
const filteredOrderList = computed(() => {
  let filtered = orderList.value

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(
      (order) =>
        String(order.id).toLowerCase().includes(keyword) ||
        order.merchantOrderNo.toLowerCase().includes(keyword) ||
        order.title.toLowerCase().includes(keyword),
    )
  }

  // 状态筛选
  if (filterStatus.value) {
    filtered = filtered.filter((order) => order.status === filterStatus.value)
  }

  // 渠道筛选
  if (filterChannel.value) {
    filtered = filtered.filter((order) => order.channel === filterChannel.value)
  }

  // 日期范围筛选
  if (dateRange.value && dateRange.value.length === 2) {
    const [startDate, endDate] = dateRange.value
    filtered = filtered.filter((order) => {
      const orderDate = new Date(order.createdAt).getTime()
      return orderDate >= startDate && orderDate <= endDate
    })
  }

  return filtered
})

// 获取状态标签类型
const getStatusType = (status: string) => {
  const typeMap: Record<string, any> = {
    pending: 'warning',
    paid: 'success',
    failed: 'error',
    cancelled: 'default',
    refunded: 'info',
  }
  return typeMap[status] || 'default'
}

// 获取状态文本
const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    pending: '待支付',
    paid: '支付成功',
    failed: '支付失败',
    cancelled: '已取消',
    refunded: '已退款',
  }
  return textMap[status] || status
}

// 格式化扩展信息
const formatExtra = (extra: string) => {
  try {
    const parsed = JSON.parse(extra)
    return JSON.stringify(parsed, null, 2)
  } catch {
    return extra
  }
}

// 表格列定义
const columns: DataTableColumns<PayOrderItem> = [
  {
    title: '订单ID',
    key: 'id',
    width: 180,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '商户订单号',
    key: 'merchantOrderNo',
    width: 180,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '订单标题',
    key: 'title',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '订单金额',
    key: 'amount',
    width: 100,
    render: (row) => {
      return h(NText, { type: 'primary', strong: true }, () => `¥${row.amount}`)
    },
    sorter: (row1, row2) => row1.amount - row2.amount,
  },
  {
    title: '支付金额',
    key: 'paidAmount',
    width: 100,
    render: (row) => {
      if (row.paidAmount > 0) {
        return h(NText, { type: 'success', strong: true }, () => `¥${row.paidAmount}`)
      }
      return h(NText, { type: 'default' }, () => '-')
    },
    sorter: (row1, row2) => row1.paidAmount - row2.paidAmount,
  },
  {
    title: '支付渠道',
    key: 'channelName',
    width: 120,
    render: (row) => {
      return h(NTag, { type: 'info' }, () => row.channelName)
    },
  },
  {
    title: '订单状态',
    key: 'status',
    width: 100,
    render: (row) => {
      return h(NTag, { type: getStatusType(row.status) }, () => getStatusText(row.status))
    },
  },
  {
    title: '创建时间',
    key: 'createdAt',
    width: 160,
    sorter: (row1, row2) => new Date(row1.createdAt).getTime() - new Date(row2.createdAt).getTime(),
  },
  {
    title: '支付时间',
    key: 'paidAt',
    width: 160,
    render: (row) => {
      return row.paidAt || '-'
    },
    sorter: (row1, row2) => {
      const time1 = row1.paidAt ? new Date(row1.paidAt).getTime() : 0
      const time2 = row2.paidAt ? new Date(row2.paidAt).getTime() : 0
      return time1 - time2
    },
  },
  {
    title: '操作',
    fixed:"right",
    key: 'actions',
    width: 100,
    render: (row) => {
      return h(
        NButton,
        {
          size: 'small',
          type: 'primary',
          quaternary: true,
          onClick: () => viewOrderDetail(row),
        },
        {
          icon: () => h(NIcon, {}, () => h(EyeOutline)),
          default: () => '详情',
        },
      )
    },
  },
]

// 处理搜索
const handleSearch = () => {
  pagination.page = 1
}

// 处理筛选变化
const handleFilterChange = () => {
  pagination.page = 1
}

// 查看订单详情
const viewOrderDetail = (order: PayOrderItem) => {
  currentOrder.value = order
  showDetailModal.value = true
}

// 取消订单
const handleCancelOrder = () => {
  if (currentOrder.value && currentOrder.value.status === 'pending') {
    // 模拟API调用
    setTimeout(() => {
      const order = orderList.value.find((item) => item.id === currentOrder.value?.id)
      if (order) {
        order.status = 'cancelled'
        order.updatedAt = new Date().toLocaleString('zh-CN')
      }
      message.success('订单已取消')
      showDetailModal.value = false
    }, 500)
  }
}

// 退款订单
const handleRefundOrder = () => {
  if (currentOrder.value && currentOrder.value.status === 'paid') {
    // 模拟API调用
    setTimeout(() => {
      const order = orderList.value.find((item) => item.id === currentOrder.value?.id)
      if (order) {
        order.status = 'refunded'
        order.updatedAt = new Date().toLocaleString('zh-CN')
      }
      message.success('退款申请已提交')
      showDetailModal.value = false
    }, 500)
  }
}

// 导出订单
const exportOrders = () => {
  message.info('订单导出功能开发中...')
}

// 加载订单列表
const loadOrderList = async () => {
  loading.value = true
  try {
    const res = await payOrderApi.getPayOrderPage({
      page: pagination.page,
      page_size: pagination.pageSize
    })

    if (res.data && res.data.records) {
      orderList.value = res.data.records.map((item: BackendPayOrder) => mapBackendToFrontend(item))
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('订单列表加载失败:', error)
    message.error('订单列表加载失败')
  } finally {
    loading.value = false
  }
}

// 加载今日统计
const loadTodayStats = async () => {
  try {
    const res = await payOrderApi.getTodayStats()
    if (res.data) {
      todayStats.total = res.data.total || 0
      todayStats.amount = res.data.amount || 0
      todayStats.successRate = res.data.success_rate || 0
      todayStats.pending = res.data.pending || 0
    }
  } catch (error) {
    console.error('今日统计加载失败:', error)
  }
}

// 刷新数据
const handleRefresh = () => {
  loadOrderList()
  loadTodayStats()
}

onMounted(() => {
  loadOrderList()
  loadTodayStats()
})
</script>

<style scoped>
.pay-order-container {
  background-color: #f5f7fa;
}

.order-card {
  max-width: 1600px;
  margin: 0 auto;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.search-section {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.action-section {
  display: flex;
  align-items: center;
}

.stats-section {
  margin-bottom: 24px;
  padding: 20px;
  background-color: #fafafa;
  border-radius: 8px;
}

.order-table {
  margin-top: 16px;
}

.order-detail-content {
  padding: 16px 0;
}

@media (max-width: 768px) {
  .pay-order-container {
    padding: 16px;
  }

  .header-section {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }

  .search-section {
    width: 100%;
  }

  .action-section {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
