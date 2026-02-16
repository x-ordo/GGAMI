export interface TableSchema {
    name: string;
    columns: any[]; // 추후 구체화
}

export interface ProjectConfig {
    projectName: string;
    targetPath: string;
    dbServer: string;
    dbUser: string;
    dbPw: string;
    dbName: string;
    modules: string[]; // Selected Module IDs
    tables?: TableSchema[];
}

export interface IGenerator {
    // 1. 기본 폴더 구조 생성
    scaffold(path: string): Promise<void>;
  
    // 2. 설정 파일 생성 (go.mod 또는 package.json)
    createManifest(config: ProjectConfig): Promise<void>;
  
    // 3. 소스 코드 생성 (main.go 또는 app.module.ts)
    generateCode(config: ProjectConfig): Promise<void>;
  
    // 4. 빌드 명령어 실행 (go build 또는 npm install)
    // buildProject(path: string): Promise<void>; // 추후 구현
}
