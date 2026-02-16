// index.html 템플릿 (HTMX)
export const HTML_TEMPLATE = `<!DOCTYPE html>
<html lang="ko">
<head>
    <meta charset="UTF-8">
    <title>{{PROJECT_NAME}}</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.tailwindcss.com"></script>
    <!-- @INJECT_HEAD -->
</head>
<body class="bg-gray-100 p-10">
    <div class="max-w-md mx-auto bg-white rounded-xl shadow-md overflow-hidden md:max-w-2xl p-6">
        <div class="uppercase tracking-wide text-sm text-indigo-500 font-semibold">{{PROJECT_NAME}}</div>
        <h1 class="block mt-1 text-lg leading-tight font-medium text-black">Windows Server 배포 테스트</h1>
        <p class="mt-2 text-gray-500">Go + HTMX + MSSQL 연결 상태를 확인합니다.</p>
        
        <div class="mt-6">
            <button hx-post="/api/check" hx-target="#result" hx-swap="innerHTML"
                    class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
                DB 연결 테스트 (Click Me)
            </button>
        </div>

        <div id="result" class="mt-4 p-4 bg-gray-50 rounded border border-gray-200 min-h-[60px]">
            결과가 여기에 표시됩니다...
        </div>

        <!-- @INJECT_BODY -->
    </div>
</body>
</html>
`;
