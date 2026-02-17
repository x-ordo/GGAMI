export namespace builder {
	
	export class Component {
	    id: string;
	    type: string;
	    content: string;
	    styles?: Record<string, string>;
	    children?: Component[];
	
	    static createFrom(source: any = {}) {
	        return new Component(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.content = source["content"];
	        this.styles = source["styles"];
	        this.children = this.convertValues(source["children"], Component);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Page {
	    id: string;
	    name: string;
	    components: Component[];
	
	    static createFrom(source: any = {}) {
	        return new Page(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.components = this.convertValues(source["components"], Component);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BuilderProject {
	    name: string;
	    pages: Page[];
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new BuilderProject(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.pages = this.convertValues(source["pages"], Page);
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace domain {
	
	export class CodeSnippet {
	    target: string;
	    marker: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new CodeSnippet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.target = source["target"];
	        this.marker = source["marker"];
	        this.content = source["content"];
	    }
	}
	export class FieldDef {
	    name: string;
	    type: string;
	    gormTags: string[];
	    defaultVal: string;
	    jsonName: string;
	
	    static createFrom(source: any = {}) {
	        return new FieldDef(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.gormTags = source["gormTags"];
	        this.defaultVal = source["defaultVal"];
	        this.jsonName = source["jsonName"];
	    }
	}
	export class ModelDef {
	    name: string;
	    fields: FieldDef[];
	
	    static createFrom(source: any = {}) {
	        return new ModelDef(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.fields = this.convertValues(source["fields"], FieldDef);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RolePermission {
	    role: string;
	    create: boolean;
	    read: boolean;
	    update: boolean;
	    delete: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RolePermission(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.create = source["create"];
	        this.read = source["read"];
	        this.update = source["update"];
	        this.delete = source["delete"];
	    }
	}
	export class ModelRBAC {
	    modelName: string;
	    permissions: RolePermission[];
	
	    static createFrom(source: any = {}) {
	        return new ModelRBAC(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.modelName = source["modelName"];
	        this.permissions = this.convertValues(source["permissions"], RolePermission);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ModuleDef {
	    id: string;
	    name: string;
	    description: string;
	    category: string;
	    dependencies: string[];
	    snippets: CodeSnippet[];
	
	    static createFrom(source: any = {}) {
	        return new ModuleDef(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.category = source["category"];
	        this.dependencies = source["dependencies"];
	        this.snippets = this.convertValues(source["snippets"], CodeSnippet);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RBACConfig {
	    enabled: boolean;
	    roles: string[];
	    jwtSecret: string;
	    modelPerms: ModelRBAC[];
	
	    static createFrom(source: any = {}) {
	        return new RBACConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.roles = source["roles"];
	        this.jwtSecret = source["jwtSecret"];
	        this.modelPerms = this.convertValues(source["modelPerms"], ModelRBAC);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProjectConfig {
	    projectName: string;
	    targetPath: string;
	    dbServer: string;
	    dbUser: string;
	    dbPw: string;
	    dbName: string;
	    port?: number;
	    modules: string[];
	    gormMode?: boolean;
	    models?: ModelDef[];
	    dbType?: string;
	    rbac?: RBACConfig;
	
	    static createFrom(source: any = {}) {
	        return new ProjectConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectName = source["projectName"];
	        this.targetPath = source["targetPath"];
	        this.dbServer = source["dbServer"];
	        this.dbUser = source["dbUser"];
	        this.dbPw = source["dbPw"];
	        this.dbName = source["dbName"];
	        this.port = source["port"];
	        this.modules = source["modules"];
	        this.gormMode = source["gormMode"];
	        this.models = this.convertValues(source["models"], ModelDef);
	        this.dbType = source["dbType"];
	        this.rbac = this.convertValues(source["rbac"], RBACConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

