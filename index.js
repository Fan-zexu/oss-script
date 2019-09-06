const OSS = require('ali-oss')
const Conf = require('./config')
const fs = require('fs')
const chalk = require('chalk')

const Origin = new OSS({
  region: Conf.originRegion,
  accessKeyId: Conf.accessKeyId,
  accessKeySecret: Conf.accessKeySecret,
  bucket: Conf.originBucket
})

const Target = new OSS({
  region: Conf.targetRegion,
  accessKeyId: Conf.accessKeyId,
  accessKeySecret: Conf.accessKeySecret,
  bucket: Conf.targetBucket
})

let marker = ''

async function go () {
  const result = await Origin.list({
    marker,
    'max-keys': 100
  })

  result.objects.forEach(async (item) => {
    if (item.size) {
      try {
        const res = await Target.head(item.name)
        if (res.res.status == 200) {
          const r = await Origin.head(item.name)
          if (res.res.headers['content-md5'] !== r.res.headers['content-md5']) {
            console.log(chalk.green('同名文件的content-md5：', res.res.headers['content-md5'], r.res.headers['content-md5']))
            console.log(chalk.yellow(`同名但内容不同的文件：${item.url}`))
            console.log()
            fs.appendFile('./data.txt', `源地址：${item.url} \n目标地址：${Conf.targetEndpoint}/${item.name} \n\n`, 'utf-8', (err) => {
              if (err) console.log(err)
            })
          } else {
            console.log(`同名且内容相同的文件：${item.url}`)
            console.log()
          }
        }
      } catch (err) {
        if (err.status != 404) {
          fs.appendFile('./data.txt', `当前执行位置nextMarker：${result.nextMarker} \n`, 'utf-8', (err) => {
            if (err) console.log(err)
          })
        }
      }
    }
  })

  if (result.isTruncated) {
    marker = result.nextMarker
    go()
  }
}
console.log('开始对比....')
console.log()
go()
