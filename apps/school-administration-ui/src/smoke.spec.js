const { expect } = require('chai')

describe('administration UI test environment', () => {
    it('provides a browser document', () => {
        expect(document).to.exist
        expect(document.createElement('main').tagName).to.equal('MAIN')
    })
})

