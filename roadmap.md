 - [ x ] 会员等级
 - [ x ] 会员
 - [ x ] 优惠券
 - [ x ] 优惠券使用记录
 - [ x ] 说说
 - [ x ] 钱包
 - [ x ] 访问日志
 - [ x ] 商品管理
 - [  ] 路由区分,/console/**为后台路由,/api/**为 api路由('/api/v1','/api/public'),其余为前台路由
 - [ ] ssr支持
 - [ ] 多主题
        ``` 
        hoshikuzu-theme-xxx
        |---templates
        |   |---assets/ 静态资源
        |       |---css/
        |       |---js/
        |       |---img/
        |   |---index.html 首页
        |   |---post.html 文章页面
        |   |---page.html 自定义页面
        |   |---categories.html 分类列表页面
        |   |---category.html 分类页面
        |   |---tags.html 标签列表页面
        |   |---tag.html 标签页面
        |   |---archives.html 归档页面
        |   |---author.html 作者页面
        |   |---404.html 404页面
        |---theme.yaml 主题声明，后台在主题详情处可查看
        |---config.yaml 主题配置schema声明，将在前端渲染为表单
        ```
        