# BingImgAPI
Bing(必应)美图API

实现其实很简单，Bing本身有图片API只是由于跨域不好使用；尝试反向代理效果不好；最终打算写个小服务转发一下json数据。

请求地址：[bing.uosblog.top](https://bing.uosblog.top)

请求类型：POST/GET

参数：无

将返回14天前～现在的所有图片，共15条

|  返回值   |     释义     |
| :-------: | :----------: |
|    url    |   图片地址   |
| copyright | 图片版权信息 |
|  enddate  |     日期     |

```json

  {
    "images" : [
    {
      "urlbase" : "\/\/cn.bing.com\/th?id=OHR.MehrangarhCourtyard_ZH-CN3216739355",
      "enddate" : "20220126",
      "startdate" : "20220125",
      "quiz" : "\/\/cn.bing.com\/search?q=Bing+homepage+quiz&filters=WQOskey:%22HPQuiz_20220125_MehrangarhCourtyard%22&FORM=HPQUIZ",
      "hsh" : "",
      "fullstartdate" : "202201251600",
      "copyrightlink" : "https:\/\/www.bing.com\/search?q=%E6%A2%85%E5%85%B0%E5%8A%A0%E5%B0%94%E5%A0%A1&form=hpcapt&mkt=zh-cn",
      "url" : "\/\/cn.bing.com\/th?id=OHR.MehrangarhCourtyard_ZH-CN3216739355_1920x1080.jpg&rf=LaDigue_1920x1080.jpg&pid=hp",
      "hs" : [
        ""
      ],
      "title" : "",
      "drk" : 1,
      "wp" : true,
      "bot" : 1,
      "copyright" : "梅兰加尔堡，印度焦特布尔 (© Jayakumar\/Shutterstock)",
      "top" : 1
    },
    "tooltips" : {
    "previous" : "上一个图像",
    "loading" : "正在加载...",
    "next" : "下一个图像",
    "walle" : "此图片不能下载用作壁纸。",
    "walls" : "下载今日美图。仅限用作桌面壁纸。"
  }
}
```

