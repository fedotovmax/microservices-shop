const esbuild = require('esbuild');
const fs = require('fs');
const path = require('path');

const pagesDir = './web/dev/scripts/pages';
const outDir = './web/public/js';

// собираем entry points автоматически
const entryPoints = Object.fromEntries(
	fs
		.readdirSync(pagesDir)
		.filter(f => f.endsWith('.ts'))
		.map(f => [path.basename(f, '.ts'), path.join(pagesDir, f)])
);

esbuild
	.build({
		entryPoints,
		outdir: outDir,

		bundle: true,
		format: 'esm',
		target: 'es2020',
		splitting: true,

		sourcemap: false,
		minify: true,
		treeShaking: true,

		entryNames: '[name]',
		chunkNames: 'chunks/[name]-[hash]',

		loader: {
			'.ts': 'ts',
		},

		// define: {
		// 	__DEV__: JSON.stringify(isDev),
		// },
	})
	.catch(() => process.exit(1));
