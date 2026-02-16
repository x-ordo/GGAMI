// 코드 생성기 페이지 로직
let currentLanguage = 'go';
let selectedModules = [];
let allModules = [];

// GORM 모드 상태
let gormMode = false;
let models = [];
let activeModelIndex = -1;

const FIELD_TYPES = [
    { value: 'string', label: 'String' },
    { value: 'int', label: 'Int' },
    { value: 'uint', label: 'Uint' },
    { value: 'float64', label: 'Float64' },
    { value: 'bool', label: 'Bool' },
    { value: 'time.Time', label: 'DateTime' },
];

const GORM_TAGS = ['primaryKey', 'unique', 'not null', 'index'];

// 페이지 로드 시 초기화
document.addEventListener('DOMContentLoaded', async () => {
    try {
        allModules = await window.go.main.App.GetModules();
        renderModules();
    } catch (err) {
        console.error('모듈 로드 실패:', err);
    }
});

// ============ 모드 전환 ============

function setMode(mode) {
    gormMode = (mode === 'gorm');

    const btnSimple = document.getElementById('btn-simple');
    const btnGorm = document.getElementById('btn-gorm');
    const langSelector = document.getElementById('lang-selector');
    const dbTypeSelector = document.getElementById('db-type-selector');
    const modulePanel = document.getElementById('module-panel');
    const modelPanel = document.getElementById('model-panel');
    const rbacSection = document.getElementById('rbac-section');
    const btnGenerate = document.getElementById('btn-generate');

    if (gormMode) {
        btnSimple.className = 'join-item btn btn-sm flex-1 btn-ghost';
        btnGorm.className = 'join-item btn btn-sm flex-1 btn-primary';
        langSelector.classList.add('hidden');
        dbTypeSelector.classList.remove('hidden');
        modulePanel.classList.add('hidden');
        modelPanel.classList.remove('hidden');
        rbacSection.classList.remove('hidden');
        btnGenerate.textContent = 'GORM Full-Stack 프로젝트 생성';
        btnGenerate.className = 'btn btn-primary w-full mt-4';
        onDBTypeChange();
    } else {
        btnSimple.className = 'join-item btn btn-sm flex-1 btn-warning';
        btnGorm.className = 'join-item btn btn-sm flex-1 btn-ghost';
        langSelector.classList.remove('hidden');
        dbTypeSelector.classList.add('hidden');
        modulePanel.classList.remove('hidden');
        modelPanel.classList.add('hidden');
        rbacSection.classList.add('hidden');
        document.getElementById('db-connection-fields').classList.remove('hidden');
        setLanguage(currentLanguage);
    }
}

function onDBTypeChange() {
    const dbType = document.getElementById('dbType').value;
    const connFields = document.getElementById('db-connection-fields');
    if (dbType === 'sqlite') {
        connFields.classList.add('hidden');
    } else {
        connFields.classList.remove('hidden');
    }
}

function onRBACToggle() {
    const enabled = document.getElementById('rbacEnabled').checked;
    const details = document.getElementById('rbac-details');
    if (enabled) {
        details.classList.remove('hidden');
        renderRBACMatrix();
    } else {
        details.classList.add('hidden');
    }
}

// ============ 언어 선택 (Simple 모드) ============

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

// ============ 모듈 (Simple 모드) ============

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

// ============ 모델 정의 (GORM 모드) ============

function addModel() {
    const name = prompt('모델 이름 (PascalCase):', 'Product');
    if (!name) return;

    models.push({
        name: name,
        fields: [
            { name: 'ID', type: 'uint', gormTags: ['primaryKey'], defaultVal: '', jsonName: 'id' }
        ]
    });
    activeModelIndex = models.length - 1;
    renderModelTabs();
    renderModelEditor();
    renderRBACMatrix();
}

function removeModel(index) {
    if (!confirm(`"${models[index].name}" 모델을 삭제하시겠습니까?`)) return;
    models.splice(index, 1);
    if (activeModelIndex >= models.length) {
        activeModelIndex = models.length - 1;
    }
    renderModelTabs();
    renderModelEditor();
    renderRBACMatrix();
}

