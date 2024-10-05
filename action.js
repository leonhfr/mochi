const childProcess = require('child_process')
const os = require('os')
const process = require('process')

const distDir = 'dist'
const architectures = {
  x64: 'amd64',
  arm64: 'arm64',
}

function chooseDirectory(artifacts) {
    const platform = os.platform()
    const arch = os.arch()

    for (let i = 0; i < artifacts.length; i++) {
      const { goos, goarch, path, type } = artifacts[i];
      if (type === 'Binary' && goos === platform && goarch === architectures[arch]) {
        return path
      }
    }

    console.error(`Unsupported platform (${platform}) and architecture (${arch})`)
    process.exit(1)
}

function main() {
    const artifacts = require(`${__dirname}/${distDir}/artifacts.json`)
    const binary = chooseDirectory(artifacts)
    const mainScript = `${__dirname}/${binary}`
    const spawnSyncReturns = childProcess.spawnSync(mainScript, { stdio: 'inherit' })
    const status = spawnSyncReturns.status
    if (typeof status === 'number') {
        process.exit(status)
    }
    process.exit(1)
}

if (require.main === module) {
    main()
}
