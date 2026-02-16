// 코드 생성기 페이지 로직
let currentLanguage = 'go';
let selectedModules = [];
let allModules = [];

// 페이지 로드 시 초기화
document.addEventListener('DOMContentLoaded', async () => {
    try {
        allModules = await window.go.main.App.GetModules();
        renderModules();
    } catch (err) {
        console.error('모듈 로드 실패:', err);
    }
});

function setLanguage(lang) {
    currentLanguage = lang;

    const btnGo = document.getElementById('btn-go');
    const btnNode = document.getElementById('btn-node');
    const btnGenerate = document.getElementById('btn-generate');

    if (lang === 'go') {
        btnGo.className = 'join-item btn btn-sm flex-1 btn-warning';
        btnNode.className = 'join-item btn btn-sm flex-1 btn-ghost';
        btnGenerate.className = 'btn btn-warning w-full mt-4';
        btnGenerate.textContent = 'Go 프로젝트 생성';
    } else {
        btnGo.className = 'join-item btn btn-sm flex-1 btn-ghost';
        btnNode.className = 'join-item btn btn-sm flex-1 btn-success';
        btnGenerate.className = 'btn btn-success w-full mt-4';
        btnGenerate.textContent = 'Node.js 프로젝트 생성';
    }
}

function renderModules() {
    const list = document.getElementById('module-list');
    list.innerHTML = '';

    for (const mod of allModules) {
        const isSelected = selectedModules.includes(mod.id);
        const el = document.createElement('div');
        el.className = `card card-compact bg-base-300 border cursor-pointer transition mb-2 ${
            isSelected
            ? 'border-primary bg-primary/10'
            : 'border-base-content/10 hover:border-base-content/30'
        }`;
        el.onclick = () => toggleModule(mod.id);
        el.innerHTML = `
            <div class="card-body p-3">
                <div class="flex justify-between items-start">
                    <h3 class="font-bold text-sm">${mod.name}</h3>
                    ${isSelected ? '<span class="badge badge-primary badge-sm">&#10004;</span>' : ''}
                </div>
                <p class="text-xs text-base-content/50">${mod.description}</p>
            </div>
        `;
        list.appendChild(el);
    }
}

function toggleModule(id) {
    if (selectedModules.includes(id)) {
        selectedModules = selectedModules.filter(m => m !== id);
    } else {
        selectedModules.push(id);
    }
    renderModules();
}

async function selectDirectory() {
    try {
        const dir = await window.go.main.App.SelectDirectory();
        if (dir) {
            const projectName = document.getElementById('projectName').value;
            document.getElementById('targetPath').value = dir + '\\' + projectName;
        }
    } catch (err) {
        console.error('폴더 선택 실패:', err);
    }
}

async function generateProject() {
    const logOutput = document.getElementById('log-output');
    const logText = document.getElementById('log-text');
    logOutput.classList.remove('hidden');
    logText.textContent = `[${currentLanguage.toUpperCase()}] 생성 중... 잠시만 기다려주세요`;

    const config = {
        projectName: document.getElementById('projectName').value,
        targetPath: document.getElementById('targetPath').value,
        dbServer: document.getElementById('dbServer').value,
        dbUser: document.getElementById('dbUser').value,
        dbPw: document.getElementById('dbPw').value,
        dbName: document.getElementById('dbName').value,
        modules: selectedModules
    };

    try {
        const result = await window.go.main.App.GenerateProject(config, currentLanguage);

        if (result.success) {
            if (currentLanguage === 'go') {
                logText.textContent = 'Go 서버 생성 완료!\n터미널에서 실행:\n  go mod tidy\n  go run .';
            } else {
                logText.textContent = 'Node.js 서버 생성 완료!\n터미널에서 실행:\n  npm install\n  npm start';
            }
        } else {
            logText.textContent = '실패: ' + result.message;
        }
    } catch (err) {
        logText.textContent = '오류: ' + err;
    }
}
