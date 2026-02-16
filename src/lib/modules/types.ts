export type InjectionTarget = 'main.go' | 'go.mod' | 'index.html';

export interface CodeSnippet {
    target: InjectionTarget;
    marker: string; // e.g., "@INJECT_ROUTES"
    content: string; // The code to inject
}

export interface ModuleDef {
    id: string;
    name: string;
    description: string;
    category: 'feature' | 'ui' | 'utils';
    snippets: CodeSnippet[];
}
