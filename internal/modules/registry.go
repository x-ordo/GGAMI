package modules

// Registry contains all available modules
var Registry = []ModuleDef{
	{
		ID:          "auth-login",
		Name:        "Simple Login Form",
		Description: "Basic username/password login form with HTMX.",
		Category:    "feature",
		Snippets: []CodeSnippet{
			{
				Target: TargetIndexHTML,
				Marker: "<!-- @INJECT_BODY -->",
				Content: `
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
        </div>`,
			},
			{
				Target: TargetMainGo,
				Marker: "// @INJECT_ROUTES",
				Content: `
	mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
		// TODO: Implement actual auth logic
		if username == "admin" && password == "1234" {
			fmt.Fprintf(w, "<div class='text-green-600'>Welcome, Admin!</div>")
		} else {
			fmt.Fprintf(w, "<div class='text-red-500'>Invalid credentials</div>")
		}
	})`,
			},
		},
	},
	{
		ID:          "ui-hero",
		Name:        "Modern Hero Section",
		Description: "A large hero banner with call to action.",
		Category:    "ui",
		Snippets: []CodeSnippet{
			{
				Target: TargetIndexHTML,
				Marker: "<!-- @INJECT_BODY -->",
				Content: `
        <div class="bg-gray-900 text-white py-20 px-10 text-center mt-8 rounded-xl">
            <h1 class="text-4xl font-bold mb-4">Build Faster with Ggami</h1>
            <p class="text-xl text-gray-400 mb-8">The ultimate zero-dependency builder for Windows Server.</p>
            <button class="bg-yellow-500 text-black font-bold py-3 px-8 rounded-full hover:bg-yellow-400 transition">Get Started</button>
        </div>`,
			},
		},
	},
}
