import { ModuleDef } from './types';

// 샘플 모듈 정의
export const MODULE_REGISTRY: ModuleDef[] = [
    {
        id: 'auth-login',
        name: 'Simple Login Form',
        description: 'Basic username/password login form with HTMX.',
        category: 'feature',
        snippets: [
            {
                target: 'index.html',
                marker: '<!-- @INJECT_BODY -->',
                content: `
        <div class="mt-8 bg-white p-6 rounded-lg shadow-md">
            <h2 class="text-xl font-bold mb-4">Login</h2>
            <form hx-post="/api/login" hx-target="#login-result" hx-swap="innerHTML">
                <div class="mb-4">
                    <label class="block text-gray-700">Username</label>
                    <input type="text" name="username" class="w-full border rounded p-2" />
                </div>
                <div class="mb-4">
                    <label class="block text-gray-700">Password</label>
                    <input type="password" name="password" class="w-full border rounded p-2" />
                </div>
                <button type="submit" class="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700">Login</button>
            </form>
            <div id="login-result" class="mt-4"></div>
        </div>`
            },
            {
                target: 'main.go',
                marker: '// @INJECT_ROUTES',
                content: `
	mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
		// TODO: Implement actual auth logic
		if username == "admin" && password == "1234" {
			fmt.Fprintf(w, "<div class='text-green-600'>Welcome, Admin!</div>")
		} else {
			fmt.Fprintf(w, "<div class='text-red-500'>Invalid credentials</div>")
		}
	})`
            }
        ]
    },
    {
        id: 'ui-hero',
        name: 'Modern Hero Section',
        description: 'A large hero banner with call to action.',
        category: 'ui',
        snippets: [
            {
                target: 'index.html',
                marker: '<!-- @INJECT_BODY -->',
                content: `
        <div class="bg-gray-900 text-white py-20 px-10 text-center mt-8 rounded-xl">
            <h1 class="text-4xl font-bold mb-4">Build Faster with Ggami</h1>
            <p class="text-xl text-gray-400 mb-8">The ultimate zero-dependency builder for Windows Server.</p>
            <button class="bg-yellow-500 text-black font-bold py-3 px-8 rounded-full hover:bg-yellow-400 transition">Get Started</button>
        </div>`
            }
        ]
    }
];
