// Generator page logic
let currentLanguage = 'go';
let selectedModules = [];
let allModules = [];

// Initialize on page load
document.addEventListener('DOMContentLoaded', async () => {
    try {
        allModules = await window.go.main.App.GetModules();
        renderModules();
    } catch (err) {
        console.error('Failed to load modules:', err);
    }
});

function setLanguage(lang) {
    currentLanguage = lang;

    const btnGo = document.getElementById('btn-go');
    const btnNode = document.getElementById('btn-node');
    const btnGenerate = document.getElementById('btn-generate');

    if (lang === 'go') {
        btnGo.className = 'flex-1 px-3 py-1 rounded text-sm font-bold transition bg-yellow-500 text-black shadow';
        btnNode.className = 'flex-1 px-3 py-1 rounded text-sm font-bold transition text-gray-400 hover:text-white';
        btnGenerate.className = 'w-full font-bold py-3 rounded mt-4 transition cursor-pointer bg-yellow-500 text-black hover:bg-yellow-400';
        btnGenerate.textContent = 'Generate Go Project';
    } else {
        btnGo.className = 'flex-1 px-3 py-1 rounded text-sm font-bold transition text-gray-400 hover:text-white';
        btnNode.className = 'flex-1 px-3 py-1 rounded text-sm font-bold transition bg-green-600 text-white shadow';
        btnGenerate.className = 'w-full font-bold py-3 rounded mt-4 transition cursor-pointer bg-green-600 text-white hover:bg-green-500';
        btnGenerate.textContent = 'Generate Node.js Project';
    }
}

function renderModules() {
    const list = document.getElementById('module-list');
    list.innerHTML = '';

    for (const mod of allModules) {
        const isSelected = selectedModules.includes(mod.id);
        const el = document.createElement('div');
        el.className = `p-3 rounded border cursor-pointer transition ${
            isSelected
            ? 'bg-blue-900/40 border-blue-500'
            : 'bg-gray-800 border-gray-600 hover:border-gray-500'
        }`;
        el.onclick = () => toggleModule(mod.id);
        el.innerHTML = `
            <div class="flex justify-between items-start">
                <h3 class="font-bold text-sm">${mod.name}</h3>
                ${isSelected ? '<span class="text-blue-400 text-xs">&#10004;</span>' : ''}
            </div>
            <p class="text-xs text-gray-400 mt-1">${mod.description}</p>
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
        console.error('Directory selection failed:', err);
    }
}

async function generateProject() {
    const logEl = document.getElementById('log-output');
    logEl.classList.remove('hidden');
    logEl.textContent = `[${currentLanguage.toUpperCase()}] Generating... please wait`;

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
                logEl.textContent = 'Go server generated!\nIn terminal:\n  go mod tidy\n  go run .';
            } else {
                logEl.textContent = 'Node.js server generated!\nIn terminal:\n  npm install\n  npm start';
            }
        } else {
            logEl.textContent = 'Failed: ' + result.message;
        }
    } catch (err) {
        logEl.textContent = 'Error: ' + err;
    }
}