function selectModel(index) {
    activeModelIndex = index;
    renderModelTabs();
    renderModelEditor();
}

function renderModelTabs() {
    const tabs = document.getElementById('model-tabs');
    tabs.innerHTML = '';
    models.forEach((model, i) => {
        const tab = document.createElement('a');
        tab.className = `tab ${i === activeModelIndex ? 'tab-active' : ''}`;
        tab.textContent = model.name;
        tab.onclick = () => selectModel(i);
        tabs.appendChild(tab);
    });
}

function renderModelEditor() {
    const editor = document.getElementById('model-editor');
    if (activeModelIndex < 0 || activeModelIndex >= models.length) {
        editor.innerHTML = '<div class="text-base-content/50 text-sm text-center py-8">"+ 모델 추가" 버튼으로 첫 모델을 만드세요</div>';
        return;
    }

    const model = models[activeModelIndex];
    let html = `
        <div class="flex justify-between items-center mb-3">
            <div class="form-control">
                <input type="text" value="${model.name}" onchange="updateModelName(${activeModelIndex}, this.value)"
                    class="input input-bordered input-sm font-bold" />
            </div>
            <button onclick="removeModel(${activeModelIndex})" class="btn btn-error btn-xs">삭제</button>
        </div>
        <table class="table table-xs">
            <thead>
                <tr>
                    <th>필드명</th>
                    <th>타입</th>
                    <th>태그</th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
    `;

    model.fields.forEach((field, fi) => {
        html += `<tr>
            <td><input type="text" value="${field.name}" onchange="updateField(${activeModelIndex}, ${fi}, 'name', this.value)" class="input input-bordered input-xs w-24" /></td>
            <td>
                <select onchange="updateField(${activeModelIndex}, ${fi}, 'type', this.value)" class="select select-bordered select-xs">
                    ${FIELD_TYPES.map(t => `<option value="${t.value}" ${t.value === field.type ? 'selected' : ''}>${t.label}</option>`).join('')}
                </select>
            </td>
            <td>
                <div class="flex flex-wrap gap-1">
                    ${GORM_TAGS.map(tag => `
                        <label class="label cursor-pointer gap-1 p-0">
                            <input type="checkbox" class="checkbox checkbox-xs"
                                ${field.gormTags.includes(tag) ? 'checked' : ''}
                                onchange="toggleFieldTag(${activeModelIndex}, ${fi}, '${tag}')" />
                            <span class="text-xs">${tag}</span>
                        </label>
                    `).join('')}
                </div>
            </td>
            <td>
                <button onclick="removeField(${activeModelIndex}, ${fi})" class="btn btn-ghost btn-xs text-error" ${fi === 0 ? 'disabled' : ''}>X</button>
            </td>
        </tr>`;
    });

    html += `</tbody></table>
        <button onclick="addField(${activeModelIndex})" class="btn btn-ghost btn-xs mt-2">+ 필드 추가</button>
    `;
    editor.innerHTML = html;
}

function updateModelName(mi, name) {
    models[mi].name = name;
    renderModelTabs();
    renderRBACMatrix();
}

function addField(mi) {
    models[mi].fields.push({
        name: '',
        type: 'string',
        gormTags: [],
        defaultVal: '',
        jsonName: ''
    });
    renderModelEditor();
}

function removeField(mi, fi) {
    models[mi].fields.splice(fi, 1);
    renderModelEditor();
}

function updateField(mi, fi, key, value) {
    models[mi].fields[fi][key] = value;
    if (key === 'name') {
        models[mi].fields[fi].jsonName = toSnakeCase(value);
    }
}

function toggleFieldTag(mi, fi, tag) {
    const tags = models[mi].fields[fi].gormTags;
    const idx = tags.indexOf(tag);
    if (idx >= 0) {
        tags.splice(idx, 1);
    } else {
        tags.push(tag);
    }
}

function toSnakeCase(str) {
    return str.replace(/([A-Z])/g, (match, p1, offset) => {
        return (offset > 0 ? '_' : '') + p1.toLowerCase();
    });
}

