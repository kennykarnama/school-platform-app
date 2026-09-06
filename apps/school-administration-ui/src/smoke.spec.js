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
        const login = fs.readFileSync(
            path.join(__dirname, 'components/includes/pages/login/login.riot'),
            'utf8',
        )

        expect(appMain).to.include('<aside-menu teacher={state.teacher}>')
        expect(appMain).to.include("window.location.hash !== '#login'")
        expect(appMain).to.include("window.addEventListener('school:teacher-authenticated'")
        expect(login).to.include("window.dispatchEvent(new CustomEvent('school:teacher-authenticated'))")
        expect(asideMenu).to.include('class="teacher-avatar"')
        expect(asideMenu).to.include("axios.post('/api/v1/teacher/logout')")
        expect(asideMenu).to.include("router.push('/#login')")
    })

    it('exposes role-specific tenant administration pages', () => {
        const appContent = fs.readFileSync(
            path.join(__dirname, 'components/global/app-content/app-content.riot'),
            'utf8',
        )
        const asideMenu = fs.readFileSync(
            path.join(__dirname, 'components/includes/aside-menu/aside-menu.riot'),
            'utf8',
        )
        const administration = fs.readFileSync(
            path.join(__dirname, 'components/includes/pages/administration/administration.riot'),
            'utf8',
        )
        const schools = fs.readFileSync(
            path.join(__dirname, 'components/includes/pages/platform/platform-schools.riot'),
            'utf8',
        )

        expect(appContent).to.include('path="/#administration"')
        expect(appContent).to.include('path="/#platform/schools"')
        expect(asideMenu).to.include('"roles": ["school_admin"]')
        expect(asideMenu).to.include('"roles": ["platform_admin"]')
        expect(administration).to.include("axios.put(`/api/v1/admin/teachers/${this.state.selectedTeacher.id}/access`")
        expect(schools).to.include("axios.post('/api/v1/platform/schools'")
    })

    it('exposes student management to administrators and teachers with role-specific actions', () => {
        const appContent = fs.readFileSync(
            path.join(__dirname, 'components/global/app-content/app-content.riot'),
            'utf8',
        )
        const asideMenu = fs.readFileSync(
            path.join(__dirname, 'components/includes/aside-menu/aside-menu.riot'),
            'utf8',
        )
        const students = fs.readFileSync(
            path.join(__dirname, 'components/includes/pages/students/students.riot'),
            'utf8',
        )
        const attendanceCreate = fs.readFileSync(
            path.join(__dirname, 'components/includes/pages/absensi/absensi-input-form/absensi-input-form.riot'),
            'utf8',
        )

        expect(appContent).to.include('path="/#students"')
        expect(asideMenu).to.include('"label": "Siswa"')
        expect(asideMenu).to.include('"roles": ["school_admin", "teacher"]')
        expect(students).to.include("isAdmin() { return this.state.role === 'school_admin' }")
        expect(students).to.include("axios.get('/api/v1/students'")
        expect(students).to.include("axios.post('/api/v1/students'")
        expect(students).to.include('axios.patch(`/api/v1/admin/students/${item.id}`')
        expect(students).to.include('axios.patch(`/api/v1/admin/students/${d.studentID}/status`')
        expect(students).to.include('axios.patch(`/api/v1/student/class/${d.studentClassID}/deactivate`')
        expect(students).to.include('axios.patch(`/api/v1/student/class/${d.studentClassID}/restore`')
        expect(students).to.include("axios.post('/api/v1/student/transfer'")
        expect(students).to.include('role="dialog"')
        expect(students).to.include("aria-modal=\"true\"")
        expect(attendanceCreate).to.include("url: '/api/v1/students'")
    })

    it('supports per-student transfer from the attendance table', () => {
        const absensiTable = fs.readFileSync(
            path.join(__dirname, 'components/includes/pages/absensi/absensi-table/absensi-table.riot'),
            'utf8',
        )
        expect(absensiTable).to.include("axios.post('/api/v1/student/transfer'")
        expect(absensiTable).to.include('role="dialog"')
    })

    it('sends class labels to class-scoped attendance APIs', () => {
        const attendancePages = [
            'components/includes/pages/absensi/absensi-input/absensi-input.riot',
            'components/includes/pages/absensi/absensi-input-form/absensi-input-form.riot',
            'components/includes/pages/absensi/absensi-transfer-form/absensi-transfer-form.riot',
            'components/includes/pages/absensi/rekap/rekap.riot',
            'components/includes/pages/absensi/rekap/rekap-klasikal.riot',
        ].map(file => fs.readFileSync(path.join(__dirname, file), 'utf8'))

        attendancePages.forEach(page => {
            expect(page).not.to.include('value={classItem.ID}')
            expect(page).to.include('value={classItem.label}')
        })
    })

    it('configures consistent button spacing and flex layout for dialog boxes', () => {
        const responsiveCss = fs.readFileSync(
            path.join(__dirname, 'responsive.css'),
            'utf8',
        )
        const absensiTable = fs.readFileSync(
            path.join(__dirname, 'components/includes/pages/absensi/absensi-table/absensi-table.riot'),
            'utf8',
        )

        expect(responsiveCss).to.match(/\.modal-container \.modal-footer\s*\{[^}]*display:\s*flex;/)
        expect(responsiveCss).to.match(/\.modal-container \.modal-footer\s*\{[^}]*justify-content:\s*flex-end;/)
        expect(responsiveCss).to.match(/\.modal-container \.modal-footer\s*\{[^}]*gap:\s*\.5rem;/)
        expect(absensiTable).not.to.include('style="margin-right: 2%;"')
    })
})

