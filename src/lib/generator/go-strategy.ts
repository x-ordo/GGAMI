import { IGenerator, ProjectConfig } from './types';
import { GO_MAIN_TEMPLATE, GO_MOD_TEMPLATE } from '../../templates/go/main';
import { HTML_TEMPLATE } from '../../templates/go/htmx';
import { mkdir, writeTextFile } from '@tauri-apps/plugin-fs';

export class GoGenerator implements IGenerator {
    async scaffold(path: string): Promise<void> {
        // 메인 폴더 생성
        await mkdir(path, { recursive: true });
        // templates 폴더 생성
        await mkdir(`${path}/templates`, { recursive: true });
        // assets 폴더 생성
        await mkdir(`${path}/assets`, { recursive: true });
    }

    async createManifest(config: ProjectConfig): Promise<void> {
        const goMod = GO_MOD_TEMPLATE.replace('{{PROJECT_NAME}}', config.projectName);
        await writeTextFile(`${config.targetPath}/go.mod`, goMod);
    }

    async generateCode(config: ProjectConfig): Promise<void> {
        let mainGo = GO_MAIN_TEMPLATE
            .replace('{{DB_SERVER}}', config.dbServer)
            .replace('{{DB_USER}}', config.dbUser)
            .replace('{{DB_PW}}', config.dbPw)
            .replace('{{DB_NAME}}', config.dbName);

        let indexHtml = HTML_TEMPLATE.replaceAll('{{PROJECT_NAME}}', config.projectName);

        // --- Module Injection Logic ---
        const { MODULE_REGISTRY } = await import('../modules/registry');
        
        // 사용자가 선택한 모듈만 필터링
        const activeModules = MODULE_REGISTRY.filter(m => config.modules.includes(m.id));

        for (const module of activeModules) {
            for (const snippet of module.snippets) {
                if (snippet.target === 'main.go') {
                    // 마커 뒤에 코드를 삽입하고, 마커를 다시 붙여서 다음 모듈도 삽입 가능하게 함
                    mainGo = mainGo.replace(snippet.marker, snippet.content + '\n' + snippet.marker);
                } else if (snippet.target === 'index.html') {
                    indexHtml = indexHtml.replace(snippet.marker, snippet.content + '\n' + snippet.marker);
                }
            }
        }
        // ------------------------------

        await writeTextFile(`${config.targetPath}/main.go`, mainGo);
        await writeTextFile(`${config.targetPath}/templates/index.html`, indexHtml);
    }
}