// ============ RBAC 매트릭스 ============

function renderRBACMatrix() {
    const container = document.getElementById('rbac-matrix');
    if (!container) return;
    if (models.length === 0) {
        container.innerHTML = '<p class="text-xs text-base-content/50">모델을 먼저 추가하세요</p>';
        return;
    }

    const rolesStr = document.getElementById('rbacRoles').value;
    const roles = rolesStr.split(',').map(r => r.trim()).filter(Boolean);

    let html = '<div class="overflow-x-auto"><table class="table table-xs"><thead><tr><th>모델</th>';
    roles.forEach(role => {
        html += `<th colspan="4" class="text-center">${role}</th>`;
    });
    html += '</tr><tr><th></th>';
    roles.forEach(() => {
        html += '<th class="text-xs">C</th><th class="text-xs">R</th><th class="text-xs">U</th><th class="text-xs">D</th>';
    });
    html += '</tr></thead><tbody>';

    models.forEach((model, mi) => {
        html += `<tr><td class="font-bold">${model.name}</td>`;
        roles.forEach((role, ri) => {
            ['create', 'read', 'update', 'delete'].forEach(action => {
                const id = `rbac-${mi}-${ri}-${action}`;
                const isAdmin = role === 'admin';
                html += `<td><input type="checkbox" id="${id}" class="checkbox checkbox-xs" ${isAdmin ? 'checked' : (action === 'read' ? 'checked' : '')} /></td>`;
            });
        });
        html += '</tr>';
    });

    html += '</tbody></table></div>';
    container.innerHTML = html;
}

// ============ 공통 ============

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

    const config = {
        projectName: document.getElementById('projectName').value,
        targetPath: document.getElementById('targetPath').value,
        dbServer: document.getElementById('dbServer').value,
        dbUser: document.getElementById('dbUser').value,
        dbPw: document.getElementById('dbPw').value,
        dbName: document.getElementById('dbName').value,
        modules: selectedModules,
        gormMode: gormMode,
    };

    if (gormMode) {
        config.dbType = document.getElementById('dbType').value;
        config.models = models.map(m => ({
            name: m.name,
            fields: m.fields.map(f => ({
                name: f.name,
                type: f.type,
                gormTags: f.gormTags,
                defaultVal: f.defaultVal || '',
                jsonName: f.jsonName || toSnakeCase(f.name),
            }))
        }));

        // RBAC config
        if (document.getElementById('rbacEnabled').checked) {
            const rolesStr = document.getElementById('rbacRoles').value;
            const roles = rolesStr.split(',').map(r => r.trim()).filter(Boolean);

            const modelPerms = models.map((model, mi) => ({
                modelName: model.name,
                permissions: roles.map((role, ri) => ({
                    role: role,
                    create: document.getElementById(`rbac-${mi}-${ri}-create`)?.checked || false,
                    read: document.getElementById(`rbac-${mi}-${ri}-read`)?.checked || false,
                    update: document.getElementById(`rbac-${mi}-${ri}-update`)?.checked || false,
                    delete: document.getElementById(`rbac-${mi}-${ri}-delete`)?.checked || false,
                }))
            }));

            config.rbac = {
                enabled: true,
                roles: roles,
                jwtSecret: document.getElementById('jwtSecret').value,
                modelPerms: modelPerms,
            };
        }

        logText.textContent = '[GORM] 생성 중... 잠시만 기다려주세요';
    } else {
        logText.textContent = `[${currentLanguage.toUpperCase()}] 생성 중... 잠시만 기다려주세요`;
    }

    const lang = gormMode ? 'go' : currentLanguage;

    try {
        const result = await window.go.main.App.GenerateProject(config, lang);

        if (result.success) {
            if (gormMode) {
                logText.textContent = 'GORM Full-Stack 프로젝트 생성 완료!\n터미널에서 실행:\n  cd ' + config.targetPath + '\n  go mod tidy\n  go run .';
            } else if (currentLanguage === 'go') {
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
