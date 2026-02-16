// Builder page logic
let componentCounter = 0;
let selectedComponentId = null;

// Drag & Drop handlers
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

    // Hide empty state
    const emptyEl = document.getElementById('canvas-empty');
    if (emptyEl) emptyEl.style.display = 'none';

    // Add to project manager via Wails binding
    try {
        await window.go.main.App.AddComponent('page-1', {
            id: compID,
            type: compType,
            content: '',
            styles: { 'class': '' }
        });
    } catch (err) {
        console.log('No project loaded, working in-memory only');
    }

    // Render via HTMX
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

// Component selection
function selectComponent(id, type) {
    selectedComponentId = id;

    // Highlight selected
    document.querySelectorAll('#canvas > div').forEach(el => {
        el.classList.remove('border-blue-400');
        el.classList.add('border-transparent');
    });
    const el = document.getElementById('comp-' + id) || document.getElementById(id);
    if (el) {
        el.classList.remove('border-transparent');
        el.classList.add('border-blue-400');
    }

    // Load property form via HTMX
    htmx.ajax('GET', '/api/properties/' + id + '?type=' + type, '#property-panel');
}

// Delete component
async function deleteComponent(id) {
    const el = document.getElementById('comp-' + id) || document.getElementById(id);
    if (el) el.remove();

    try {
        await window.go.main.App.DeleteComponent('page-1', id);
    } catch (err) {
        // Ignore if no project loaded
    }

    // Reset property panel if deleted component was selected
    if (selectedComponentId === id) {
        selectedComponentId = null;
        document.getElementById('property-panel').innerHTML =
            '<p class="text-sm text-gray-500">Select a component to edit its properties</p>';
    }

    // Show empty state if no components left
    const canvas = document.getElementById('canvas');
    const components = canvas.querySelectorAll('[id^="comp-"]');
    if (components.length === 0) {
        const emptyEl = document.getElementById('canvas-empty');
        if (emptyEl) emptyEl.style.display = 'flex';
    }

    updateComponentCount();
}

// Project operations
async function newProject() {
    const name = prompt('Project name:', 'My Website');
    if (!name) return;

    try {
        await window.go.main.App.CreateBuilderProject(name);
        document.getElementById('project-name').textContent = name;
        clearCanvas();
    } catch (err) {
        alert('Failed to create project: ' + err);
    }
}

async function saveProject() {
    try {
        await window.go.main.App.SaveBuilderProject();
    } catch (err) {
        alert('Failed to save: ' + err);
    }
}

async function loadProject() {
    try {
        const project = await window.go.main.App.LoadBuilderProject();
        if (project) {
            document.getElementById('project-name').textContent = project.name;
            clearCanvas();
            // Re-render all components
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
        alert('Failed to load: ' + err);
    }
}

async function exportHTML() {
    try {
        await window.go.main.App.ExportBuilderHTML();
    } catch (err) {
        alert('Failed to export: ' + err);
    }
}

function clearCanvas() {
    const canvas = document.getElementById('canvas');
    // Remove all component elements but keep the empty state
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
    document.getElementById('component-count').textContent = count + ' component' + (count !== 1 ? 's' : '');
}
