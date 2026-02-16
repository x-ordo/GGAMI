import { useState } from 'react';
import { generateProject, LanguageType } from './lib/generator'; 
import { open } from '@tauri-apps/plugin-dialog';
import { MODULE_REGISTRY } from './lib/modules/registry';

function App() {
  const [log, setLog] = useState("");
  const [language, setLanguage] = useState<LanguageType>('go');
  const [selectedModules, setSelectedModules] = useState<string[]>([]);
  
  const [config, setConfig] = useState({
    projectName: "GgamiDemo",
    targetPath: "C:\\Projects\\GgamiDemo",
    dbServer: "localhost",
    dbUser: "sa",
    dbPw: "",
    dbName: "MasterDB"
  });

  const handleSelectDir = async () => {
    try {
      const selected = await open({
        directory: true,
        multiple: false,
        defaultPath: "C:\\"
      });
      if (selected && typeof selected === 'string') {
        setConfig({ ...config, targetPath: selected + "\\" + config.projectName });
      }
    } catch (err) {
      console.error("Folder selection failed:", err);
    }
  };

  const toggleModule = (id: string) => {
      setSelectedModules(prev => 
        prev.includes(id) ? prev.filter(m => m !== id) : [...prev, id]
      );
  };

  const handleGenerate = async () => {
    setLog(`[${language.toUpperCase()}] 생성 중... (잠시만 기다려주세요)`);
    
    // Config에 모듈 정보 포함
    const fullConfig = {
        ...config,
        modules: selectedModules
    };

    const result = await generateProject(fullConfig, language);
    
    if (result.success) {
      if (language === 'go') {
        setLog("✅ Go 서버 생성 성공! \n터미널에서:\n  go mod tidy\n  go run .");
      } else {
        setLog("✅ Node.js 서버 생성 성공! \n터미널에서:\n  npm install\n  npm start");
      }
    } else {
      setLog("❌ 실패: " + result.message);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white flex flex-col items-center justify-center p-8">
      <div className="w-full max-w-4xl bg-gray-800 p-8 rounded-lg shadow-xl border border-gray-700 flex gap-8">
        {/* 왼쪽: 설정 패널 */}
        <div className="flex-1 space-y-4">
            <h1 className="text-3xl font-bold mb-6 text-yellow-400">🐶 Ggami Builder</h1>
            
            <div className="flex bg-gray-700 rounded p-1 mb-4">
                <button 
                    onClick={() => setLanguage('go')}
                    className={`flex-1 px-3 py-1 rounded text-sm font-bold transition ${language === 'go' ? 'bg-yellow-500 text-black shadow' : 'text-gray-400 hover:text-white'}`}
                >
                    Go (권장)
                </button>
                <button 
                    onClick={() => setLanguage('node')}
                    className={`flex-1 px-3 py-1 rounded text-sm font-bold transition ${language === 'node' ? 'bg-green-600 text-white shadow' : 'text-gray-400 hover:text-white'}`}
                >
                    Node.js
                </button>
            </div>

            <div>
                <label className="block text-sm font-medium mb-1">프로젝트 이름</label>
                <input 
                type="text" 
                className="w-full bg-gray-700 p-2 rounded border border-gray-600 focus:border-yellow-400 focus:outline-none"
                value={config.projectName}
                onChange={(e) => setConfig({...config, projectName: e.target.value})}
                />
            </div>

            <div className="grid grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium mb-1">DB 서버 IP</label>
                    <input type="text" className="w-full bg-gray-700 p-2 rounded border border-gray-600"
                    value={config.dbServer} onChange={(e) => setConfig({...config, dbServer: e.target.value})} />
                </div>
                <div>
                    <label className="block text-sm font-medium mb-1">DB 이름</label>
                    <input type="text" className="w-full bg-gray-700 p-2 rounded border border-gray-600"
                    value={config.dbName} onChange={(e) => setConfig({...config, dbName: e.target.value})} />
                </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium mb-1">DB 계정 (ID)</label>
                    <input type="text" className="w-full bg-gray-700 p-2 rounded border border-gray-600"
                    value={config.dbUser} onChange={(e) => setConfig({...config, dbUser: e.target.value})} />
                </div>
                <div>
                    <label className="block text-sm font-medium mb-1">DB 암호 (PW)</label>
                    <input type="password" className="w-full bg-gray-700 p-2 rounded border border-gray-600"
                    value={config.dbPw} onChange={(e) => setConfig({...config, dbPw: e.target.value})} />
                </div>
            </div>

            <div>
                <label className="block text-sm font-medium mb-1">저장 경로</label>
                <div className="flex gap-2">
                    <input type="text" readOnly className="flex-1 bg-gray-700 p-2 rounded border border-gray-600 text-gray-400"
                    value={config.targetPath} />
                    <button onClick={handleSelectDir} className="bg-gray-600 px-3 rounded hover:bg-gray-500">📂</button>
                </div>
            </div>
        </div>

        {/* 오른쪽: 모듈 선택 패널 */}
        <div className="w-80 border-l border-gray-700 pl-8">
            <h2 className="text-xl font-bold mb-4 text-gray-300">📦 Modules</h2>
            <div className="space-y-3 h-[400px] overflow-y-auto pr-2">
                {MODULE_REGISTRY.map(module => (
                    <div 
                        key={module.id}
                        onClick={() => toggleModule(module.id)}
                        className={`p-3 rounded border cursor-pointer transition ${
                            selectedModules.includes(module.id) 
                            ? 'bg-blue-900/40 border-blue-500' 
                            : 'bg-gray-800 border-gray-600 hover:border-gray-500'
                        }`}
                    >
                        <div className="flex justify-between items-start">
                            <h3 className="font-bold text-sm">{module.name}</h3>
                            {selectedModules.includes(module.id) && <span className="text-blue-400 text-xs">✔</span>}
                        </div>
                        <p className="text-xs text-gray-400 mt-1">{module.description}</p>
                    </div>
                ))}
            </div>
            
            <button 
                onClick={handleGenerate}
                className={`w-full font-bold py-3 rounded mt-4 transition cursor-pointer ${language === 'go' ? 'bg-yellow-500 text-black hover:bg-yellow-400' : 'bg-green-600 text-white hover:bg-green-500'}`}
            >
                {language === 'go' ? 'Go 프로젝트 생성' : 'Node.js 프로젝트 생성'}
            </button>
        </div>
      </div>
          
        {log && (
          <div className="mt-6 w-full max-w-4xl p-4 bg-black rounded text-sm font-mono text-green-400 whitespace-pre-wrap">
            {log}
          </div>
        )}
    </div>
  );
}

export default App;
