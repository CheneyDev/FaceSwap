![image-20240416110242995](./assets/README/image-20240416110242995.png)

# FaceSwap

FaceSwap是一个能够将任何人脸图片转换到另一张图片中的工具

## 使用方法

1. 准备两张图片：一张包含要转换的人脸图片，另一张作为目标的背景图片 （注意顺序，第一张图检测不到人脸程序会出错）
2. 进入 [FaceSwap](https://swap.qqdd.dev/) 网页，上传准备好的两张图片
3. 点击 "Submit" 按钮，等待片刻
4. 生成的换脸图片将会显示在界面上

## 示例

| 原始图片                                                | 目标图片                                                     | 结果                                          |
| ------------------------------------------------------- | ------------------------------------------------------------ | --------------------------------------------- |
| ![image](./assets/README/MTk4MTczMTkzNzI1Mjg5NjYy.webp) | ![image_to_become](./assets/README/cHJpdmF0ZS9sci9pbWFnZXMvd2Vic2l0ZS8yMDIyLTA1L2pvYjU4NS12MjE2LXRhbmctYXVtLTAxMC1leWUtYXJ0cHJpbnRzLmpwZw.webp) | ![output](./assets/README/ComfyUI_00001_.png) |

## 技术栈和服务

| 技术/工具                                                   | 部署/服务                                                    |
| ----------------------------------------------------------- | ------------------------------------------------------------ |
| 后端：Go Gin Gorm                                           | [Render](https://render.com/)                                |
| 前端：Next.js                                               | [Vercel](https://vercel.com/)                                |
| 样式组件：Tailwind CSS, [shadcn/ui](https://ui.shadcn.com/) |                                                              |
| 数据库：Postgres                                            | [Supabase](https://supabase.com/)                            |
| 对象存储：CloudFlare R2                                     | [CloudFlare](https://www.cloudflare.com/developer-platform/r2/) |
| 换脸 API：fofr/become-image                                 | [Replicate](https://replicate.com/fofr/become-image)         |



