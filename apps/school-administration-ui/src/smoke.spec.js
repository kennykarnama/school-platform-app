const { expect } = require('chai')
const fs = require('fs')
const path = require('path')

describe('administration UI test environment', () => {
    it('provides a browser document', () => {
        expect(document).to.exist
        expect(document.createElement('main').tagName).to.equal('MAIN')
    })

    it('serves the application as UTF-8', () => {
        const html = fs.readFileSync(path.join(__dirname, 'index.html'), 'utf8')
        const nginx = fs.readFileSync(path.join(__dirname, '..', 'nginx.conf'), 'utf8')

        expect(html).to.include('<meta charset="UTF-8">')
        expect(nginx).to.match(/\bcharset utf-8;/)
    })
})
