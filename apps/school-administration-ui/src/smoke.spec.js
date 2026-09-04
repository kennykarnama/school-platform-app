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

    it('uses an in-application confirmation dialog for applying setup data', () => {
        const setupPage = fs.readFileSync(
            path.join(__dirname, 'components/includes/pages/setup/setup.riot'),
            'utf8',
        )

        expect(setupPage).not.to.match(/window\.(alert|confirm|prompt)\s*\(/)
        expect(setupPage).to.include('Konfirmasi Terapkan Data')
        expect(setupPage).to.include('onclick={confirmApply}')
        expect(setupPage).to.include("event.key === 'Escape'")
    })

    it('uses same-origin API routing in production and an explicit local development URL', () => {
        const appMain = fs.readFileSync(
            path.join(__dirname, 'components/global/app-main/app-main.riot'),
            'utf8',
        )
        const packageJson = JSON.parse(
            fs.readFileSync(path.join(__dirname, '..', 'package.json'), 'utf8'),
        )

        expect(appMain).to.include("process.env.SCHOOL_ADMINISTRATION_API_BASE_URL || ''")
        expect(appMain).not.to.include('window.location.hostname}:8081')
        expect(packageJson.scripts['start-dev']).to.include(
            'SCHOOL_ADMINISTRATION_API_BASE_URL=http://127.0.0.1:8081',
        )
    })

    it('shows the authenticated teacher profile and uses the logout endpoint', () => {
        const appMain = fs.readFileSync(
            path.join(__dirname, 'components/global/app-main/app-main.riot'),
            'utf8',
        )
        const asideMenu = fs.readFileSync(
            path.join(__dirname, 'components/includes/aside-menu/aside-menu.riot'),
            'utf8',
        )

        expect(appMain).to.include('<aside-menu teacher={state.teacher}>')
        expect(asideMenu).to.include('class="teacher-avatar"')
        expect(asideMenu).to.include("axios.post('/api/v1/teacher/logout')")
        expect(asideMenu).to.include("router.push('/#login')")
    })
})
