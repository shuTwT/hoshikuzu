# 易支付

易支付提供统一的支付接口，支持接入多种第三方支付平台。

## 功能特性

- 统一接口：提供统一的支付接口，方便接入
- 多平台支持：支持接入多个易支付平台
- 订单查询：查询易支付订单状态
- 申请退款：支持订单退款

## 配置说明

在 `config.yaml` 中配置易支付参数：

```yaml
payment:
  epay:
    url: your_epay_url
    pid: your_epay_pid
    key: your_epay_key
```

### 配置参数说明

| 参数 | 说明 | 获取方式 |
|------|------|----------|
| url | 易支付接口地址 | 从易支付平台获取 |
| pid | 易支付商户 ID | 在易支付平台注册后获取 |
| key | 易支付密钥 | 在易支付平台设置 |

## API 接口

### 创建易支付订单

```http
POST /api/v1/payment/epay/create
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| order_id | string | 是 | 订单号 |
| amount | decimal | 是 | 支付金额 |
| subject | string | 是 | 支付标题 |
| type | string | 否 | 支付类型 (alipay/wechat/qqpay)，默认 alipay |
| return_url | string | 否 | 同步返回地址 |
| notify_url | string | 否 | 异步通知地址 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": "ORDER_20240101",
    "pay_url": "https://your-epay-url/submit.php?xxx"
  }
}
```

### 查询易支付订单

```http
GET /api/v1/payment/epay/query/:order_id
```

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": "ORDER_20240101",
    "trade_status": "TRADE_SUCCESS",
    "total_amount": "100.00"
  }
}
```

## 支付流程

1. 创建支付订单
2. 获取支付跳转链接
3. 用户跳转到易支付页面完成支付
4. 易支付异步通知支付结果
5. 系统更新订单状态

## 使用示例

### 创建支付订单

```javascript
import { epayApi } from '@/api/payment/epay'

const createPayment = async () => {
  const data = {
    order_id: 'ORDER_20240101',
    amount: 100,
    subject: '商品名称',
    type: 'alipay',
    return_url: 'http://your-domain.com/return',
    notify_url: 'http://your-domain.com/notify'
  }
  
  const result = await epayApi.create(data)
  window.location.href = result.pay_url
}
```

### 查询订单

```javascript
import { epayApi } from '@/api/payment/epay'

const queryOrder = async () => {
  const result = await epayApi.query('ORDER_20240101')
  console.log(result)
}
```

## 回调处理

### 异步通知

易支付成功后，会向 `notify_url` 发送异步通知。系统会自动处理通知并更新订单状态。

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
- [支付宝支付](alipay.md) - 支付宝支付文档
- [微信支付](wechat.md) - 微信支付文档
