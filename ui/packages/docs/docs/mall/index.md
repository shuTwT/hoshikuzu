# 商城模块

商城模块提供完整的电商功能，包括商品管理、会员系统、优惠券等功能。

## 功能特性

- **商品管理**：商品的创建、编辑、上下架、库存管理
- **会员系统**：会员等级、积分、权益管理
- **优惠券**：优惠券创建、发放、核销
- **订单管理**：订单创建、查询、退款
- **支付集成**：集成多种支付方式

## 快速开始

### 1. 创建商品

```javascript
import { productApi } from '@/api/mall/product'

const createProduct = async () => {
  const data = {
    name: '商品名称',
    price: 100,
    stock: 100,
    description: '商品描述'
  }
  
  const result = await productApi.create(data)
  console.log(result)
}
```

### 2. 创建会员等级

```javascript
import { memberApi } from '@/api/mall/member'

const createMemberLevel = async () => {
  const data = {
    name: '黄金会员',
    discount: 0.9,
    points: 1000
  }
  
  const result = await memberApi.createLevel(data)
  console.log(result)
}
```

### 3. 创建优惠券

```javascript
import { couponApi } from '@/api/mall/coupon'

const createCoupon = async () => {
  const data = {
    name: '新用户优惠券',
    type: 'discount',
    value: 10,
    min_amount: 100
  }
  
  const result = await couponApi.create(data)
  console.log(result)
}
```

## 模块概览

### 商品管理

- 商品创建与编辑
- 商品分类管理
- 商品上下架
- 库存管理
- 商品搜索

### 会员系统

- 会员等级管理
- 会员权益管理
- 积分系统
- 会员统计

### 优惠券

- 优惠券创建
- 优惠券发放
- 优惠券核销
- 优惠券统计

## API 接口

### 商品接口

- `POST /api/v1/mall/product/create` - 创建商品
- `PUT /api/v1/mall/product/update/:id` - 更新商品
- `DELETE /api/v1/mall/product/delete/:id` - 删除商品
- `GET /api/v1/mall/product/page` - 获取商品列表
- `GET /api/v1/mall/product/query/:id` - 查询商品详情

### 会员接口

- `POST /api/v1/mall/member/level/create` - 创建会员等级
- `PUT /api/v1/mall/member/level/update/:id` - 更新会员等级
- `DELETE /api/v1/mall/member/level/delete/:id` - 删除会员等级
- `GET /api/v1/mall/member/level/list` - 获取会员等级列表
- `GET /api/v1/mall/member/info` - 获取会员信息

### 优惠券接口

- `POST /api/v1/mall/coupon/create` - 创建优惠券
- `PUT /api/v1/mall/coupon/update/:id` - 更新优惠券
- `DELETE /api/v1/mall/coupon/delete/:id` - 删除优惠券
- `GET /api/v1/mall/coupon/list` - 获取优惠券列表
- `POST /api/v1/mall/coupon/receive` - 领取优惠券

## 相关文档

- [商品](product.md) - 商品管理详细文档
- [会员](member.md) - 会员系统详细文档
- [优惠券](coupon.md) - 优惠券详细文档
- [支付模块](../payment/index.md) - 支付模块文档
