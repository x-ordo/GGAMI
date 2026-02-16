// 비주얼 빌더 페이지 로직
let componentCounter = 0;
let selectedComponentId = null;

// 드래그 앤 드롭 핸들러
function onDragStart(event, type) {
    event.dataTransfer.setData('component-type', type);
    event.dataTransfer.effectAllowed = 'copy';
}

function onDragOver(event) {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'copy';
    document.getElementById('canvas').classList.add('drag-over');
}

function onDragLeave(event) {
    document.getElementById('canvas').classList.remove('drag-over');
}

async function onDrop(event) {
    event.preventDefault();
    document.getElementById('canvas').classList.remove('drag-over');

    const compType = event.dataTransfer.getData('component-type');
    if (!compType) return;

    componentCounter++;
    const compID = 'comp-' + Date.now() + '-' + componentCounter;

    // 빈 상태 숨기기
    const emptyEl = document.getElementById('canvas-empty');
    if (emptyEl) emptyEl.style.display = 'none';

    // Wails 바인딩을 통해 프로젝트 매니저에 추가
    try {
        await window.go.main.App.AddComponent('page-1', {
            id: compID,
            type: compType,
            content: '',
            styles: { 'class': '' }
        });
    } catch (err) {
        console.log('프로젝트 미로드 상태, 메모리에서만 작업');
    }

    // HTMX로 렌더링
    const formData = new FormData();
    formData.append('type', compType);
    formData.append('id', compID);

    const response = await fetch('/api/canvas/add', {
        method: 'POST',
        body: formData
    });
    const html = await response.text();

    const canvas = document.getElementById('canvas');
    const wrapper = document.createElement('div');
    wrapper.innerHTML = html;
    canvas.appendChild(wrapper.firstElementChild);

    updateComponentCount();
}

// 컴포넌트 선택
function selectComponent(id, type) {
    selectedComponentId = id;

    // 선택 하이라이트
    document.querySelectorAll('#canvas > div').forEach(el => {
        el.classList.remove('border-primary');
        el.classList.add('border-transparent');
    });
    const el = document.getElementById('comp-' + id) || document.getElementById(id);
    if (el) {
        el.classList.remove('border-transparent');
        el.classList.add('border-primary');
    }

    // HTMX로 속성 폼 로드
    htmx.ajax('GET', '/api/properties/' + id + '?type=' + type, '#property-panel');
}

// 컴포넌트 삭제
async function deleteComponent(id) {
    const el = document.getElementById('comp-' + id) || document.getElementById(id);
    if (el) el.remove();

    try {
        await window.go.main.App.DeleteComponent('page-1', id);
    } catch (err) {
        // 프로젝트 미로드 시 무시
    }

    // 삭제된 컴포넌트가 선택 상태였으면 속성 패널 초기화
    if (selectedComponentId === id) {
        selectedComponentId = null;
        document.getElementById('property-panel').innerHTML =
            '<p class="text-sm text-base-content/50">컴포넌트를 선택하면 속성을 편집할 수 있습니다</p>';
    }

    // 컴포넌트가 없으면 빈 상태 표시
    const canvas = document.getElementById('canvas');
    const components = canvas.querySelectorAll('[id^="comp-"]');
    if (components.length === 0) {
        const emptyEl = document.getElementById('canvas-empty');
        if (emptyEl) emptyEl.style.display = 'flex';
    }

    updateComponentCount();
}

// 프로젝트 작업
async function newProject() {
    const name = prompt('프로젝트 이름:', '내 웹사이트');
    if (!name) return;

    try {
        await window.go.main.App.CreateBuilderProject(name);
        document.getElementById('project-name').textContent = name;
        clearCanvas();
    } catch (err) {
        alert('프로젝트 생성 실패: ' + err);
    }
}

async function saveProject() {
    try {
        await window.go.main.App.SaveBuilderProject();
    } catch (err) {
        alert('저장 실패: ' + err);
    }
}

async function loadProject() {
    try {
        const project = await window.go.main.App.LoadBuilderProject();
        if (project) {
            document.getElementById('project-name').textContent = project.name;
            clearCanvas();
            // 모든 컴포넌트 다시 렌더링
            if (project.pages && project.pages.length > 0) {
                for (const comp of project.pages[0].components) {
                    const formData = new FormData();
                    formData.append('type', comp.type);
                    formData.append('id', comp.id);

                    const response = await fetch('/api/canvas/add', {
                        method: 'POST',
                        body: formData
                    });
                    const html = await response.text();
                    const canvas = document.getElementById('canvas');
                    const wrapper = document.createElement('div');
                    wrapper.innerHTML = html;
                    canvas.appendChild(wrapper.firstElementChild);
                }
                const emptyEl = document.getElementById('canvas-empty');
                if (emptyEl && project.pages[0].components.length > 0) {
                    emptyEl.style.display = 'none';
                }
            }
            updateComponentCount();
        }
    } catch (err) {
        alert('불러오기 실패: ' + err);
    }
}

async function exportHTML() {
    try {
        await window.go.main.App.ExportBuilderHTML();
    } catch (err) {
        alert('내보내기 실패: ' + err);
    }
}

function clearCanvas() {
    const canvas = document.getElementById('canvas');
    const children = Array.from(canvas.children);
    for (const child of children) {
        if (child.id !== 'canvas-empty') {
            child.remove();
        }
    }
    const emptyEl = document.getElementById('canvas-empty');
    if (emptyEl) emptyEl.style.display = 'flex';
    updateComponentCount();
}

function updateComponentCount() {
    const canvas = document.getElementById('canvas');
    const count = canvas.querySelectorAll('[id^="comp-"]').length;
    document.getElementById('component-count').textContent = '컴포넌트 ' + count + '개';
}
