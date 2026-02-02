# 支付宝支付

支付宝支付支持支付宝扫码支付、网页支付、移动支付等多种支付方式。

## 功能特性

- 扫码支付：生成二维码，用户扫码支付
- 网页支付：跳转到支付宝网页完成支付
- 移动支付：支持移动端支付宝支付
- 订单查询：查询支付宝订单状态
- 申请退款：支持订单退款
- 对账下载：下载对账单

## 配置说明

在 `config.yaml` 中配置支付宝参数：

```yaml
payment:
  alipay:
    app_id: your_alipay_app_id
    private_key: your_alipay_private_key
    public_key: your_alipay_public_key
```

### 配置参数说明

| 参数 | 说明 | 获取方式 |
|------|------|----------|
| app_id | 支付宝应用 ID | 在支付宝开放平台创建应用后获取 |
| private_key | 应用私钥 | 使用支付宝提供的工具生成 |
| public_key | 支付宝公钥 | 在支付宝开放平台获取 |

## API 接口

### 创建支付宝支付订单

```http
POST /api/v1/payment/alipay/create
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| order_id | string | 是 | 订单号 |
| amount | decimal | 是 | 支付金额 |
| subject | string | 是 | 支付标题 |
| body | string | 否 | 支付描述 |
| type | string | 否 | 支付类型 (qr/page/app)，默认 qr |
| return_url | string | 否 | 同步返回地址 |
| notify_url | string | 否 | 异步通知地址 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": "ORDER_20240101",
    "qr_code": "data:image/png;base64,...",
    "pay_url": "https://qr.alipay.com/xxx"
  }
}
```

### 查询支付宝订单

```http
GET /api/v1/payment/alipay/query/:order_id
```

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": "ORDER_20240101",
    "trade_status": "TRADE_SUCCESS",
    "total_amount": "100.00",
    "receipt_amount": "100.00"
  }
}
```

### 申请支付宝退款

```http
POST /api/v1/payment/alipay/refund
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| order_id | string | 是 | 订单号 |
| refund_amount | decimal | 是 | 退款金额 |
| reason | string | 否 | 退款原因 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": "ORDER_20240101",
    "refund_id": "REFUND_20240101",
    "refund_amount": "100.00",
    "status": "SUCCESS"
  }
}
```

## 支付流程

### 扫码支付流程

1. 创建支付订单，指定支付类型为 `qr`
2. 获取支付二维码
3. 用户使用支付宝扫码完成支付
4. 支付宝异步通知支付结果
5. 系统更新订单状态

### 网页支付流程

1. 创建支付订单，指定支付类型为 `page`
2. 获取支付跳转链接
3. 用户跳转到支付宝网页完成支付
4. 支付宝异步通知支付结果
5. 系统更新订单状态

## 使用示例

### 创建扫码支付

```javascript
import { alipayApi } from '@/api/payment/alipay'

const createQrPayment = async () => {
  const data = {
    order_id: 'ORDER_20240101',
    amount: 100,
    subject: '商品名称',
    type: 'qr'
  }
  
  const result = await alipayApi.create(data)
  console.log(result.qr_code)
}
```

### 创建网页支付

```javascript
import { alipayApi } from '@/api/payment/alipay'

const createPagePayment = async () => {
  const data = {
    order_id: 'ORDER_20240101',
    amount: 100,
    subject: '商品名称',
    type: 'page',
    return_url: 'http://your-domain.com/return',
    notify_url: 'http://your-domain.com/notify'
  }
  
  const result = await alipayApi.create(data)
  window.location.href = result.pay_url
}
```

### 查询订单

```javascript
import { alipayApi } from '@/api/payment/alipay'

const queryOrder = async () => {
  const result = await alipayApi.query('ORDER_20240101')
  console.log(result)
}
```

### 申请退款

```javascript
import { alipayApi } from '@/api/payment/alipay'

const refundOrder = async () => {
  const data = {
    order_id: 'ORDER_20240101',
    refund_amount: 100,
    reason: '用户申请退款'
  }
  
  const result = await alipayApi.refund(data)
  console.log(result)
}
```

## 回调处理

### 异步通知

支付宝支付成功后，会向 `notify_url` 发送异步通知。系统会自动处理通知并更新订单状态。

### 同步跳转

支付完成后，用户会跳转到 `return_url`，可以通过查询订单状态确认支付结果。

## 错误处理

### 常见错误

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| INVALID_PARAMETER | 参数错误 | 检查请求参数是否正确 |
| ORDER_NOT_FOUND | 订单不存在 | 检查订单号是否正确 |
| SIGN_ERROR | 签名错误 | 检查签名配置是否正确 |
| SYSTEM_ERROR | 系统错误 | 联系技术支持 |

## 安全建议

1. **签名验证**：所有回调必须验证签名
2. **金额校验**：处理回调时校验订单金额
3. **重复通知**：处理重复的异步通知
4. **日志记录**：记录所有支付相关日志

## 相关文档

- [支付模块](index.md) - 支付模块总览
- [微信支付](wechat.md) - 微信支付文档
- [易支付](epay.md) - 易支付文档
