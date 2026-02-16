import { IGenerator, ProjectConfig } from './types';
import { GoGenerator } from './go-strategy';
import { NodeGenerator } from './node-strategy';

export type LanguageType = 'go' | 'node';

export class GeneratorFactory {
    static create(type: LanguageType): IGenerator {
        switch (type) {
            case 'go':
                return new GoGenerator();
            case 'node':
                return new NodeGenerator();
            default:
                throw new Error("Unsupported language type: " + type);
        }
    }
}

export async function generateProject(config: ProjectConfig, language: LanguageType = 'go') {
    try {
        console.log(`🚀 프로젝트 생성 시작 (${language}):`, config.projectName);

        const generator = GeneratorFactory.create(language);
        
        // 1. Scaffold
        await generator.scaffold(config.targetPath);
        
        // 2. Mainfest
        await generator.createManifest(config);
        
        // 3. Code
        await generator.generateCode(config);

        console.log("✅ 파일 생성 완료!");
        return { success: true, message: "생성 완료: " + config.targetPath };

    } catch (error) {
        console.error("❌ 생성 실패:", error);
        return { success: false, message: String(error) };
    }
}
