function __GLOBAL_SET_THEME(theme) {
	localStorage.setItem('theme', theme);

	const resolvedTheme =
		theme === 'system' ? __GLOBAL_GET_SYSTEM_THEME() : theme;

	const meta = document.querySelector('meta[name="theme-color"]');

	if (meta) {
		meta.setAttribute(
			'content',
			resolvedTheme === 'dark' ? '#09090b' : '#ffffff'
		);
	}

	document.documentElement.setAttribute('data-theme', resolvedTheme);
	document.documentElement.setAttribute(
		'style',
		`color-scheme: ${resolvedTheme};`
	);
}

function __GLOBAL_GET_THEME_FROM_LOCAL_STORAGE() {
	return localStorage.getItem('theme');
}

function __GLOBAL_GET_SYSTEM_THEME() {
	return window.matchMedia('(prefers-color-scheme: dark)').matches
		? 'dark'
		: 'light';
}

function __GLOBAL_FIRST_THEME_LOAD() {
	const savedTheme = __GLOBAL_GET_THEME_FROM_LOCAL_STORAGE();
	__GLOBAL_SET_THEME(savedTheme ?? 'system');
}

__GLOBAL_FIRST_THEME_LOAD();

document.addEventListener('DOMContentLoaded', () => {
	window
		.matchMedia('(prefers-color-scheme: dark)')
		.addEventListener('change', e => {
			const savedTheme = __GLOBAL_GET_THEME_FROM_LOCAL_STORAGE();
			if (!savedTheme || savedTheme === 'system') {
				__GLOBAL_SET_THEME('system');
			}
		});
	window.addEventListener('storage', e => {
		if (e.key !== 'theme') return;
		__GLOBAL_SET_THEME(e.newValue ?? 'system');
	});
});
