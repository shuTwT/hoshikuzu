# 微信支付

微信支付是 支付模块的重要组成部分，支持微信扫码支付、公众号支付、小程序支付等多种支付方式。

## 功能特性

- 扫码支付：生成二维码，用户扫码支付
- 公众号支付：在微信公众号内完成支付
- 小程序支付：在微信小程序内完成支付
- 订单查询：查询微信订单状态
- 申请退款：支持订单退款
- 对账下载：下载对账单

## 配置说明

在 `config.yaml` 中配置微信支付参数：

```yaml
payment:
  wechat:
    app_id: your_wechat_app_id
    mch_id: your_wechat_mch_id
    api_key: your_wechat_api_key
    cert_path: /path/to/apiclient_cert.pem
    key_path: /path/to/apiclient_key.pem
```

### 配置参数说明

| 参数 | 说明 | 获取方式 |
|------|------|----------|
| app_id | 微信应用 ID | 在微信开放平台创建应用后获取 |
| mch_id | 微信商户号 | 在微信支付商户平台获取 |
| api_key | API 密钥 | 在微信支付商户平台设置 |
| cert_path | 证书路径 | 在微信支付商户平台下载 |
| key_path | 密钥路径 | 在微信支付商户平台下载 |

## API 接口

### 创建微信支付订单

```http
POST /api/v1/payment/wechat/create
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| order_id | string | 是 | 订单号 |
| amount | decimal | 是 | 支付金额 |
| subject | string | 是 | 支付标题 |
| body | string | 否 | 支付描述 |
| type | string | 否 | 支付类型 (qr/jsapi/app)，默认 qr |
| openid | string | 否 | 用户 openid（jsapi 支付需要） |
| notify_url | string | 否 | 异步通知地址 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": "ORDER_20240101",
    "qr_code": "data:image/png;base64,...",
    "code_url": "weixin://wxpay/bizpayurl?pr=xxx",
    "prepay_id": "wx2024010100000000000000000000"
  }
}
```

### 查询微信订单

```http
GET /api/v1/payment/wechat/query/:order_id
```

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": "ORDER_20240101",
    "trade_state": "SUCCESS",
    "total_fee": 10000,
    "transaction_id": "4200001234567890"
  }
}
```

### 申请微信退款

```http
POST /api/v1/payment/wechat/refund
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
    "refund_amount": 10000,
    "status": "SUCCESS"
  }
}
```

## 支付流程

### 扫码支付流程

1. 创建支付订单，指定支付类型为 `qr`
2. 获取支付二维码
3. 用户使用微信扫码完成支付
4. 微信异步通知支付结果
5. 系统更新订单状态

### 公众号支付流程

1. 创建支付订单，指定支付类型为 `jsapi`，传入用户 `openid`
2. 获取预支付 ID
3. 调用微信 JSAPI 发起支付
4. 用户在公众号内完成支付
5. 微信异步通知支付结果
6. 系统更新订单状态

### 小程序支付流程

1. 创建支付订单，指定支付类型为 `app`，传入用户 `openid`
2. 获取预支付 ID
3. 调用小程序支付 API 发起支付
4. 用户在小程序内完成支付
5. 微信异步通知支付结果
6. 系统更新订单状态

## 使用示例

### 创建扫码支付

```javascript
import { wechatApi } from '@/api/payment/wechat'

const createQrPayment = async () => {
  const data = {
    order_id: 'ORDER_20240101',
    amount: 100,
    subject: '商品名称',
    type: 'qr'
  }
  
  const result = await wechatApi.create(data)
  console.log(result.qr_code)
}
```

### 创建公众号支付

```javascript
import { wechatApi } from '@/api/payment/wechat'

const createJsapiPayment = async () => {
  const data = {
    order_id: 'ORDER_20240101',
    amount: 100,
    subject: '商品名称',
    type: 'jsapi',
    openid: 'user_openid'
  }
  
  const result = await wechatApi.create(data)
  console.log(result.prepay_id)
}
```

### 查询订单

```javascript
import { wechatApi } from '@/api/payment/wechat'

const queryOrder = async () => {
  const result = await wechatApi.query('ORDER_20240101')
  console.log(result)
}
```

### 申请退款

```javascript
import { wechatApi } from '@/api/payment/wechat'

const refundOrder = async () => {
  const data = {
    order_id: 'ORDER_20240101',
    refund_amount: 100,
    reason: '用户申请退款'
  }
  
  const result = await wechatApi.refund(data)
  console.log(result)
}
```

## 回调处理

### 异步通知

微信支付成功后，会向 `notify_url` 发送异步通知。系统会自动处理通知并更新订单状态。

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
5. **证书安全**：妥善保管支付证书

## 相关文档

- [支付模块](index.md) - 支付模块总览
- [支付宝支付](alipay.md) - 支付宝支付文档
- [易支付](epay.md) - 易支付文档
